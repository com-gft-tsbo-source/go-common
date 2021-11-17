package dispatcher

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/netutil"
)

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var _build_customer string = "<???_build_customer???>"
var _build_module string = "<???_build_module???>"
var _build_component string = "<???_build_component???>"
var _build_project string = "<???_build_project???>"
var _build_stamp string = "<???_build_stamp???>"
var _build_commit string = "<???_build_commit???>"

// ###########################################################################
// ###########################################################################
// Dispatcher
// ###########################################################################
// ###########################################################################

type WrapperFunc func(w http.ResponseWriter, r *http.Request) bool

// Dispatcher encapsulates the data of a general REST dispatcher
type Dispatcher struct {
	IConfiguration
	muxer          *http.ServeMux
	defaultHandler *HandlerGroup
	logger         *log.Logger
	maxPathLen     int
	tlsInfo        *TLSInfo
	HTTPClient     *http.Client

	wrappers []WrapperFunc

	RequestHeaders  []Header
	ResponseHeaders []Header
	CopyHeaders     []HeaderOperation

	prometheusOps       prometheus.Counter
	prometheusOpsFailed prometheus.Counter
	prometheusOps404    prometheus.Counter
	prometheusOps401    prometheus.Counter
}

// ---------------------------------------------------------------------------

// IDispatcher is the interface of a general REST dispatcher
type IDispatcher interface {
	Run()
	AddHandler(string, *HandlerGroup)
	AddHandlerRaw(string, *HandlerGroup, string)
	Reply(http.ResponseWriter, interface{}) int
	ReplyData(http.ResponseWriter, interface{}, []byte) int
	DefaultHandler() *HandlerGroup
	GetLogger() *log.Logger
	GetBaseURL() string
}

// ###########################################################################

