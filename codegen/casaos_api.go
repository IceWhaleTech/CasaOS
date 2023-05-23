// Package codegen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package codegen

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	Access_tokenScopes = "access_token.Scopes"
)

// Defines values for SetZerotierNetworkStatusJSONBodyStatus.
const (
	Offline SetZerotierNetworkStatusJSONBodyStatus = "offline"
	Online  SetZerotierNetworkStatusJSONBodyStatus = "online"
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

// ZTInfo defines model for ZTInfo.
type ZTInfo struct {
	Id     *string `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
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

// GetZTInfoOK defines model for GetZTInfoOK.
type GetZTInfoOK = ZTInfo

// ResponseInternalServerError defines model for ResponseInternalServerError.
type ResponseInternalServerError = BaseResponse

// ResponseOK defines model for ResponseOK.
type ResponseOK = BaseResponse

// SetZerotierNetworkStatusJSONBody defines parameters for SetZerotierNetworkStatus.
type SetZerotierNetworkStatusJSONBody struct {
	Status *SetZerotierNetworkStatusJSONBodyStatus `json:"status,omitempty"`
}

// SetZerotierNetworkStatusJSONBodyStatus defines parameters for SetZerotierNetworkStatus.
type SetZerotierNetworkStatusJSONBodyStatus string

// SetZerotierNetworkStatusJSONRequestBody defines body for SetZerotierNetworkStatus for application/json ContentType.
type SetZerotierNetworkStatusJSONRequestBody SetZerotierNetworkStatusJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Test file methods
	// (GET /file/test)
	GetFileTest(ctx echo.Context) error
	// Get log
	// (GET /health/logs)
	GetHealthlogs(ctx echo.Context) error
	// Get port in use
	// (GET /health/ports)
	GetHealthPorts(ctx echo.Context) error
	// Get service status
	// (GET /health/services)
	GetHealthServices(ctx echo.Context) error
	// Get Zerotier info
	// (GET /zt/info)
	GetZerotierInfo(ctx echo.Context) error
	// Set Zerotier network status
	// (PUT /zt/{network_id}/status)
	SetZerotierNetworkStatus(ctx echo.Context, networkId string) error
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

// GetHealthlogs converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealthlogs(ctx echo.Context) error {
	var err error

	ctx.Set(Access_tokenScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHealthlogs(ctx)
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

// GetZerotierInfo converts echo context to params.
func (w *ServerInterfaceWrapper) GetZerotierInfo(ctx echo.Context) error {
	var err error

	ctx.Set(Access_tokenScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetZerotierInfo(ctx)
	return err
}

// SetZerotierNetworkStatus converts echo context to params.
func (w *ServerInterfaceWrapper) SetZerotierNetworkStatus(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "network_id" -------------
	var networkId string

	err = runtime.BindStyledParameterWithLocation("simple", false, "network_id", runtime.ParamLocationPath, ctx.Param("network_id"), &networkId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter network_id: %s", err))
	}

	ctx.Set(Access_tokenScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SetZerotierNetworkStatus(ctx, networkId)
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
	router.GET(baseURL+"/health/logs", wrapper.GetHealthlogs)
	router.GET(baseURL+"/health/ports", wrapper.GetHealthPorts)
	router.GET(baseURL+"/health/services", wrapper.GetHealthServices)
	router.GET(baseURL+"/zt/info", wrapper.GetZerotierInfo)
	router.PUT(baseURL+"/zt/:network_id/status", wrapper.SetZerotierNetworkStatus)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xX72/bNhD9VwhuH5pBsbxk3ToB/dAfaxoUW4Ilw4bGhktLZ4mNRKp3pyRuoP99oCjb",
	"sq386pqhn2xRp7v3Hnnk47WMbVFaA4ZJRtcSgUprCJqHA+C3oHLOji0yHb1zY7E1DIbdX1WWuY4Va2vC",
	"j2SNG6M4g0I1b/P8aCajs2v5PcJMRvK7cFUq9HEUvlQEf7Y1ZR1cyxJtCcjaI0gUN8luS9GBKOu6Htd1",
	"HcgEKEZdOmwykkfvZB2s6JwAXugYvnFGC5R3knp/emhm9oFkbqvvE8obay7oHRoGNCp3SAF/Q7T41TCs",
	"y7iNZFFb+OLCV++A+4p63IXFqVIHbbJmlte+iDbXQAFEKm1erCdqXwgErtBAIqZzQZ4f6QSEngnOAEFo",
	"EsrMZSDhShVlDjKSMpAIKjky+VxGjBUEkuele0OM2qQeeLdZtnBxXLofzVA0z8vkz4bLZNowpNAo3Y4o",
	"ROWgXO2mdteowo2dvjp2EVVyQ8Kn+w9M+Nfr4y6BZW9scTCWJ1gZ4xj3lpaxImVpQD6F3JJpA0cdyPvk",
	"200Vw6Wa3z+vY9N22hYLnazX+HFv/6enP//y7NdhX16vUTf+lSJ1dNIXS6y42mBgTa5ND+IGIkFcoeb5",
	"iVvdHp2KYyCasD2Hpou0W7sZqARQLtDIFxVnFvXnpuFWuVWp30HLXrfc15tgVA2H+3GpY64QmgcYGSGE",
	"8C/IVhiDKCDR6vlIPikRZoC0G9vc4m7TgxCJROH5zkgKwpiAn49kxlxSFIaoLgep5qyaVgTYbg+D2Bbh",
	"YQx/ZyqHU4izMLepDQulTeint/2ZTJUxgBOXfmJ0mvHk2XBYXg1Kk47kl4LNXaJHRMuXuikxmeYV3A5Y",
	"F6lQuYPgl5AH9f8j8mjCjVUwMh6VeHF8KEq0FzoBEoWmGPJcGbAViQI4swmJmUWR6NkMEAwLisEo1JYG",
	"Lssbi0ITVeC20UQkmuKKSFtDgShzUATiQpNmt9uKswPNb6upQCgtabY4Hz9ZqOGV2KbvYe4Ii+Kj1Uac",
	"2QrFa02xxWT1deIHBmkanptPL6bTl1P4Z2cwatpFc6eTHWEZyAtA8k1yseea2ZZgVKllJPcHw8G+DGSp",
	"OGt6NJzpHEIGas6+FHi70U6BWLiwhWYD2aTEpmUPExk5d/FGO07EzfnSMYd7w+FN5+YyLuwcxnUgnz7k",
	"kz5z0exHVVEonPfhd7KplGR0Jt90h8fuuzBrTg63MqkjyRZff8A0Uf2Mb7AUNmbgXWIEVaxbi5nFQrGM",
	"5FQbh7z3YO7zWV9ZrwNgkdu0o5Ln2q9TufAItwvlrcSXrI2t28XjMHY8hDaiIrgnc+qYi97GcWlbTyD8",
	"eSrsTICKM/GhdQM/fBBtmt6m2nAx/0m+zm3mcRRsibRU7xLxM4eLc/1G8d4DWtaAwkX2CrSIaOzRF8qz",
	"vBc9jixrJDqqLMe3dLk2wJcWzyc6qcOVESurHplOuhXa79oZ2BbsZCXYHz70ZDFXpUJVAANSc3NdL7LI",
	"qxMZeBfnzo+Vh1vhbebgU6URksXtYrXBbW5nYx8MxC9tMn/QNWzdBnfMqqkKJ+7SqtrZrPk3Du5jZJcj",
	"dvoRYm6t7be2qm6Z8zvWV8enN9O87tDPxm5K/EXSL4MKcxnJ8GKv9WXSBbQFNhfJk9Oj10c7q0Wx0fV1",
	"cNcHa2exK3S1yyo9QFuVvl4b9/vWKb61wYzrfwMAAP//9tkexLESAAA=",
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
