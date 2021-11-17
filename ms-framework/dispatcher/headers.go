package dispatcher

import (
	"net/http"
	"strings"
)

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type Header struct {
	Key   string
	Value string
}

func (h *Header) appendRequest(out *http.Request) {
	out.Header.Add(h.Key, h.Value)
}

func (h *Header) appendResponse(out http.ResponseWriter) {
	out.Header().Add(h.Key, h.Value)
}

func (h *Header) setRequest(out *http.Request) {
	out.Header.Set(h.Key, h.Value)
}

func (h *Header) setResponse(out http.ResponseWriter) {
	out.Header().Set(h.Key, h.Value)
}

// HeaderFromString ...
func HeaderFromString(line string) *Header {
	var header Header
	pairs := strings.SplitN(line, ":", 2)
	header.Key = pairs[0]
	header.Value = strings.Trim(pairs[1], " ")
	return &header
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

const (
	HO_Copy = iota
	HO_Append
	HO_Force
)

type HeaderOperation struct {
	Key string
	Op  int
}

func (h *HeaderOperation) setToResponse(out http.ResponseWriter, in *http.Request) {
	values := in.Header.Values(h.Key)

	if h.Op == HO_Force {
		out.Header().Del(h.Key)
	}

	if len(values) > 0 {

		if h.Op == HO_Copy {
			out.Header().Del(h.Key)
		}

		for _, value := range values {
			out.Header().Add(h.Key, value)
		}
	}
}

func (h *HeaderOperation) setToRequest(out *http.Request, in *http.Request) {
	values := in.Header.Values(h.Key)

	if h.Op == HO_Force {
		out.Header.Del(h.Key)
	}

	if len(values) > 0 {

		if h.Op == HO_Copy {
			out.Header.Del(h.Key)
		}

		for _, value := range values {
			out.Header.Add(h.Key, value)
		}
	}
}

var fixedCopyHeaders = []HeaderOperation{
	{Key: "x-environment", Op: HO_Copy},
	{Key: "x-session-id", Op: HO_Copy},
	{Key: "x-correlation-id", Op: HO_Copy},
	{Key: "x-sequence-nr", Op: HO_Copy},
	{Key: "x-debug-info", Op: HO_Copy},

	// istio/envoy tracing headers
	{Key: "x-request-id", Op: HO_Copy},
	{Key: "x-b3-traceid", Op: HO_Copy},
	{Key: "x-b3-spanid", Op: HO_Copy},
	{Key: "x-b3-parentspanid", Op: HO_Copy},
	{Key: "x-b3-sampled", Op: HO_Copy},
	{Key: "x-b3-flags", Op: HO_Copy},
	{Key: "b3", Op: HO_Copy},
	{Key: "x-ot-span-context", Op: HO_Copy},
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (ds *Dispatcher) SetResponseHeaders(contentType string, out http.ResponseWriter, in *http.Request) {

	if len(contentType) != 0 {
		out.Header().Set("Content-Type", contentType)
	}

	for _, fn := range ds.GetResponseHeaderFunctions() {
		fn(out, in)
	}

	if len(ds.GetNamespace()) > 0 {
		out.Header().Set("X-Namespace", ds.GetNamespace())
	}

	if in != nil {
		for _, h := range fixedCopyHeaders {
			h.setToResponse(out, in)
		}
		for _, h := range ds.GetCopyHeaderOperations() {
			h.setToResponse(out, in)
		}
	}

	if len(ds.GetResponseHeaders()) > 0 {
		for _, h := range ds.GetResponseHeaders() {
			h.setResponse(out)
		}
	}

}

func (ds *Dispatcher) SetRequestHeaders(contentType string, out *http.Request, in *http.Request) {

	if len(contentType) != 0 {
		out.Header.Set("Content-Type", contentType)
	}

	for _, fn := range ds.GetRequestHeaderFunctions() {
		fn(out, in)
	}

	if len(ds.GetNamespace()) > 0 {
		out.Header.Set("X-Namespace", ds.GetNamespace())
	}

	if in != nil {
		for _, h := range fixedCopyHeaders {
			h.setToRequest(out, in)
		}
		for _, h := range ds.GetCopyHeaderOperations() {
			h.setToRequest(out, in)
		}
	}

	if len(ds.GetRequestHeaders()) > 0 {
		for _, h := range ds.GetRequestHeaders() {
			h.appendRequest(out)
		}
	}
}