// InitDispatcher is the constructor of a Dispatcher
func Init(
	ds *Dispatcher,
	configuration *Configuration,
	defaultHandler *HandlerGroup,
	muxer *http.ServeMux,
	logger *log.Logger) {
	ds.IConfiguration = configuration

	if logger == nil {
		logFormat := fmt.Sprintf("[%-12.12s] ", "dispatcher")
		logFlags := log.Ldate | log.Ltime | log.LUTC | log.Lmsgprefix
		if len(configuration.GetLogfile()) == 0 {
			logger = log.New(os.Stdout, logFormat, logFlags)
		} else if configuration.GetLogfile() == "-" {
			logger = log.New(os.Stdout, logFormat, logFlags)
		} else if configuration.GetLogfile() == ":stdout" {
			logger = log.New(os.Stdout, logFormat, logFlags)
		} else if configuration.GetLogfile() == ":stderr" {
			logger = log.New(os.Stderr, logFormat, logFlags)
		} else {
			file, err := os.OpenFile(configuration.GetLogfile(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
			if err != nil {
				log.Panic(fmt.Printf("Could not open logfile '%s', error was '%s'!", configuration.GetLogfile(), err.Error()))
			}
			logger = log.New(&logFlushWriter{Writer: file}, logFormat, logFlags)
		}
	}

	if muxer == nil {
		muxer = http.NewServeMux()
	}

	ds.muxer = muxer
	ds.logger = logger
	ds.maxPathLen = 10
	ds.defaultHandler = defaultHandler

	// ds.HeaderConfiguration = &configuration.HeaderConfiguration
	for _, line := range configuration.RequestHeaderStrings {
		h := HeaderFromString(line)
		ds.AddRequestHeader(h)
	}

	for _, line := range configuration.ResponseHeaderStrings {
		h := HeaderFromString(line)
		ds.AddResponseHeader(h)
	}

	for _, line := range configuration.CopyHeaderStrings {
		ds.AddCopyHeaderOperation(&HeaderOperation{Key: line, Op: HO_Copy})
	}

	ds.prometheusOps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ops_total",
		Help: "The total number of processed events",
	})

	ds.prometheusOpsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ops_failed",
		Help: "The total number of failed processed events",
	})

	ds.prometheusOps404 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ops_not_found",
		Help: "The total number of not found events",
	})

	ds.prometheusOps401 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ops_not_authorized",
		Help: "The total number of not authorized events",
	})

	prometheusLogFn := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ds.GetLogger().Printf("> %-6.6s | %3d | %6d | /metrics   | %s\n", r.Method, 0, 0, "Prometheus metrics served.")
			h.ServeHTTP(w, r)
		}
	}

	if !ds.GetNoMetrics() {
		ds.Handle("/metrics", prometheusLogFn(promhttp.Handler()))
	}

	// ds.Handle("/metrics", promhttp.Handler())

	if defaultHandler != nil {
		ds.AddHandlerRaw("/", ds.defaultHandler, "")

		if len(ds.GetNamespace()) > 0 {
			ds.AddHandlerRaw("/", ds.defaultHandler, ds.GetNamespace())
		}
	} else {
		ds.defaultHandler = ds.DefaultHandler()
		ds.AddHandlerRaw("/", ds.defaultHandler, "")
	}

	if len(ds.GetCertChainFile()) > 0 && len(ds.GetKeyFile()) > 0 {

		if ds.tlsInfo == nil {
			ds.tlsInfo = &TLSInfo{}
		}

		certificate, err := tls.LoadX509KeyPair(ds.GetCertChainFile(), ds.GetKeyFile())
		if err != nil {
			ds.GetLogger().Fatal(err)
		}
		ds.tlsInfo.certificate = &certificate
	}

	if len(ds.GetCAFile()) > 0 {

		if ds.tlsInfo == nil {
			ds.tlsInfo = &TLSInfo{}
		}

		caCert, err := ioutil.ReadFile(ds.GetCAFile())
		if err != nil {
			ds.GetLogger().Fatal(err)
		}
		ds.tlsInfo.caCertPool = x509.NewCertPool()
		ds.tlsInfo.caCertPool.AppendCertsFromPEM(caCert)

	}

	if ds.tlsInfo == nil {

		ds.HTTPClient = &http.Client{Timeout: time.Duration(ds.GetClientTimeout()) * time.Millisecond}
		return

	}

	if ds.tlsInfo.certificate != nil {

		ds.tlsInfo.tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{*ds.tlsInfo.certificate},
			RootCAs:      ds.tlsInfo.caCertPool,
		}

	} else {

		ds.tlsInfo.tlsConfig = &tls.Config{
			RootCAs: ds.tlsInfo.caCertPool,
		}

	}

	ds.tlsInfo.tlsConfig.VerifyPeerCertificate = checkCertificate
	ds.tlsInfo.tlsConfig.ClientAuth = tls.RequireAnyClientCert
	ds.tlsInfo.tlsConfig.ClientCAs = ds.tlsInfo.caCertPool
	ds.tlsInfo.transport = &http.Transport{TLSClientConfig: ds.tlsInfo.tlsConfig}
	ds.HTTPClient = &http.Client{Transport: ds.tlsInfo.transport, Timeout: time.Duration(ds.GetClientTimeout()) * time.Millisecond}

}

// ---------------------------------------------------------------------------

// HasTLS ...
func (ds *Dispatcher) HasTLS() bool {
	return ds.tlsInfo != nil
}

// ---------------------------------------------------------------------------

// GetBaseURL ...
func (ds *Dispatcher) GetBaseURL() string {
	var url string

	if ds.tlsInfo != nil && ds.tlsInfo.certificate != nil {
		url = fmt.Sprintf("https://%s:%d", ds.GetHost(), ds.GetPort())
	} else {
		url = fmt.Sprintf("http://%s:%d", ds.GetHost(), ds.GetPort())
	}

	if len(ds.GetNamespace()) > 0 {
		url = fmt.Sprintf("%s/%s", url, ds.GetNamespace())
	}

	return url
}

// ---------------------------------------------------------------------------

// GetLogger ...
func (ds *Dispatcher) GetLogger() *log.Logger { return ds.logger }

// ---------------------------------------------------------------------------

// DefaultHandler ...
func (ds *Dispatcher) DefaultHandler() *HandlerGroup {
	return &HandlerGroup{Any: ds.PageNotFound, Options: ds.defaultOptions}
}

// ---------------------------------------------------------------------------

