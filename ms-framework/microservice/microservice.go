package microservice

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/com-gft-tsbo-source/go-common/ms-framework/dispatcher"
)

// ###########################################################################
// ###########################################################################
// MicroService
// ###########################################################################
// ###########################################################################

// MicroService ...
type MicroService struct {
	dispatcher.Dispatcher
	*DBConfiguration
	*ServiceConfiguration
	*FileConfiguration
	UserEntries map[string]UserEntry
}

// ---------------------------------------------------------------------------

// IMicroService is the general ms interface
type IMicroService interface {
	dispatcher.IDispatcher
	IConfiguration
	httpGetStatus(http.ResponseWriter, *http.Request) (int, int, string)
}

// ###########################################################################

// Init ...
func Init(ms *MicroService,
	configuration *Configuration,
	defaultHandler *dispatcher.HandlerGroup) {

	defaultRequestHeaderFn := func(out *http.Request, in *http.Request) {
		out.Header.Set("X-cid", ms.GetName())
		out.Header.Set("X-chost", ms.GetHostname())
		out.Header.Set("X-version", ms.GetVersion())
	}
	defaultResponseHeaderFn := func(out http.ResponseWriter, in *http.Request) {
		out.Header().Set("X-cid", ms.GetName())
		out.Header().Set("X-chost", ms.GetHostname())
		out.Header().Set("X-version", ms.GetVersion())
	}

	dispatcher.Init(&ms.Dispatcher, &configuration.Configuration, defaultHandler, nil, nil)
	ms.GetLogger().SetPrefix(fmt.Sprintf("[%-12.12s] ", configuration.GetName()))
	ms.DBConfiguration = &configuration.DBConfiguration
	ms.ServiceConfiguration = &configuration.ServiceConfiguration
	ms.AddRequestHeaderFunction(defaultRequestHeaderFn)
	ms.AddResponseHeaderFunction(defaultResponseHeaderFn)
	// ms.HeaderConfiguration = &configuration.HeaderConfiguration
	statusHandler := dispatcher.HandlerGroup{Get: ms.httpGetStatus}
	ms.UserEntries = nil

	if len(configuration.Passwordfile) > 0 {

		ms.UserEntries = UserEntriesFromFile(configuration.Passwordfile)

		checkUserAccessFn := func(w http.ResponseWriter, r *http.Request) bool {
			var entry UserEntry
			var exists bool

			username, password, ok := r.BasicAuth()

			if !ok {
				goto error
			}

			entry, exists = (ms.UserEntries)[username]

			if !exists {
				goto error
			}

			if !entry.CheckPassword(password) {
				goto error
			}

			return true
		error:
			ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
			ms.PageNotAuthorized(w, r)
			return false
		}

		ms.AddWrapper(checkUserAccessFn)
	}

	ms.AddHandler("/status", &statusHandler)
}

// ---------------------------------------------------------------------------

// InitFromArgs ...
func InitFromArgs(ms *MicroService, args []string, flagset *flag.FlagSet, defaultHandler *dispatcher.HandlerGroup) {
	var configuration Configuration
	InitConfigurationFromArgs(&configuration, args, flagset)
	Init(ms, &configuration, defaultHandler)

}

// ###########################################################################

// httpGetStatus ...
func (ms *MicroService) httpGetStatus(w http.ResponseWriter, r *http.Request) (status int, contentLen int, msg string) {
	var response Response
	InitResponseFromMicroService(&response, ms, "OK")
	status = http.StatusOK
	ms.SetResponseHeaders("application/json; charset=utf-8", w, r)
	w.WriteHeader(status)
	contentLen = ms.Reply(w, response)
	return status, contentLen, "Status is good."
}

// ---------------------------------------------------------------------------

// GetEndpoint ...
func (ms *MicroService) GetEndpoint(name string) string {
	if strings.HasPrefix(name, "/") {
		return fmt.Sprintf("%s%s", ms.GetBaseURL(), name)
	} else {
		return fmt.Sprintf("%s/%s", ms.GetBaseURL(), name)
	}
}

// ---------------------------------------------------------------------------

// Run ...
func (ms *MicroService) Run() {
	ms.GetLogger().Println(fmt.Sprintf("Starting MS '%s' at version '%s'.", ms.GetName(), ms.GetVersion()))
	ms.Dispatcher.Run()
}
