package services

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/zoobz-io/rocco"
	"github.com/zoobz-io/sctx"
)

// AgentIdentity implements rocco.Identity for mTLS-authenticated agents.
type AgentIdentity struct {
	name   string
	role   string
	scopes []string
}

// ID returns the agent's common name.
func (a *AgentIdentity) ID() string { return a.name }

// TenantID is unused — returns empty.
func (a *AgentIdentity) TenantID() string { return "" }

// Email is unused — returns empty.
func (a *AgentIdentity) Email() string { return "" }

// Scopes returns the agent's tool permissions derived from its role.
func (a *AgentIdentity) Scopes() []string { return a.scopes }

// Roles returns the agent's role.
func (a *AgentIdentity) Roles() []string { return []string{a.role} }

// HasScope checks if the agent has a specific permission.
func (a *AgentIdentity) HasScope(scope string) bool {
	for _, s := range a.scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// HasRole checks if the agent has a specific role.
func (a *AgentIdentity) HasRole(role string) bool { return a.role == role }

// Stats is unused.
func (a *AgentIdentity) Stats() map[string]int { return nil }

// Role returns the agent's role string.
func (a *AgentIdentity) Role() string { return a.role }

// RoleExpander resolves a role string into its tool permissions.
type RoleExpander func(role string) []string

// Authenticator creates a rocco-compatible identity extractor that reads
// the client certificate from the mTLS connection, generates an sctx token,
// and returns an AgentIdentity with role-based scopes.
func Authenticator(admin sctx.Admin[AgentMeta], expand RoleExpander) func(context.Context, *http.Request) (rocco.Identity, error) {
	return func(ctx context.Context, r *http.Request) (rocco.Identity, error) {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			return nil, fmt.Errorf("client certificate required")
		}

		cert := r.TLS.PeerCertificates[0]

		// Generate or retrieve cached token.
		_, err := admin.GenerateTrusted(ctx, cert)
		if err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}

		// Resolve role from cert OU.
		role := "admin"
		if cert.Subject.CommonName != "kuang" && len(cert.Subject.OrganizationalUnit) > 0 {
			role = cert.Subject.OrganizationalUnit[0]
		}

		scopes := expand(role)
		return &AgentIdentity{
			name:   cert.Subject.CommonName,
			role:   role,
			scopes: scopes,
		}, nil
	}
}

// AgentMeta is the sctx metadata type for agent security contexts.
type AgentMeta struct {
	AgentName string `json:"agentName"`
	Role      string `json:"role"`
}

// AgentPolicy creates an sctx context policy that maps cert OU to permissions.
func AgentPolicy(expand RoleExpander) sctx.ContextPolicy[AgentMeta] {
	return func(cert *x509.Certificate) (*sctx.Context[AgentMeta], error) {
		role := "admin"
		if cert.Subject.CommonName != "kuang" {
			if len(cert.Subject.OrganizationalUnit) == 0 {
				return nil, fmt.Errorf("certificate %q has no OU field (role required)", cert.Subject.CommonName)
			}
			role = cert.Subject.OrganizationalUnit[0]
		}

		perms := expand(role)
		if perms == nil {
			return nil, fmt.Errorf("unknown role %q", role)
		}

		return &sctx.Context[AgentMeta]{
			Permissions: perms,
			Metadata:    AgentMeta{AgentName: cert.Subject.CommonName, Role: role},
			ExpiresAt:   cert.NotAfter,
		}, nil
	}
}

// LoadTLSConfig creates an mTLS server config.
func LoadTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("loading server cert: %w", err)
	}

	caCert, err := os.ReadFile(caFile) // #nosec G304 -- controlled path //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("reading CA cert: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA cert")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
	}, nil
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

// LoadCAPool loads a CA certificate pool from a PEM file.
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
