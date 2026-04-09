package handlers

import (
	"github.com/zoobz-io/rocco"
	"github.com/zoobz-io/sum"
	"github.com/zoobzio/dotfiles/kuang/contracts"
	"github.com/zoobzio/dotfiles/kuang/wire"
)

var pnpmRun = rocco.POST[wire.PnpmRunRequest, wire.PnpmRunResponse]("/v1/pnpm/run", func(r *rocco.Request[wire.PnpmRunRequest]) (wire.PnpmRunResponse, error) {
	pnpm := sum.MustUse[contracts.Pnpm](r)
	resp, err := pnpm.Run(r, r.Body)
	if err != nil {
		return wire.PnpmRunResponse{}, err
	}
	return *resp, nil
}).
	WithSummary("Run a pnpm script").
	WithTags("pnpm").
	WithAuthentication().
	WithScopes("pnpm_run")
