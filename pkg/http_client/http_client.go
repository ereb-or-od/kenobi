package http_client

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"crypto/x509"
	marshallers "github.com/ereb-or-od/kenobi/pkg/marshalling/interfaces"
	"github.com/ereb-or-od/kenobi/pkg/marshalling/json"

	"errors"
	"fmt"
	"github.com/ereb-or-od/kenobi/pkg/logging"
	logger "github.com/ereb-or-od/kenobi/pkg/logging/interfaces"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	// MethodGet HTTP method
	MethodGet = "GET"

	// MethodPost HTTP method
	MethodPost = "POST"

	// MethodPut HTTP method
	MethodPut = "PUT"

	// MethodDelete HTTP method
	MethodDelete = "DELETE"

	// MethodPatch HTTP method
	MethodPatch = "PATCH"

	// MethodHead HTTP method
	MethodHead = "HEAD"

	// MethodOptions HTTP method
	MethodOptions = "OPTIONS"
)

var (
	hdrUserAgentKey       = http.CanonicalHeaderKey("User-Agent")
	hdrAcceptKey          = http.CanonicalHeaderKey("Accept")
	hdrContentTypeKey     = http.CanonicalHeaderKey("Content-Type")
	hdrContentLengthKey   = http.CanonicalHeaderKey("Content-Length")
	hdrContentEncodingKey = http.CanonicalHeaderKey("Content-Encoding")

	plainTextType   = "text/plain; charset=utf-8"
	jsonContentType = "application/json"
	formContentType = "application/x-www-form-urlencoded"

	jsonCheck = regexp.MustCompile(`(?i:(application|text)/(json|.*\+json|json\-.*)(;|$))`)
	xmlCheck  = regexp.MustCompile(`(?i:(application|text)/(xml|.*\+xml)(;|$))`)

	bufPool = &sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}
)

type (
	// RequestMiddleware type is for request middleware, called before a request is sent
	RequestMiddleware func(*HttpClient, *Request) error

	// ResponseMiddleware type is for response middleware, called after a response has been received
	ResponseMiddleware func(*HttpClient, *Response) error

	// PreRequestHook type is for the request hook, called right before the request is sent
	PreRequestHook func(*HttpClient, *http.Request) error

	// RequestLogCallback type is for request logs, called before the request is logged
	RequestLogCallback func(*RequestLog) error

	// ResponseLogCallback type is for response logs, called before the response is logged
	ResponseLogCallback func(*ResponseLog) error

	// ErrorHook type is for reacting to request errors, called after all retries were attempted
	ErrorHook func(*Request, error)
)

type HttpClient struct {
	BaseUrl               string
	QueryParam            url.Values
	FormData              url.Values
	Header                http.Header
	UserInfo              *User
	Token                 string
	AuthScheme            string
	Cookies               []*http.Cookie
	Error                 reflect.Type
	Debug                 bool
	DisableWarn           bool
	AllowGetMethodPayload bool
	RetryCount            int
	RetryWaitTime         time.Duration
	RetryMaxWaitTime      time.Duration
	RetryConditions       []RetryConditionFunc
	RetryHooks            []OnRetryFunc
	RetryAfter            RetryAfterFunc
	Marshaller            marshallers.Marshaller

	HeaderAuthorizationKey string

	jsonEscapeHTML     bool
	setContentLength   bool
	closeConnection    bool
	notParseResponse   bool
	trace              bool
	debugBodySizeLimit int64
	outputDirectory    string
	scheme             string
	pathParams         map[string]string
	log                logger.Logger
	httpClient         *http.Client
	proxyURL           *url.URL
	beforeRequest      []RequestMiddleware
	udBeforeRequest    []RequestMiddleware
	preReqHook         PreRequestHook
	afterResponse      []ResponseMiddleware
	requestLog         RequestLogCallback
	responseLog        ResponseLogCallback
	errorHooks         []ErrorHook
}

type User struct {
	Username, Password string
}

func (c *HttpClient) UseBaseUrl(url string) *HttpClient {
	c.BaseUrl = strings.TrimRight(url, "/")
	return c
}

func (c *HttpClient) UseHeader(header, value string) *HttpClient {
	c.Header.Set(header, value)
	return c
}

func (c *HttpClient) UseHeaders(headers map[string]string) *HttpClient {
	for h, v := range headers {
		c.Header.Set(h, v)
	}
	return c
}

func (c *HttpClient) UseHeaderVerbatim(header, value string) *HttpClient {
	c.Header[header] = []string{value}
	return c
}

func (c *HttpClient) UseCookieJar(jar http.CookieJar) *HttpClient {
	c.httpClient.Jar = jar
	return c
}

