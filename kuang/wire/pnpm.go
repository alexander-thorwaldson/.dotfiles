package wire

// PnpmRunRequest is the input for running a pnpm script.
type PnpmRunRequest struct {
	Dir    string `json:"dir"`
	Script string `json:"script"`
}

// PnpmRunResponse is returned after running a pnpm script.
type PnpmRunResponse struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exitCode"`
}

// Clone returns a copy of PnpmRunResponse.
func (r PnpmRunResponse) Clone() PnpmRunResponse { return r }
