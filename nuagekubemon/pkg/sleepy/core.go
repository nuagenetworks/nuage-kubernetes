package sleepy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	HEAD   = "HEAD"
	PATCH  = "PATCH"
)

// GetSupported is the interface that provides the Get
// method a resource must support to receive HTTP GETs.
type GetSupported interface {
	Get(map[string]string, url.Values, http.Header) (
		int, interface{}, http.Header)
}

// PostSupported is the interface that provides the Post
// method a resource must support to receive HTTP POSTs.
type PostSupported interface {
	Post(map[string]string, url.Values, http.Header, map[string]interface{}) (
		int, interface{}, http.Header)
}

// PutSupported is the interface that provides the Put
// method a resource must support to receive HTTP PUTs.
type PutSupported interface {
	Put(map[string]string, url.Values, http.Header, map[string]interface{}) (
		int, interface{}, http.Header)
}

// DeleteSupported is the interface that provides the Delete
// method a resource must support to receive HTTP DELETEs.
type DeleteSupported interface {
	Delete(map[string]string, url.Values, http.Header) (
		int, interface{}, http.Header)
}

// HeadSupported is the interface that provides the Head
// method a resource must support to receive HTTP HEADs.
type HeadSupported interface {
	Head(map[string]string, url.Values, http.Header) (
		int, interface{}, http.Header)
}

// PatchSupported is the interface that provides the Patch
// method a resource must support to receive HTTP PATCHs.
type PatchSupported interface {
	Patch(map[string]string, url.Values, http.Header) (
		int, interface{}, http.Header)
}

// An API manages a group of resources by routing requests
// to the correct method on a matching resource and marshalling
// the returned data to JSON for the HTTP response.
//
// You can instantiate multiple APIs on separate ports. Each API
// will manage its own set of resources.
type API struct {
	router            *mux.Router
	routerInitialized bool
}

// NewAPI allocates and returns a new API.
func NewAPI() *API {
	return &API{}
}

func (api *API) requestHandler(resource interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {

		if request.ParseForm() != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		var code int
		var data interface{}
		var header http.Header
		vars := mux.Vars(request)

		switch request.Method {
		case GET:
			if resource, ok := resource.(GetSupported); ok {
				code, data, header = resource.Get(vars, request.Form, request.Header)
			}
		case POST:
			if resource, ok := resource.(PostSupported); ok {
				if request.ContentLength == 0 {
					// If the content length is 0, just pass in an empty map
					// because json.Unmarshal() will fail if there's no json.
					code, data, header = resource.Post(vars, request.Form,
						request.Header, make(map[string]interface{}))
				} else {
					bodyMap := make(map[string]interface{})
					body := make([]byte, request.ContentLength)
					if _, err := request.Body.Read(body); err != nil && err.Error() != "EOF" {
						code, data, header = http.StatusInternalServerError, nil, nil
						break
					}
					if err := json.Unmarshal(body, &bodyMap); err != nil {
						code, data, header = http.StatusBadRequest, map[string]string{
							"error": err.Error(),
						}, nil
						break
					}
					code, data, header = resource.Post(vars, request.Form,
						request.Header, bodyMap)
				}
			}
		case PUT:
			if resource, ok := resource.(PutSupported); ok {
				if request.ContentLength == 0 {
					// If the content length is 0, just pass in an empty map
					// because json.Unmarshal() will fail if there's no json.
					code, data, header = resource.Put(vars, request.Form,
						request.Header, make(map[string]interface{}))
				} else {
					bodyMap := make(map[string]interface{})
					body := make([]byte, request.ContentLength)
					if _, err := request.Body.Read(body); err != nil && err.Error() != "EOF" {
						code, data, header = http.StatusInternalServerError, nil, nil
						break
					}
					if err := json.Unmarshal(body, &bodyMap); err != nil {
						code, data, header = http.StatusBadRequest, map[string]string{
							"error": err.Error(),
						}, nil
						break
					}
					code, data, header = resource.Put(vars, request.Form,
						request.Header, bodyMap)
				}
			}
		case DELETE:
			if resource, ok := resource.(DeleteSupported); ok {
				code, data, header = resource.Delete(vars, request.Form, request.Header)
			}
		case HEAD:
			if resource, ok := resource.(HeadSupported); ok {
				code, data, header = resource.Head(vars, request.Form, request.Header)
			}
		case PATCH:
			if resource, ok := resource.(PatchSupported); ok {
				code, data, header = resource.Patch(vars, request.Form, request.Header)
			}
		default:
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var content []byte
		if data != nil {
			var err error
			content, err = json.MarshalIndent(data, "", "  ")
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		for name, values := range header {
			for _, value := range values {
				rw.Header().Add(name, value)
			}
		}
		rw.WriteHeader(code)
		if data != nil {
			rw.Write(content)
		}
	}
}

// Mux returns the http.ServeMux used by an API. If a ServeMux has
// does not yet exist, a new one will be created and returned.
func (api *API) Mux() *mux.Router {
	if api.routerInitialized {
		return api.router
	} else {
		api.router = mux.NewRouter()
		api.routerInitialized = true
		return api.router
	}
}

// AddResource adds a new resource to an API. The API will route
// requests that match one of the given paths to the matching HTTP
// method on the resource.
func (api *API) AddResource(resource interface{}, paths ...string) {
	for _, path := range paths {
		api.Mux().HandleFunc(path, api.requestHandler(resource))
	}
}

// AddResourceWithWrapper behaves exactly like AddResource but wraps
// the generated handler function with a give wrapper function to allow
// to hook in Gzip support and similar.
func (api *API) AddResourceWithWrapper(resource interface{}, wrapper func(handler http.HandlerFunc) http.HandlerFunc, paths ...string) {
	for _, path := range paths {
		api.Mux().HandleFunc(path, wrapper(api.requestHandler(resource)))
	}
}

// Start causes the API to begin serving requests on the given port.
func (api *API) Start(port int) error {
	if !api.routerInitialized {
		return errors.New("You must add at least one resource to this API.")
	}
	portString := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(portString, api.Mux())
}