// Run is the main entry point for the dispatcher
func (ds *Dispatcher) Run() {

	addAccessControlAllowOriginFn := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		}
	}

	delayReplyFn := func(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Duration(ds.GetDelayReply()) * time.Millisecond)
			h.ServeHTTP(w, r)
		}
	}

	addWrapper := func(h http.Handler, wrapper WrapperFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			success := wrapper(w, r)

			if !success {
				return
			}

			h.ServeHTTP(w, r)
		}
	}

	var err error
	var listener net.Listener
	var wrappedHandler http.HandlerFunc

	ds.GetLogger().Println(fmt.Sprintf("This is '%s' in module '%s' for project '%s' of customer '%s' built at '%s' from '%s'.", _build_component, _build_module, _build_project, _build_customer, _build_stamp, _build_commit))

	if ds.tlsInfo != nil && ds.tlsInfo.certificate != nil {
		listener, err = tls.Listen("tcp", fmt.Sprintf("%s:%d", ds.GetHost(), ds.GetPort()), ds.tlsInfo.tlsConfig)
		if err != nil {
			ds.GetLogger().Fatal(err)
			return
		}
	} else {
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", ds.GetHost(), ds.GetPort()))
		if err != nil {
			ds.GetLogger().Fatal(err)
			return
		}
	}

	if ds.GetMaxConnections() > 0 {
		listener = netutil.LimitListener(listener, ds.GetMaxConnections())
	}

	wrappedHandler = addAccessControlAllowOriginFn(ds.muxer)

	for _, wrapper := range ds.wrappers {
		wrappedHandler = addWrapper(wrappedHandler, wrapper)
	}

	if ds.GetDelayReply() > 0 {
		wrappedHandler = delayReplyFn(wrappedHandler)
	}

	if ds.GetMaxConnections() > 0 {
		ds.GetLogger().Println(fmt.Sprintf("Allowing %d concurrent connections.", ds.GetMaxConnections()))
	}

	if ds.GetDelayReply() > 0 {
		ds.GetLogger().Println(fmt.Sprintf("Delaying replies by %dms..", ds.GetDelayReply()))
	}

	if ds.tlsInfo != nil && ds.tlsInfo.certificate != nil {
		ds.GetLogger().Println(fmt.Sprintf("Starting listener on 'https://%s:%d'", ds.GetHost(), ds.GetPort()))
	} else {
		ds.GetLogger().Println(fmt.Sprintf("Starting listener on 'http://%s:%d'", ds.GetHost(), ds.GetPort()))
	}

	err = http.Serve(listener, wrappedHandler)
	ds.GetLogger().Fatal(err)
}

// ---------------------------------------------------------------------------

// AddHandlerRaw adds a HTTP handler to the current dispatcher
func (ds *Dispatcher) AddHandlerRaw(path string, handlers *HandlerGroup, namespace string) {
	if handlers.Any == nil {
		handlers.Any = ds.PageNotFound
	}
	if handlers.Put == nil {
		handlers.Put = ds.PageNotFound
	}
	if handlers.Get == nil {
		handlers.Get = ds.PageNotFound
	}
	if handlers.Post == nil {
		handlers.Post = ds.PageNotFound
	}
	if handlers.Delete == nil {
		handlers.Delete = ds.PageNotFound
	}
	if handlers.Head == nil {
		handlers.Head = ds.PageNotFound
	}
	if handlers.Connect == nil {
		handlers.Connect = ds.PageNotFound
	}
	if handlers.Options == nil {
		handlers.Options = ds.PageNotFound
	}

	if len(namespace) > 0 {
		path = fmt.Sprintf("/%s%s", namespace, path)
	}
	l := len(path)

	if ds.maxPathLen < l {
		ds.maxPathLen = l
	}

	ds.muxer.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) { ds.handler(handlers, w, r) })
}

// AddHandler adds a HTTP handler to the current dispatcher
func (ds *Dispatcher) AddHandler(path string, handlers *HandlerGroup) {
	ds.AddHandlerRaw(path, handlers, ds.GetNamespace())
}

// ---------------------------------------------------------------------------

// Handle ...
func (ds *Dispatcher) Handle(path string, handler http.Handler) {
	if len(ds.GetNamespace()) > 0 {
		path = fmt.Sprintf("/%s%s", ds.GetNamespace(), path)
	}
	ds.muxer.Handle(path, http.StripPrefix(path, handler))
}

