package wire

// DevStartRequest is the input for starting a dev server.
type DevStartRequest struct {
	Repo    string `json:"repo"`
	Dir     string `json:"dir"`
	Command string `json:"command,omitempty"`
	Port    int    `json:"port"`
	Timeout int    `json:"timeout,omitempty"`
}

// DevStartResponse is returned when a dev server is started.
type DevStartResponse struct {
	Repo   string `json:"repo"`
	Status string `json:"status"`
	URL    string `json:"url"`
	Port   int    `json:"port"`
}

// Clone returns a copy of DevStartResponse.
func (r DevStartResponse) Clone() DevStartResponse { return r }

// DevStopRequest is the input for stopping a dev server.
type DevStopRequest struct {
	Repo string `json:"repo"`
}

// DevStopResponse is returned when a dev server is stopped.
type DevStopResponse struct {
	Repo    string `json:"repo"`
	Stopped bool   `json:"stopped"`
}

// Clone returns a copy of DevStopResponse.
func (r DevStopResponse) Clone() DevStopResponse { return r }

// DevStatusResponse is returned by the dev status endpoint.
type DevStatusResponse struct {
	Services []DevServiceInfo `json:"services"`
}

// DevServiceInfo describes a running dev server.
type DevServiceInfo struct {
	Repo      string `json:"repo"`
	Dir       string `json:"dir"`
	Command   string `json:"command"`
	Status    string `json:"status"`
	URL       string `json:"url"`
	StartedAt string `json:"startedAt"`
	Port      int    `json:"port"`
	PID       int    `json:"pid"`
}

// Clone returns a deep copy of DevStatusResponse.
func (r DevStatusResponse) Clone() DevStatusResponse {
	svcs := make([]DevServiceInfo, len(r.Services))
	copy(svcs, r.Services)
	return DevStatusResponse{Services: svcs}
}

// DevLogRequest is the input for reading dev server logs.
type DevLogRequest struct {
	Repo  string `json:"repo"`
	Lines int    `json:"lines,omitempty"`
}

// DevLogResponse is returned with recent log output.
type DevLogResponse struct {
	Repo  string   `json:"repo"`
	Lines []string `json:"lines"`
}

// Clone returns a deep copy of DevLogResponse.
func (r DevLogResponse) Clone() DevLogResponse {
	lines := make([]string, len(r.Lines))
	copy(lines, r.Lines)
	return DevLogResponse{Repo: r.Repo, Lines: lines}
}
