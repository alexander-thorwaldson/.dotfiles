package handlers

import (
	"github.com/zoobz-io/rocco"
	"github.com/zoobz-io/sum"
	"github.com/zoobzio/dotfiles/kuang/contracts"
	"github.com/zoobzio/dotfiles/kuang/wire"
)

var devStart = rocco.POST[wire.DevStartRequest, wire.DevStartResponse]("/v1/dev/start", func(r *rocco.Request[wire.DevStartRequest]) (wire.DevStartResponse, error) {
	dev := sum.MustUse[contracts.Dev](r)
	resp, err := dev.Start(r, r.Body)
	if err != nil {
		return wire.DevStartResponse{}, err
	}
	return *resp, nil
}).
	WithSummary("Start a dev server").
	WithTags("dev").
	WithAuthentication().
	WithScopes("dev_start")

var devStop = rocco.POST[wire.DevStopRequest, wire.DevStopResponse]("/v1/dev/stop", func(r *rocco.Request[wire.DevStopRequest]) (wire.DevStopResponse, error) {
	dev := sum.MustUse[contracts.Dev](r)
	resp, err := dev.Stop(r, r.Body)
	if err != nil {
		return wire.DevStopResponse{}, err
	}
	return *resp, nil
}).
	WithSummary("Stop a dev server").
	WithTags("dev").
	WithAuthentication().
	WithScopes("dev_stop")

var devStatus = rocco.GET[rocco.NoBody, wire.DevStatusResponse]("/v1/dev/status", func(r *rocco.Request[rocco.NoBody]) (wire.DevStatusResponse, error) {
	dev := sum.MustUse[contracts.Dev](r)
	resp, err := dev.Status(r)
	if err != nil {
		return wire.DevStatusResponse{}, err
	}
	return *resp, nil
}).
	WithSummary("List running dev servers").
	WithTags("dev").
	WithAuthentication().
	WithScopes("dev_status")

var devLog = rocco.POST[wire.DevLogRequest, wire.DevLogResponse]("/v1/dev/log", func(r *rocco.Request[wire.DevLogRequest]) (wire.DevLogResponse, error) {
	dev := sum.MustUse[contracts.Dev](r)
	resp, err := dev.Log(r, r.Body)
	if err != nil {
		return wire.DevLogResponse{}, err
	}
	return *resp, nil
}).
	WithSummary("Read dev server logs").
	WithTags("dev").
	WithAuthentication().
	WithScopes("dev_log")
