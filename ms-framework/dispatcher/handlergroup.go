package dispatcher

import (
	"net/http"
)

// ###########################################################################
// ###########################################################################
// Dispatcher HandlerGroup
// ###########################################################################
// ###########################################################################

// HTTPHandler is the function type for all HTTP handlers
type HTTPHandler = func(http.ResponseWriter, *http.Request) (int, int, string)

// HandlerGroup encapsulates the possible http handlers
type HandlerGroup struct {
	Get     HTTPHandler
	Put     HTTPHandler
	Post    HTTPHandler
	Delete  HTTPHandler
	Head    HTTPHandler
	Connect HTTPHandler
	Options HTTPHandler
	Any     HTTPHandler
}