// ---------------------------------------------------------------------------

// Reply writes a reply to a request
func (ds *Dispatcher) Reply(w http.ResponseWriter, msg interface{}) int {
	msgBytes, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0
	}
	w.Write(msgBytes)
	return len(msgBytes)
}

// ---------------------------------------------------------------------------

// ReplyData writes a reply with additional data
func (ds *Dispatcher) ReplyData(w http.ResponseWriter, msg interface{}, payload []byte) int {
	msgBytes, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0
	}
	w.Write(msgBytes)
	if len(payload) > 0 {
		w.Write(payload)
	}
	return len(msgBytes) + len(payload)
}

// ---------------------------------------------------------------------------

func (ds *Dispatcher) handler(handlers *HandlerGroup, w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {

	ds.prometheusOps.Inc()
	fps := r.Header.Get("X-FailurePercent")

	if len(fps) > 0 {
		fpi, _ := strconv.Atoi(fps)

		if fpi > seededRand.Intn(100) {

			fcs := r.Header.Get("X-FailureCode")
			var fci int = 418

			if len(fcs) > 0 {
				fci, _ = strconv.Atoi(fcs)
			}

			ds.prometheusOpsFailed.Inc()
			msg = fmt.Sprintf("Forcing error %d.", fci)
			ds.GetLogger().Printf("> %-6.6s | %3d | %6d | %-*.*s | %s\n", r.Method, fci, len(msg), ds.maxPathLen, ds.maxPathLen, r.URL.Path, msg)
			http.Error(w, msg, fci)
			return fci, 0, ""
		}
	}

	switch r.Method {
	case http.MethodGet:
		status, contentLen, msg = handlers.Get(w, r)
	case http.MethodPut:
		status, contentLen, msg = handlers.Put(w, r)
	case http.MethodPost:
		status, contentLen, msg = handlers.Post(w, r)
	case http.MethodDelete:
		status, contentLen, msg = handlers.Delete(w, r)
	case http.MethodHead:
		status, contentLen, msg = handlers.Head(w, r)
	case http.MethodConnect:
		status, contentLen, msg = handlers.Connect(w, r)
	case http.MethodOptions:
		status, contentLen, msg = handlers.Options(w, r)
	default:
		status, contentLen, msg = handlers.Any(w, r)
	}

	ds.GetLogger().Printf("> %-6.6s | %3d | %6d | %-*.*s | %s\n", r.Method, status, contentLen, ds.maxPathLen, ds.maxPathLen, r.URL.Path, msg)
	return status, contentLen, msg
}

// ---------------------------------------------------------------------------

func (ds *Dispatcher) PageNotFound(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	var response Response
	ds.prometheusOps404.Inc()
	InitResponseFromDispatcher(&response, ds, "Error")
	status = http.StatusNotFound
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	contentLen = ds.Reply(w, response)
	return status, contentLen, "Path not registered"
}

// ---------------------------------------------------------------------------

func (ds *Dispatcher) PageNotAuthorized(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	var response Response
	ds.prometheusOps401.Inc()
	InitResponseFromDispatcher(&response, ds, "Not authorized")
	status = http.StatusUnauthorized
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	contentLen = ds.Reply(w, response)
	return status, contentLen, "Not authorized"
}

// ---------------------------------------------------------------------------

func (ds *Dispatcher) defaultOptions(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	var response Response
	InitResponseFromDispatcher(&response, ds, "Allow all")
	status = http.StatusOK
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Headers", "X-Cid, X-Chost, X-Version, X-Namespace, X-Environment, x-session-id, x-correlation-id, x-sequence-nr, x-request-id, x-b3-traceid, x-b3-spanid, x-b3-parentspanid, x-b3-sampled, x-b3-flags, b3, x-ot-span-context")
	w.WriteHeader(status)
	contentLen = ds.Reply(w, response)
	return status, contentLen, "Sent default options."
}

// ---------------------------------------------------------------------------

func (ds *Dispatcher) AddWrapper(wrapper WrapperFunc) {
	ds.wrappers = append(ds.wrappers, wrapper)
}
