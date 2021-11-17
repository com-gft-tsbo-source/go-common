package dispatcher

import (
	"net/http"
)

type RequestHeaderFunction func(out *http.Request, in *http.Request)
type ResponseHeaderFunction func(out http.ResponseWriter, in *http.Request)

// ConnectionConfiguration ...
type ConnectionConfiguration struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Hostname  string `json:"hostname"`
	Version   string `json:"version"`
}

// IConnectionConfiguration ...
type IConnectionConfiguration interface {
	// Debug(s string)
	GetHost() string
	GetPort() int
	GetNamespace() string
	GetName() string
	GetHostname() string
	GetVersion() string
}

// TLSConfiguration ...
type TLSConfiguration struct {
	CertChainFile string `json:"certchainfile"`
	KeyFile       string `json:"keyfile"`
	CAFile        string `json:"cafile"`
}

// ITLSConfiguration ...
type ITLSConfiguration interface {
	GetCertChainFile() string
	GetKeyFile() string
	GetCAFile() string
}

// LimitConfiguration ...
type LimitConfiguration struct {
	MaxConnections int `json:"maxconnections"`
	DelayReply     int `json:"delayreply"`
	ClientTimeout  int `json:"clienttimeout"`
}

// ILimitConfiguration ...
type ILimitConfiguration interface {
	GetMaxConnections() int
	GetDelayReply() int
	GetClientTimeout() int
}

// LogConfiguration ...
type LogConfiguration struct {
	Logfile string `json:"logfile"`
}

// ILogConfiguration ...
type ILogConfiguration interface {
	GetLogfile() string
}

// MetricsConfiguration ...
type MetricsConfiguration struct {
	NoMetrics bool `json:"nometrics"`
}

// IMetricsConfiguration ...
type IMetricsConfiguration interface {
	GetNoMetrics() bool
}

// AuthConfiguration ...
type AuthConfiguration struct {
	Passwordfile string `json:"passwordfile"`
}

// ILogConfiguration ...
type IAuthConfiguration interface {
	GetPasswordfile() string
}

type HeaderList []string

func (cs *HeaderList) String() string { return "" }
func (cs *HeaderList) Set(value string) error {
	*cs = append(*cs, value)
	return nil
}

// HeaderConfiguration ...
type HeaderConfiguration struct {
	RequestHeaderFunctions  [](RequestHeaderFunction)
	RequestHeaderStrings    []string
	RequestHeaders          []*Header
	ResponseHeaderFunctions [](ResponseHeaderFunction)
	ResponseHeaderStrings   []string
	ResponseHeaders         []*Header
	CopyHeaderStrings       []string
	CopyHeaderOperations    []*HeaderOperation
}

type IHeaderConfiguration interface {
	GetRequestHeaders() []*Header
	AddRequestHeaderFunction(fn RequestHeaderFunction)
	AddRequestHeader(h *Header)
	GetRequestHeaderFunctions() []RequestHeaderFunction
	AddResponseHeaderFunction(fn ResponseHeaderFunction)
	AddResponseHeader(h *Header)
	GetResponseHeaders() []*Header
	GetResponseHeaderFunctions() []ResponseHeaderFunction
	GetCopyHeaderOperations() []*HeaderOperation
	AddCopyHeaderOperation(op *HeaderOperation)
}

// Configuration ...
type Configuration struct {
	ConnectionConfiguration
	TLSConfiguration
	LimitConfiguration
	LogConfiguration
	MetricsConfiguration
	AuthConfiguration
	HeaderConfiguration
}

// IConfiguration
type IConfiguration interface {
	IConnectionConfiguration
	ITLSConfiguration
	ILimitConfiguration
	ILogConfiguration
	IMetricsConfiguration
	IAuthConfiguration
	IHeaderConfiguration
}

// // Debug ...
// func (cfg *ConnectionConfiguration) Debug(s string) { fmt.Printf("%s [%p]\n", s, &cfg) }

// GetHost ...
func (cfg *ConnectionConfiguration) GetHost() string { return cfg.Host }

// GetPort ...
func (cfg *ConnectionConfiguration) GetPort() int { return cfg.Port }

// GetNamespace ...
func (cfg *ConnectionConfiguration) GetNamespace() string { return cfg.Namespace }

// GetName ...
func (cfg ConnectionConfiguration) GetName() string { return cfg.Name }

// GetHostname ...
func (cfg ConnectionConfiguration) GetHostname() string { return cfg.Hostname }

// GetVersion ...
func (cfg ConnectionConfiguration) GetVersion() string { return cfg.Version }

// GetCertChainFile ...
func (cfg *TLSConfiguration) GetCertChainFile() string { return cfg.CertChainFile }

// GetKeyFile ...
func (cfg *TLSConfiguration) GetKeyFile() string { return cfg.KeyFile }

// GetCAFile ...
func (cfg *TLSConfiguration) GetCAFile() string { return cfg.CAFile }

// GetMaxConnections ...
func (cfg *LimitConfiguration) GetMaxConnections() int { return cfg.MaxConnections }

// GetDelayReply ...
func (cfg *LimitConfiguration) GetDelayReply() int { return cfg.DelayReply }

// GetDelayReply ...
func (cfg *LimitConfiguration) GetClientTimeout() int { return cfg.ClientTimeout }

// GetLogfile ...
func (cfg *LogConfiguration) GetLogfile() string { return cfg.Logfile }

// GetNoMetrics ...
func (cfg *MetricsConfiguration) GetNoMetrics() bool { return cfg.NoMetrics }

// GetPasswordfile ...
func (cfg *AuthConfiguration) GetPasswordfile() string { return cfg.Passwordfile }

// AddRequestHeaderFunction ...
func (cfg *HeaderConfiguration) AddRequestHeaderFunction(fn RequestHeaderFunction) {
	cfg.RequestHeaderFunctions = append(cfg.RequestHeaderFunctions, fn)
}

// AddResponseHeaderFunction ...
func (cfg *HeaderConfiguration) AddResponseHeaderFunction(fn ResponseHeaderFunction) {
	cfg.ResponseHeaderFunctions = append(cfg.ResponseHeaderFunctions, fn)
}

// GetResponseHeaders ...
func (cfg *HeaderConfiguration) GetResponseHeaders() []*Header {
	return cfg.ResponseHeaders
}

// GetResponseHeaderFunctions ...
func (cfg *HeaderConfiguration) GetResponseHeaderFunctions() []ResponseHeaderFunction {
	return cfg.ResponseHeaderFunctions
}

// GetRequestHeaders ...
func (cfg *HeaderConfiguration) GetRequestHeaders() []*Header {
	return cfg.RequestHeaders
}

// GetRequestHeaderFunctions ...
func (cfg *HeaderConfiguration) GetRequestHeaderFunctions() []RequestHeaderFunction {
	return cfg.RequestHeaderFunctions
}

// AddRequestHeader ...
func (cfg *HeaderConfiguration) AddRequestHeader(h *Header) {
	cfg.RequestHeaders = append(cfg.RequestHeaders, h)
}

// AddResponseHeader ...
func (cfg *HeaderConfiguration) AddResponseHeader(h *Header) {
	cfg.ResponseHeaders = append(cfg.ResponseHeaders, h)
}

// AddCopyHeaderOperation ...
func (cfg *HeaderConfiguration) AddCopyHeaderOperation(op *HeaderOperation) {
	cfg.CopyHeaderOperations = append(cfg.CopyHeaderOperations, op)
}

// GetCopyHeaderOperations ...
func (cfg *HeaderConfiguration) GetCopyHeaderOperations() []*HeaderOperation {
	return cfg.CopyHeaderOperations
}