func (c *HttpClient) UseCookie(hc *http.Cookie) *HttpClient {
	c.Cookies = append(c.Cookies, hc)
	return c
}

func (c *HttpClient) UseCookies(cs []*http.Cookie) *HttpClient {
	c.Cookies = append(c.Cookies, cs...)
	return c
}

func (c *HttpClient) UseQueryParam(param, value string) *HttpClient {
	c.QueryParam.Set(param, value)
	return c
}

func (c *HttpClient) UseQueryParams(params map[string]string) *HttpClient {
	for p, v := range params {
		c.UseQueryParam(p, v)
	}
	return c
}

func (c *HttpClient) UseFormData(data map[string]string) *HttpClient {
	for k, v := range data {
		c.FormData.Set(k, v)
	}
	return c
}

func (c *HttpClient) UseBasicAuthentication(username, password string) *HttpClient {
	c.UserInfo = &User{Username: username, Password: password}
	return c
}

func (c *HttpClient) UseAuthenticationToken(token string) *HttpClient {
	c.Token = token
	return c
}

func (c *HttpClient) UseAuthenticationSchema(scheme string) *HttpClient {
	c.AuthScheme = scheme
	return c
}

func (c *HttpClient) NewRequest() *Request {
	r := &Request{
		QueryParam: url.Values{},
		FormData:   url.Values{},
		Header:     http.Header{},
		Cookies:    make([]*http.Cookie, 0),

		client:          c,
		multipartFiles:  []*File{},
		multipartFields: []*MultipartField{},
		pathParams:      map[string]string{},
		jsonEscapeHTML:  true,
	}
	return r
}

func (c *HttpClient) new() *Request {
	return c.NewRequest()
}

func (c *HttpClient) OnBeforeRequest(m RequestMiddleware) *HttpClient {
	c.udBeforeRequest = append(c.udBeforeRequest, m)
	return c
}

func (c *HttpClient) OnAfterResponse(m ResponseMiddleware) *HttpClient {
	c.afterResponse = append(c.afterResponse, m)
	return c
}

func (c *HttpClient) OnError(h ErrorHook) *HttpClient {
	c.errorHooks = append(c.errorHooks, h)
	return c
}

func (c *HttpClient) UsePreRequestHook(h PreRequestHook) *HttpClient {
	if c.preReqHook != nil {
		c.log.Warn(fmt.Sprintf("Overwriting an existing pre-request hook: %s", functionName(h)))
	}
	c.preReqHook = h
	return c
}

func (c *HttpClient) UseDebug(d bool) *HttpClient {
	c.Debug = d
	return c
}

func (c *HttpClient) UseDebugBodyLimit(sl int64) *HttpClient {
	c.debugBodySizeLimit = sl
	return c
}

func (c *HttpClient) OnRequestLog(rl RequestLogCallback) *HttpClient {
	if c.requestLog != nil {
		c.log.Warn(fmt.Sprintf("Overwriting an existing on-request-log callback from=%s to=%s",
			functionName(c.requestLog), functionName(rl)))
	}
	c.requestLog = rl
	return c
}

func (c *HttpClient) OnResponseLog(rl ResponseLogCallback) *HttpClient {
	if c.responseLog != nil {
		c.log.Warn(fmt.Sprintf("Overwriting an existing on-response-log callback from=%s to=%s",
			functionName(c.responseLog), functionName(rl)))
	}
	c.responseLog = rl
	return c
}

func (c *HttpClient) UseDisableWarnings(d bool) *HttpClient {
	c.DisableWarn = d
	return c
}

func (c *HttpClient) UseAllowGetMethodPayload(a bool) *HttpClient {
	c.AllowGetMethodPayload = a
	return c
}

func (c *HttpClient) UseLogger(l logger.Logger) *HttpClient {
	c.log = l
	return c
}

func (c *HttpClient) UseContentLength(l bool) *HttpClient {
	c.setContentLength = l
	return c
}

func (c *HttpClient) UseTimeout(timeout time.Duration) *HttpClient {
	c.httpClient.Timeout = timeout
	return c
}

func (c *HttpClient) UseError(err interface{}) *HttpClient {
	c.Error = typeOf(err)
	return c
}

func (c *HttpClient) UseRedirectPolicy(policies ...interface{}) *HttpClient {
	for _, p := range policies {
		if _, ok := p.(RedirectPolicy); !ok {
			c.log.Warn(fmt.Sprintf("%v does not implemented.RedirectPolicy (missing Apply method)",
				functionName(p)))
		}
	}

	c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		for _, p := range policies {
			if err := p.(RedirectPolicy).Apply(req, via); err != nil {
				return err
			}
		}
		return nil // looks good, go ahead
	}

	return c
}

