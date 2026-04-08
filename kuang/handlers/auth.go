package handlers

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zoobz-io/sctx"
)

// AgentMeta is the metadata type for agent security contexts.
type AgentMeta struct {
	AgentName string `json:"agentName"`
	Role      Role   `json:"role"`
}

// Auth holds the sctx admin service and kuang's principal for creating guards.
type Auth struct {
	Admin     sctx.Admin[AgentMeta]
	Principal sctx.Principal
	Logger    *slog.Logger
}

// NewAuth initializes the sctx admin service and authenticates kuang's own certificate.
func NewAuth(logger *slog.Logger, privateKey crypto.PrivateKey, cert *x509.Certificate, caPool *x509.CertPool) (*Auth, error) {
	admin, err := sctx.NewAdminService[AgentMeta](privateKey, caPool)
	if err != nil {
		return nil, fmt.Errorf("creating admin service: %w", err)
	}

	// Policy: map certificate OU (role) to tool permissions.
	if policyErr := admin.SetPolicy(agentPolicy); policyErr != nil {
		return nil, fmt.Errorf("setting policy: %w", policyErr)
	}

	ctx := context.Background()
	principal, err := sctx.NewPrincipal[AgentMeta](ctx, admin, privateKey, cert)
	if err != nil {
		return nil, fmt.Errorf("creating kuang principal: %w", err)
	}

	return &Auth{
		Admin:     admin,
		Principal: principal,
		Logger:    logger,
	}, nil
}

// agentPolicy maps certificate fields to permissions and metadata.
// The cert's first OU field is treated as the role. Permissions are
// resolved from the role via ExpandRole.
func agentPolicy(cert *x509.Certificate) (*sctx.Context[AgentMeta], error) {
	// Kuang itself (CN=kuang) gets admin role.
	// Context expiry matches the certificate's NotAfter.
	expiry := cert.NotAfter

	if cert.Subject.CommonName == "kuang" {
		return &sctx.Context[AgentMeta]{
			Permissions: AllToolPermissions(),
			Metadata:    AgentMeta{AgentName: "kuang", Role: RoleAdmin},
			IssuedAt:    time.Now(),
			ExpiresAt:   expiry,
		}, nil
	}

	if len(cert.Subject.OrganizationalUnit) == 0 {
		return nil, fmt.Errorf("certificate %q has no OU field (role required)", cert.Subject.CommonName)
	}

	role := Role(cert.Subject.OrganizationalUnit[0])
	perms := ExpandRole(role)
	if perms == nil {
		return nil, fmt.Errorf("unknown role %q for certificate %q", role, cert.Subject.CommonName)
	}

	return &sctx.Context[AgentMeta]{
		Permissions: perms,
		Metadata:    AgentMeta{AgentName: cert.Subject.CommonName, Role: role},
		IssuedAt:    time.Now(),
		ExpiresAt:   expiry,
	}, nil
}

// CreateGuard creates a guard for the given permission using kuang's principal.
func (a *Auth) CreateGuard(ctx context.Context, permission string) (sctx.Guard, error) {
	return a.Principal.Guard(ctx, permission)
}

// GuardedHandler wraps an MCP tool handler with sctx guard validation.
// The guard checks that the caller's token (injected by TokenInjector) has
// the required permission.
func GuardedHandler[In, Out any](guard sctx.Guard, logger *slog.Logger, name string, fn mcp.ToolHandlerFor[In, Out]) mcp.ToolHandlerFor[In, Out] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input In) (*mcp.CallToolResult, Out, error) {
		var zero Out

		if err := guard.Validate(ctx); err != nil {
			logger.Warn("guard rejected tool call", "tool", name, "err", err)
			r := &mcp.CallToolResult{}
			r.SetError(fmt.Errorf("unauthorized: %w", err))
			return r, zero, nil
		}

		return fn(ctx, req, input)
	}
}

// TLSConfig creates a TLS config for kuang's HTTPS server with mTLS.
func TLSConfig(certFile, keyFile string, caPool *x509.CertPool) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("loading server cert: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// LoadCAPool loads the CA root certificate into a cert pool.
func LoadCAPool(caFile string) (*x509.CertPool, error) {
	caCert, err := os.ReadFile(caFile) // #nosec G304 -- controlled path //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("reading CA cert: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA cert")
	}
	return pool, nil
}

// LoadKeyPair loads a certificate and private key from PEM files.
func LoadKeyPair(certFile, keyFile string) (*x509.Certificate, crypto.PrivateKey, error) {
	tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("loading key pair: %w", err)
	}
	cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, nil, fmt.Errorf("parsing certificate: %w", err)
	}
	return cert, tlsCert.PrivateKey, nil
}

// TokenInjector is an HTTP middleware that extracts the client certificate from
// the mTLS connection, generates an sctx token via GenerateTrusted, and injects
// it into the request context for downstream guards to validate.
func TokenInjector(admin sctx.Admin[AgentMeta], logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			http.Error(w, "client certificate required", http.StatusUnauthorized)
			return
		}

		clientCert := r.TLS.PeerCertificates[0]
		token, err := admin.GenerateTrusted(r.Context(), clientCert)
		if err != nil {
			logger.Error("failed to generate token", "cn", clientCert.Subject.CommonName, "err", err)
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}

		ctx := sctx.InjectToken(r.Context(), token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
