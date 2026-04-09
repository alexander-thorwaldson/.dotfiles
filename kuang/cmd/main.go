// Package main is the entry point for the kuang API server.
package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/zoobz-io/sctx"
	"github.com/zoobz-io/sum"
	"github.com/zoobzio/dotfiles/kuang/contracts"
	"github.com/zoobzio/dotfiles/kuang/handlers"
	"github.com/zoobzio/dotfiles/kuang/services"
)

func main() {
	os.Exit(run())
}

func run() int {
	// ── 1. Certs ─────────────────────────────────
	certsDir := os.Getenv("KUANG_CERTS_DIR")
	if certsDir == "" {
		certsDir = "certs"
	}

	caPool, err := services.LoadCAPool(certsDir + "/root_ca.crt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading CA pool: %v\n", err)
		return 1
	}

	kuangCert, kuangKey, err := services.LoadKeyPair(certsDir+"/kuang.crt", certsDir+"/kuang.key")
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading kuang key pair: %v\n", err)
		return 1
	}

	tlsCfg, err := services.LoadTLSConfig(certsDir+"/kuang.crt", certsDir+"/kuang.key", certsDir+"/root_ca.crt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuring TLS: %v\n", err)
		return 1
	}

	// ── 2. sctx ──────────────────────────────────
	roleExpand := func(role string) []string {
		return services.ExpandRole(services.Role(role))
	}

	admin, err := sctx.NewAdminService[services.AgentMeta](kuangKey, caPool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating sctx admin: %v\n", err)
		return 1
	}
	if policyErr := admin.SetPolicy(services.AgentPolicy(roleExpand)); policyErr != nil {
		fmt.Fprintf(os.Stderr, "setting sctx policy: %v\n", policyErr)
		return 1
	}

	// Authenticate kuang itself so it has a valid context.
	if _, authErr := admin.GenerateTrusted(context.Background(), kuangCert); authErr != nil {
		fmt.Fprintf(os.Stderr, "authenticating kuang: %v\n", authErr)
		return 1
	}

	// ── 3. ICE + Registry ─────────────────────────
	iceEndpoint := os.Getenv("KUANG_ICE_ENDPOINT")
	if iceEndpoint == "" {
		iceEndpoint = "http://127.0.0.1:9119"
	}
	ice := services.NewICEClient(iceEndpoint)

	k := sum.Start()

	sum.Register[contracts.GH](k, services.NewGH(ice))
	sum.Register[contracts.Jira](k, services.NewJira(ice))
	sum.Register[contracts.Dev](k, services.NewDev(ice))
	sum.Register[contracts.Pnpm](k, services.NewPnpm(ice))

	sum.Freeze(k)

	// ── 4. Server ────────────────────────────────
	svc := sum.New()
	svc.Engine().
		WithTLSConfig(tlsCfg).
		WithAuthenticator(services.Authenticator(admin, roleExpand))
	svc.Handle(handlers.All()...)

	port := 7117
	if p := os.Getenv("KUANG_PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	fmt.Fprintf(os.Stderr, "kuang listening on :%d (mTLS)\n", port)
	if err := svc.Run("", port); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		return 1
	}
	return 0
}