func (c *HttpClient) UseRetryCount(count int) *HttpClient {
	c.RetryCount = count
	return c
}

func (c *HttpClient) UseRetryWaitTime(waitTime time.Duration) *HttpClient {
	c.RetryWaitTime = waitTime
	return c
}

func (c *HttpClient) UseRetryMaxWaitTime(maxWaitTime time.Duration) *HttpClient {
	c.RetryMaxWaitTime = maxWaitTime
	return c
}

func (c *HttpClient) UseRetryAfter(callback RetryAfterFunc) *HttpClient {
	c.RetryAfter = callback
	return c
}

func (c *HttpClient) AddRetryCondition(condition RetryConditionFunc) *HttpClient {
	c.RetryConditions = append(c.RetryConditions, condition)
	return c
}

func (c *HttpClient) AddRetryAfterErrorCondition() *HttpClient {
	c.AddRetryCondition(func(response *Response, err error) bool {
		return response.IsError()
	})
	return c
}

func (c *HttpClient) AddRetryHook(hook OnRetryFunc) *HttpClient {
	c.RetryHooks = append(c.RetryHooks, hook)
	return c
}

func (c *HttpClient) UseTLSClientConfiguration(config *tls.Config) *HttpClient {
	transport, err := c.transport()
	if err != nil {
		c.log.Error("an unexpected error occurred", err)
		return c
	}
	transport.TLSClientConfig = config
	return c
}

func (c *HttpClient) UseProxy(proxyURL string) *HttpClient {
	transport, err := c.transport()
	if err != nil {
		c.log.Error("an unexpected error occurred", err)
		return c
	}

	pURL, err := url.Parse(proxyURL)
	if err != nil {
		c.log.Error("an unexpected error occurred", err)
		return c
	}

	c.proxyURL = pURL
	transport.Proxy = http.ProxyURL(c.proxyURL)
	return c
}

func (c *HttpClient) RemoveProxy() *HttpClient {
	transport, err := c.transport()
	if err != nil {
		c.log.Error("an unexpected error occurred", err)
		return c
	}
	c.proxyURL = nil
	transport.Proxy = nil
	return c
}
func (c *HttpClient) UseCertificates(certs ...tls.Certificate) *HttpClient {
	config, err := c.tlsConfig()
	if err != nil {
		c.log.Error("an unexpected error occurred", err)
		return c
	}
	config.Certificates = append(config.Certificates, certs...)
	return c
}

func (c *HttpClient) UseRootCertificate(pemFilePath string) *HttpClient {
	rootPemData, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		c.log.Error("an unexpected error occurred", err)
		return c
	}

	config, err := c.tlsConfig()
	if err != nil {
		c.log.Error("an unexpected error occurred", err)
		return c
	}
	if config.RootCAs == nil {
		config.RootCAs = x509.NewCertPool()
	}

	config.RootCAs.AppendCertsFromPEM(rootPemData)
	return c
}

func (c *HttpClient) UseRootCertificateFromString(pemContent string) *HttpClient {
	config, err := c.tlsConfig()
	if err != nil {
		c.log.Error("an unexpected error occurred", err)
		return c
	}
	if config.RootCAs == nil {
		config.RootCAs = x509.NewCertPool()
	}

	config.RootCAs.AppendCertsFromPEM([]byte(pemContent))
	return c
}
func (c *HttpClient) UseOutputDirectory(dirPath string) *HttpClient {
	c.outputDirectory = dirPath
	return c
}

func (c *HttpClient) UseTransport(transport http.RoundTripper) *HttpClient {
	if transport != nil {
		c.httpClient.Transport = transport
	}
	return c
}

func (c *HttpClient) UseSchema(scheme string) *HttpClient {
	if !IsStringEmpty(scheme) {
		c.scheme = scheme
	}
	return c
}

func (c *HttpClient) UseCloseConnection(close bool) *HttpClient {
	c.closeConnection = close
	return c
}

func (c *HttpClient) UseDoNotParseResponse(parse bool) *HttpClient {
	c.notParseResponse = parse
	return c
}

func (c *HttpClient) UsePathParam(param, value string) *HttpClient {
	c.pathParams[param] = value
	return c
}

func (c *HttpClient) UsePathParams(params map[string]string) *HttpClient {
	for p, v := range params {
		c.UsePathParam(p, v)
	}
	return c
}

func (c *HttpClient) UseJSONEscapeHTML(b bool) *HttpClient {
	c.jsonEscapeHTML = b
	return c
}

func (c *HttpClient) EnableTrace() *HttpClient {
	c.trace = true
	return c
}

func (c *HttpClient) DisableTrace() *HttpClient {
	c.trace = false
	return c
}

func (c *HttpClient) IsProxySet() bool {
	return c.proxyURL != nil
}

