package dispatcher

// ###########################################################################
// ###########################################################################
// Dispatcher Response
// ###########################################################################
// ###########################################################################

// Response encapsulates the data o a response
type Response struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Namespace string `json:"namespace,omitempty"`
	Status    string `json:"status"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
}

// ###########################################################################

// InitResponse ...
func InitResponse(r *Response, host string, port int, namespace string, name string, hostname string, version string, status string) {
	r.Host = host
	r.Port = port
	r.Namespace = namespace
	r.Name = name
	r.Hostname = hostname
	r.Version = version
	r.Status = status
}

// // ---------------------------------------------------------------------------

// InitResponseFromDispatcher ...
func InitResponseFromDispatcher(r *Response, ds IConfiguration, status string) {
	InitResponse(r, ds.GetHost(), ds.GetPort(), ds.GetNamespace(), ds.GetName(), ds.GetHostname(), ds.GetVersion(), status)
}
