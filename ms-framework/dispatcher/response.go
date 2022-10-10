package dispatcher

// ###########################################################################
// ###########################################################################
// Dispatcher Response
// ###########################################################################
// ###########################################################################

// Response encapsulates the data o a response
type BasicResponse struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Namespace string `json:"namespace,omitempty"`
	Status    string `json:"status"`
	Code      int    `json:"code"`
	Name      string `json:"name"`
	Hostname  string `json:"hostname"`
	Version   string `json:"version"`
}

type Trace struct {
	Namespace string  `json:"namespace,omitempty"`
	Name      string  `json:"name"`
	Version   string  `json:"version"`
	Hostname  string  `json:"hostname"`
	Code      int     `json:"code"`
	Status    string  `json:"status"`
	Traces    []Trace `json:"traces,omitempty"`
}

type Response struct {
	BasicResponse
	Trace Trace `json:"trace"`
}

// ###########################################################################

// InitResponse ...
func InitResponse(r *Response, host string, port int, namespace string, name string, hostname string, version string, code int, status string) {
	r.Host = host
	r.Port = port
	r.Namespace = namespace
	r.Name = name
	r.Hostname = hostname
	r.Version = version
	r.Code = code
	r.Status = status
	InitTrace(&r.Trace, namespace, name, hostname, version, code, status)
}

// InitTrace ...
func InitTrace(t *Trace, namespace string, name string, hostname string, version string, code int, status string) {
	t.Namespace = namespace
	t.Name = name
	t.Version = version
	t.Hostname = hostname
	t.Code = code
	t.Status = status
}

// // ---------------------------------------------------------------------------

// InitResponseFromDispatcher ...
func InitResponseFromDispatcher(r *Response, ds IConfiguration, code int, status string) {
	InitResponse(r, ds.GetHost(), ds.GetPort(), ds.GetNamespace(), ds.GetName(), ds.GetHostname(), ds.GetVersion(), code, status)
}

// InitTraceFromDispatcher ...
func InitTraceFromDispatcher(t *Trace, ds IConfiguration, code int, status string) {
	InitTrace(t, ds.GetNamespace(), ds.GetName(), ds.GetHostname(), ds.GetVersion(), code, status)
}