func (c *HttpClient) GetClient() *http.Client {
	return c.httpClient
}

func (c *HttpClient) execute(req *Request) (*Response, error) {
	defer releaseBuffer(req.bodyBuf)
	// Apply Request middleware
	var err error

	// user defined on before request methods
	for _, f := range c.udBeforeRequest {
		if err = f(c, req); err != nil {
			return nil, wrapNoRetryErr(err)
		}
	}

	for _, f := range c.beforeRequest {
		if err = f(c, req); err != nil {
			return nil, wrapNoRetryErr(err)
		}
	}

	if hostHeader := req.Header.Get("Host"); hostHeader != "" {
		req.RawRequest.Host = hostHeader
	}

	// call pre-request if defined
	if c.preReqHook != nil {
		if err = c.preReqHook(c, req.RawRequest); err != nil {
			return nil, wrapNoRetryErr(err)
		}
	}

	if err = requestLogger(c, req); err != nil {
		return nil, wrapNoRetryErr(err)
	}

	req.Time = time.Now()
	resp, err := c.httpClient.Do(req.RawRequest)

	response := &Response{
		Request:     req,
		RawResponse: resp,
	}

	if err != nil || req.notParseResponse || c.notParseResponse {
		response.setReceivedAt()
		return response, err
	}

	if !req.isSaveResponse {
		defer closeq(resp.Body)
		body := resp.Body

		// GitHub #142 & #187
		if strings.EqualFold(resp.Header.Get(hdrContentEncodingKey), "gzip") && resp.ContentLength != 0 {
			if _, ok := body.(*gzip.Reader); !ok {
				body, err = gzip.NewReader(body)
				if err != nil {
					response.setReceivedAt()
					return response, err
				}
				defer closeq(body)
			}
		}

		if response.body, err = ioutil.ReadAll(body); err != nil {
			response.setReceivedAt()
			return response, err
		}

		response.size = int64(len(response.body))
	}

	response.setReceivedAt() // after we read the body

	// Apply Response middleware
	for _, f := range c.afterResponse {
		if err = f(c, response); err != nil {
			break
		}
	}

	return response, wrapNoRetryErr(err)
}

func (c *HttpClient) tlsConfig() (*tls.Config, error) {
	transport, err := c.transport()
	if err != nil {
		return nil, err
	}
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{}
	}
	return transport.TLSClientConfig, nil
}

func (c *HttpClient) transport() (*http.Transport, error) {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		return transport, nil
	}
	return nil, errors.New("current transport is not an *http.Transport instance")
}

type ResponseError struct {
	Response *Response
	Err      error
}

func (e *ResponseError) Error() string {
	return e.Err.Error()
}

func (e *ResponseError) Unwrap() error {
	return e.Err
}

func (c *HttpClient) onErrorHooks(req *Request, resp *Response, err error) {
	if err != nil {
		if resp != nil { // wrap with ResponseError
			err = &ResponseError{Response: resp, Err: err}
		}
		for _, h := range c.errorHooks {
			h(req, err)
		}
	}
}

type File struct {
	Name      string
	ParamName string
	io.Reader
}

func (f *File) String() string {
	return fmt.Sprintf("ParamName: %v; FileName: %v", f.ParamName, f.Name)
}

type MultipartField struct {
	Param       string
	FileName    string
	ContentType string
	io.Reader
}

func createClient(hc *http.Client) *HttpClient {
	if hc.Transport == nil {
		hc.Transport = createTransport(nil)
	}

	c := &HttpClient{ // not setting lang default values
		QueryParam:             url.Values{},
		FormData:               url.Values{},
		Header:                 http.Header{},
		Cookies:                make([]*http.Cookie, 0),
		RetryWaitTime:          defaultWaitTime,
		RetryMaxWaitTime:       defaultMaxWaitTime,
		Marshaller:             json.New(),
		HeaderAuthorizationKey: http.CanonicalHeaderKey("Authorization"),
		jsonEscapeHTML:         true,
		httpClient:             hc,
		debugBodySizeLimit:     math.MaxInt32,
		pathParams:             make(map[string]string),
	}

	// Logger
	defaultLogger, _ := logging.New()
	c.UseLogger(defaultLogger)

	// default before request middlewares
	c.beforeRequest = []RequestMiddleware{
		parseRequestURL,
		parseRequestHeader,
		parseRequestBody,
		createHTTPRequest,
		addCredentials,
	}

	// user defined request middlewares
	c.udBeforeRequest = []RequestMiddleware{}

	// default after response middlewares
	c.afterResponse = []ResponseMiddleware{
		responseLogger,
		parseResponseBody,
		saveResponseIntoFile,
	}

	return c
}
