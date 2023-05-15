// Package codegen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package codegen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	Access_tokenScopes = "access_token.Scopes"
)

// BaseResponse defines model for BaseResponse.
type BaseResponse struct {
	// Message message returned by server side if there is any
	Message *string `json:"message,omitempty"`
}

// HealthPorts defines model for HealthPorts.
type HealthPorts struct {
	TCP *[]int `json:"tcp,omitempty"`
	UDP *[]int `json:"udp,omitempty"`
}

// HealthServices defines model for HealthServices.
type HealthServices struct {
	NotRunning *[]string `json:"not_running,omitempty"`
	Running    *[]string `json:"running,omitempty"`
}

// GetHealthPortsOK defines model for GetHealthPortsOK.
type GetHealthPortsOK struct {
	Data *HealthPorts `json:"data,omitempty"`

	// Message message returned by server side if there is any
	Message *string `json:"message,omitempty"`
}

// GetHealthServicesOK defines model for GetHealthServicesOK.
type GetHealthServicesOK struct {
	Data *HealthServices `json:"data,omitempty"`

	// Message message returned by server side if there is any
	Message *string `json:"message,omitempty"`
}

// ResponseInternalServerError defines model for ResponseInternalServerError.
type ResponseInternalServerError = BaseResponse

// ResponseOK defines model for ResponseOK.
type ResponseOK = BaseResponse

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Test file methods
	// (GET /file/test)
	GetFileTest(ctx echo.Context) error
	// Get port in use
	// (GET /health/ports)
	GetHealthPorts(ctx echo.Context) error
	// Get service status
	// (GET /health/services)
	GetHealthServices(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetFileTest converts echo context to params.
func (w *ServerInterfaceWrapper) GetFileTest(ctx echo.Context) error {
	var err error

	ctx.Set(Access_tokenScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetFileTest(ctx)
	return err
}

// GetHealthPorts converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealthPorts(ctx echo.Context) error {
	var err error

	ctx.Set(Access_tokenScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHealthPorts(ctx)
	return err
}

// GetHealthServices converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealthServices(ctx echo.Context) error {
	var err error

	ctx.Set(Access_tokenScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHealthServices(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/file/test", wrapper.GetFileTest)
	router.GET(baseURL+"/health/ports", wrapper.GetHealthPorts)
	router.GET(baseURL+"/health/services", wrapper.GetHealthServices)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xW4W/bthP9Vwj+fh+aQbaMBAUKAf2QNksaBIODJcMGxIZLU2eJjURyd6QTL9D/PpCS",
	"Y8V2g6Rdh31yRB7fvfdyx+MDl6a2RoN2xLMHjkDWaIL4cQbuE4jKlZcGHY0vwpo02oF24U9hbaWkcMro",
	"9AsZHdZIllCLuFtV4wXPbh74/xEWPOP/Szep0jaO0g+C4NcuJ2+SB27RWECnWga5cBHsOYgeRd40zbRp",
	"moTnQBKVDdx4xscXvEk2cq4Al0rCf1zRmuXzotapzrUD1KIKpwB/RjT4KnEvl7TLZJ2btclZm71H7pVG",
	"fw+X4EqTdGDR8Scnsu3/Rw1EoogbT4G6DYbgPGrI2XzFqNVHKgemFsyVgMAUMaFXPOFwL2pbAc84TziC",
	"yMe6WvHMoYeEu5UNO+RQ6aIl3i/cHV5O2vCjHNTx+xH83egRTGkHBUSnuxWBKAKV+0FhBlrUYe3642WI",
	"8PlXAN8evRLwt5PLvoDHOt3RoI2bodc6KN6bmktBwtCQWgi+Y9MWjybhL8EbFMLBnVi9HDfKIZAelVtd",
	"hdppFQgpgWjmzC3EGlWhMkoQOSBPeOfHsXelQfVXLOdNLmHVBaxap5RemN0Sm/jR6EhaJZ1HiB8w0Ywx",
	"1m6Q8SiB1ZAr8X7C31iEBSANpKkMDmKFQ8ZygbcHE84IJYF7P+Glc5ayNEVxNyyUK/3cE2DXfENp6vRc",
	"wu+lqOAaZJlWpjBpLZROW/O6n9lcaA04C/AzrYrSzd6NRvZ+aHUx4d9KtgpAP5Ctu1MxxWxeeXiesKoL",
	"JqpA4aMgMb5qSf37jFo26VYVTHTLih1fnjOLZqlyIFYrklBVQoPxxGpwpcmJLQyyXC0WgKAdIwlaoDI0",
	"DCinBpki8hAuqZzliqQnUkZTwmwFgoAtFSkX7jJ2c6bcJz9nCNaQcgZX0zdrN1onduW3NA+YQfbFKM1u",
	"jEd2okgazDen83ZhWBTprf7zeD7/MIc/DoaT2C7Kxd7dCOYJXwJS2yTLw9CuxoIWVvGMHw1HwyOecCtc",
	"GXs0XagKUgcUJ0sBbrfRroEcC2Frz4Y8QmJs2fOcZ+FxcKqCJnLx9u49gw5Ho69Npce4tDfqmoS/fc2R",
	"faM73ke+rgWu9vEPtomCeHbDT/vL03AuLeO9nNr1ZOk82RHcH0DfonnnffjPKz8Dx4IOpjTzBD3dbeb9",
	"yqk3kvYWRIDtJgkjJ5wnZhYMhCzZ526G/PSZdTB7i2Vr9n2Xfb336I9xsBPSSX3WxN4cjE/dpxPwZtpM",
	"Q0BIRnHfY8Uzni4Pu3uPh4AOftv1N9fjk/HBZnBuZQ+P5ecPPKn1kOh+4ERxhsbbNl8X98tOl+wInTZ/",
	"BwAA///xoj5w+wwAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
