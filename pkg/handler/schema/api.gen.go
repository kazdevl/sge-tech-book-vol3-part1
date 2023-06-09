// Package schema provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package schema

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// APIResponse defines model for APIResponse.
type APIResponse struct {
	Common   *CommonResponse   `json:"common,omitempty"`
	Original *OriginalResponse `json:"original,omitempty"`
}

// CommonResponse defines model for CommonResponse.
type CommonResponse struct {
	Delete *struct {
		UserCoin    *[]UserCoin    `json:"user_coin,omitempty"`
		UserItem    *[]UserItem    `json:"user_item,omitempty"`
		UserMonster *[]UserMonster `json:"user_monster,omitempty"`
	} `json:"delete,omitempty"`
	Update *struct {
		UserCoin    *[]UserCoin    `json:"user_coin,omitempty"`
		UserItem    *[]UserItem    `json:"user_item,omitempty"`
		UserMonster *[]UserMonster `json:"user_monster,omitempty"`
	} `json:"update,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// OriginalResponse defines model for OriginalResponse.
type OriginalResponse struct {
	UserRegister *UserRegisterResponseContent `json:"user_register,omitempty"`
}

// UserCoin defines model for UserCoin.
type UserCoin struct {
	Currency *int64 `json:"currency,omitempty"`
	UserId   *int64 `json:"user_id,omitempty"`
}

// UserItem defines model for UserItem.
type UserItem struct {
	Count  *int64 `json:"count,omitempty"`
	ItemId *int64 `json:"item_id,omitempty"`
	UserId *int64 `json:"user_id,omitempty"`
}

// UserMonster defines model for UserMonster.
type UserMonster struct {
	Exp       *int64 `json:"exp,omitempty"`
	MonsterId *int64 `json:"monster_id,omitempty"`
	UserId    *int64 `json:"user_id,omitempty"`
}

// UserRegisterResponseContent defines model for UserRegisterResponseContent.
type UserRegisterResponseContent struct {
	UserId *int64 `json:"user_id,omitempty"`
}

// MonsterEnhanceJSONBody defines parameters for MonsterEnhance.
type MonsterEnhanceJSONBody struct {
	Items *[]struct {
		Count  *int64 `json:"count,omitempty"`
		ItemId *int64 `json:"item_id,omitempty"`
	} `json:"items,omitempty"`
	MonsterId *int64 `json:"monster_id,omitempty"`
}

// MonsterEnhanceParams defines parameters for MonsterEnhance.
type MonsterEnhanceParams struct {
	UserId int64 `form:"userId" json:"userId"`
}

// UserGetDataParams defines parameters for UserGetData.
type UserGetDataParams struct {
	UserId int64 `form:"userId" json:"userId"`
}

// MonsterEnhanceJSONRequestBody defines body for MonsterEnhance for application/json ContentType.
type MonsterEnhanceJSONRequestBody MonsterEnhanceJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// MonsterEnhance request with any body
	MonsterEnhanceWithBody(ctx context.Context, params *MonsterEnhanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	MonsterEnhance(ctx context.Context, params *MonsterEnhanceParams, body MonsterEnhanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UserGetData request
	UserGetData(ctx context.Context, params *UserGetDataParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UserRegister request
	UserRegister(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) MonsterEnhanceWithBody(ctx context.Context, params *MonsterEnhanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewMonsterEnhanceRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) MonsterEnhance(ctx context.Context, params *MonsterEnhanceParams, body MonsterEnhanceJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewMonsterEnhanceRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UserGetData(ctx context.Context, params *UserGetDataParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUserGetDataRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UserRegister(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUserRegisterRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewMonsterEnhanceRequest calls the generic MonsterEnhance builder with application/json body
func NewMonsterEnhanceRequest(server string, params *MonsterEnhanceParams, body MonsterEnhanceJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewMonsterEnhanceRequestWithBody(server, params, "application/json", bodyReader)
}

// NewMonsterEnhanceRequestWithBody generates requests for MonsterEnhance with any type of body
func NewMonsterEnhanceRequestWithBody(server string, params *MonsterEnhanceParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/monster/enhance")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "userId", runtime.ParamLocationQuery, params.UserId); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewUserGetDataRequest generates requests for UserGetData
func NewUserGetDataRequest(server string, params *UserGetDataParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/user/data")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "userId", runtime.ParamLocationQuery, params.UserId); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewUserRegisterRequest generates requests for UserRegister
func NewUserRegisterRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/user/register")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// MonsterEnhance request with any body
	MonsterEnhanceWithBodyWithResponse(ctx context.Context, params *MonsterEnhanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*MonsterEnhanceResponse, error)

	MonsterEnhanceWithResponse(ctx context.Context, params *MonsterEnhanceParams, body MonsterEnhanceJSONRequestBody, reqEditors ...RequestEditorFn) (*MonsterEnhanceResponse, error)

	// UserGetData request
	UserGetDataWithResponse(ctx context.Context, params *UserGetDataParams, reqEditors ...RequestEditorFn) (*UserGetDataResponse, error)

	// UserRegister request
	UserRegisterWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*UserRegisterResponse, error)
}

type MonsterEnhanceResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *APIResponse
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r MonsterEnhanceResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r MonsterEnhanceResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UserGetDataResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *APIResponse
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r UserGetDataResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UserGetDataResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UserRegisterResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *APIResponse
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r UserRegisterResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UserRegisterResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// MonsterEnhanceWithBodyWithResponse request with arbitrary body returning *MonsterEnhanceResponse
func (c *ClientWithResponses) MonsterEnhanceWithBodyWithResponse(ctx context.Context, params *MonsterEnhanceParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*MonsterEnhanceResponse, error) {
	rsp, err := c.MonsterEnhanceWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseMonsterEnhanceResponse(rsp)
}

func (c *ClientWithResponses) MonsterEnhanceWithResponse(ctx context.Context, params *MonsterEnhanceParams, body MonsterEnhanceJSONRequestBody, reqEditors ...RequestEditorFn) (*MonsterEnhanceResponse, error) {
	rsp, err := c.MonsterEnhance(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseMonsterEnhanceResponse(rsp)
}

// UserGetDataWithResponse request returning *UserGetDataResponse
func (c *ClientWithResponses) UserGetDataWithResponse(ctx context.Context, params *UserGetDataParams, reqEditors ...RequestEditorFn) (*UserGetDataResponse, error) {
	rsp, err := c.UserGetData(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUserGetDataResponse(rsp)
}

// UserRegisterWithResponse request returning *UserRegisterResponse
func (c *ClientWithResponses) UserRegisterWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*UserRegisterResponse, error) {
	rsp, err := c.UserRegister(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUserRegisterResponse(rsp)
}

// ParseMonsterEnhanceResponse parses an HTTP response from a MonsterEnhanceWithResponse call
func ParseMonsterEnhanceResponse(rsp *http.Response) (*MonsterEnhanceResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &MonsterEnhanceResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest APIResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseUserGetDataResponse parses an HTTP response from a UserGetDataWithResponse call
func ParseUserGetDataResponse(rsp *http.Response) (*UserGetDataResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UserGetDataResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest APIResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseUserRegisterResponse parses an HTTP response from a UserRegisterWithResponse call
func ParseUserRegisterResponse(rsp *http.Response) (*UserRegisterResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UserRegisterResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest APIResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// モンスター強化
	// (POST /monster/enhance)
	MonsterEnhance(ctx echo.Context, params MonsterEnhanceParams) error
	// ユーザデータ取得
	// (GET /user/data)
	UserGetData(ctx echo.Context, params UserGetDataParams) error
	// ユーザ登録
	// (POST /user/register)
	UserRegister(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// MonsterEnhance converts echo context to params.
func (w *ServerInterfaceWrapper) MonsterEnhance(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params MonsterEnhanceParams
	// ------------- Required query parameter "userId" -------------

	err = runtime.BindQueryParameter("form", true, true, "userId", ctx.QueryParams(), &params.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.MonsterEnhance(ctx, params)
	return err
}

// UserGetData converts echo context to params.
func (w *ServerInterfaceWrapper) UserGetData(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params UserGetDataParams
	// ------------- Required query parameter "userId" -------------

	err = runtime.BindQueryParameter("form", true, true, "userId", ctx.QueryParams(), &params.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UserGetData(ctx, params)
	return err
}

// UserRegister converts echo context to params.
func (w *ServerInterfaceWrapper) UserRegister(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UserRegister(ctx)
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

	router.POST(baseURL+"/monster/enhance", wrapper.MonsterEnhance)
	router.GET(baseURL+"/user/data", wrapper.UserGetData)
	router.POST(baseURL+"/user/register", wrapper.UserRegister)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xWwWsTTxT+V368n8eQjVY87E1rkRyKUvRUioy7L8mU7Mz0zWxpKDmkAXup6KniSVCo",
	"B/GioBWL/WOWWPwvZGa22bTZ2I2UnnrK7O6bb7753vfeyzZEMlFSoDAawm3QUQcT5pZ3HzVXUCspNNpH",
	"RVIhGY7uYySTRAq7ukHYghD+DwqgIEcJFl3UGKVfA0m8zQXrXrTzYR5X7O3XwPQUQgjy2TpGxqKdw59i",
	"GWMXTcn7VCM9jSR3F+AGE30RnycaadFuKGgwItazzw7NosyF1rQbZqElUmiDNBfgcr5nCrNMuVTF7FqZ",
	"aWXKtFoiklRWArETsCUpYQZC4MLcuQ1jAC4Mtv2pCWrN2i46/6gNcdF25xFupJwwhnDVYxbxayVspiqj",
	"PIeEbX4q1EX6rOSxp5CLUhgUM9QY53takJQIRdSrKIo3R1wpehaRZu6t85lJhanIwpqoKovL4bxcOPgs",
	"bdxSVf3kIa6W9yyXlPvvnw+zr7hoSd+/dURcGW5HDTzmpmurw7jf4nkTSfuIm/VGveHGjELBFIcQFuqN",
	"+gLUQDHTceyCXLwARYeJyNeP1O4e9hbMntaMIYQ8UUt5nMUglqBB0hCuboOtAdhIkXpQA8ESy8nevRnD",
	"ZFUbSrGWD9Zqmqz57ajNPRn3vKPHWjOlujxyLIN17WdwAX42E+MOOV5cTaGUWeh8N5/Tw2VGcc3Tm9Fd",
	"51ajMZdYf+uLk/9+3Elnvfjr0/vR4WE2+DjaPTh59TwbfMgGL7LBW3CRLZZ2zaVR8fOnhEQqcEthZDD+",
	"D/OYGug0SRj1IIRs+C4bfsl2vmc7x9nwaHT0bbS372ICa9MgZsad3cYS89uKf4Dmvo25YudfZ3R2Rg+y",
	"4VG28zUb7rrF8ejl/ujn64mkTs798rY22cvhWu9Kep+8+fF777NH0Eibp4WQUhdC6BijNPTX+n8CAAD/",
	"/wyrUBlODQAA",
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
