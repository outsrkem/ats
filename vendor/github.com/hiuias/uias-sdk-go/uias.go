// uias.go
package uias

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// 常量定义
const (
	defaultTimeout         = 10 * time.Second
	defaultCheckActionPath = "/v1/uias/action/check"
	defaultVerifyTokenPath = "/v1/verify/token"
	defaultTokenPath       = "/v1/uias/auth/token"
	maxBodySize            = 10 << 20 // 10MB
	maxBodyLogSize         = 500
	xRequestIdKey          = "X-Request-Id"
	xAuthTokenKey          = "X-Auth-Token"
	xSubjectTokenKey       = "X-Subject-Token"
	contentTypeJSON        = "application/json; charset=utf-8"
)

// 错误类型定义
var (
	ErrInvalidEndpoint    = errors.New("endpoint is required and must be a valid URL")
	ErrRequestFailed      = errors.New("request failed")
	ErrInvalidResponse    = errors.New("invalid response from server")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrInvalidCACert      = errors.New("invalid CA certificate")
	ErrRequestTimeout     = errors.New("request timeout")
	ErrAuthNotConfigured  = errors.New("authentication not configured")
	ErrResponseTooLarge   = errors.New("response body exceeds size limit")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// UIASError 自定义错误类型
type UIASError struct {
	Code     string
	Message  string
	Original error
}

func (e *UIASError) Error() string {
	if e.Original != nil {
		return fmt.Sprintf("UIAS error [%s]: %s (caused by: %v)", e.Code, e.Message, e.Original)
	}
	return fmt.Sprintf("UIAS error [%s]: %s", e.Code, e.Message)
}

func (e *UIASError) Unwrap() error {
	return e.Original
}

// WrapResponse 外层响应结构体，统一包装原始HTTP响应信息和业务数据
// T 为业务数据类型（如RespCheckActionData、RespVerifyTokenData、RespToken）
type WrapResponse[T any] struct {
	StatusCode int         `json:"status_code"`     // 原始HTTP响应状态码（如200、401、403）
	Header     http.Header `json:"header"`          // 原始HTTP响应头（键值对，值为数组形式）
	Data       *T          `json:"data,omitempty"`  // 业务层数据（请求成功时非空）
	Error      string      `json:"error,omitempty"` // 错误信息（请求失败时非空）
}

// RespCheckActionData 包含操作权限验证成功后返回的业务数据
type RespCheckActionData struct {
	Metadata struct {
		Message string `json:"message"`
		Time    int64  `json:"time"`
		Ecode   string `json:"ecode"`
	} `json:"metadata"`
	Payload struct {
		Authentication int    `json:"authentication"` // 1:有效;0:无效
		Error          string `json:"error"`
		Msg            struct {
			Action    string `json:"action"`
			Statement struct {
				Action string `json:"Action"`
				Effect string `json:"Effect"`
			} `json:"Statement"`
		} `json:"msg"`
		User struct {
			Domain struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			Id   string `json:"id"`
			Name struct {
				Account string `json:"account"`
			} `json:"name"`
		} `json:"user"`
	} `json:"payload"`
}

// RespVerifyTokenData 包含令牌验证成功后返回的业务数据
type RespVerifyTokenData struct {
	Metadata struct {
		Message string `json:"message"`
		Time    int64  `json:"time"`
		Ecode   string `json:"ecode"`
	} `json:"metadata"`
	Payload struct {
		Token struct {
			Valid     int    `json:"valid"` // 1:有效;0:无效
			IssuedAt  string `json:"issued_at"`
			ExpiresAt string `json:"expires_at"`
		} `json:"token"`
		User struct {
			Domain struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			Id   string `json:"id"`
			Name struct {
				Account string `json:"account"`
			} `json:"name"`
		} `json:"user"`
	} `json:"payload"`
}

// RespDataToken 包含令牌创建成功后返回的原始业务数据
type RespDataToken struct {
	Metadata struct {
		Message string `json:"message"`
		Time    int64  `json:"time"`
		Ecode   string `json:"ecode"`
	} `json:"metadata"`
	Payload struct {
		Token struct {
			IssuedAt  string `json:"issued_at"`
			ExpiresAt string `json:"expires_at"`
			User      struct {
				Domain struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"domain"`
				ID   string `json:"id"`
				Name struct {
					Account string `json:"account"`
				} `json:"name"`
			} `json:"user"`
			Roles []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"roles"`
		} `json:"token"`
	} `json:"payload"`
}

// RespToken 令牌创建接口的业务数据（包含响应头中的Token和原始业务数据）
type RespToken struct {
	Token string        `json:"token"` // 从响应头X-Subject-Token提取的令牌
	Body  RespDataToken `json:"body"`  // 响应体中的原始业务数据
}

// Config 包含客户端配置参数
type Config struct {
	Endpoint        string        // 服务端点 (例如: "https://api.example.com")
	UrlPath         string        // 用于请求的API路径（默认：defaultCheckActionPath）
	SkipTlsVerify   bool          // 是否跳过TLS验证 (默认: false)
	CACertPath      string        // CA证书文件路径
	Timeout         time.Duration // 请求超时时间 (默认: 10秒)
	MaxIdleConns    int           // 最大空闲连接数
	IdleConnTimeout time.Duration // 空闲连接超时时间
}

// AuthConfig 认证配置（AK/SK）
type AuthConfig struct {
	Ak string // Access Key
	Sk string // Secret Key
}

// Auth 认证信息载体
type Auth struct {
	config AuthConfig
}

// Response 内部HTTP请求结果封装（仅内部使用）
type Response struct {
	StatusCode int         // HTTP状态码
	Body       []byte      // 响应体内容
	Header     http.Header // 响应头
}

// Logger 日志接口（支持自定义日志实现）
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// DefaultLogger 默认日志实现（基于标准库log）
type DefaultLogger struct{}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// Metrics 指标接口（支持自定义监控指标收集）
type Metrics interface {
	IncRequestCounter(method, path string, statusCode int)
	ObserveRequestDuration(method, path string, duration time.Duration)
}

// NoopMetrics 空指标实现（默认无监控）
type NoopMetrics struct{}

func (m *NoopMetrics) IncRequestCounter(method, path string, statusCode int)              {}
func (m *NoopMetrics) ObserveRequestDuration(method, path string, duration time.Duration) {}

// UIASClient UIAS客户端接口（对外暴露的核心能力）
type UIASClient interface {
	VerifyAction(ctx context.Context, requestId, token string, rawBody []byte) (*WrapResponse[RespCheckActionData], error)
	CreateToken(ctx context.Context, requestId string) (*WrapResponse[RespToken], error)
	VerifyToken(ctx context.Context, requestId, token string) (*WrapResponse[RespVerifyTokenData], error)
	Close()
}

// Client UIAS服务客户端实现
type Client struct {
	auth       *Auth        // 认证信息
	config     Config       // 客户端配置
	httpClient *http.Client // HTTP客户端
	logger     Logger       // 日志实例
	metrics    Metrics      // 指标实例
}

// ClientOption 客户端选项函数（用于灵活配置客户端）
type ClientOption func(*Client)

// CredentialBuilder 凭证构建器（用于构造AK/SK认证信息）
type CredentialBuilder struct {
	config AuthConfig
}

// ClientBuilder 客户端构建器（用于构造客户端实例）
type ClientBuilder struct {
	config Config
	auth   *Auth
}

// 包装错误信息（添加上下文描述）
func wrapError(err error, message string) error {
	if err == nil {
		return errors.New(message)
	}
	return fmt.Errorf("%s: %w", message, err)
}

// 过滤敏感字段（日志中隐藏密码、密钥等信息）
func filterSensitiveFields(data map[string]interface{}) {
	sensitiveFields := []string{"password", "secret", "token", "key", "sk", "ak"}
	for _, field := range sensitiveFields {
		if _, exists := data[field]; exists {
			data[field] = "***REDACTED***"
		}
	}
}

// 清理日志中的敏感信息并截断长内容
func sanitizeBodyForLog(body []byte) string {
	// 截断超过maxBodyLogSize的内容
	if len(body) > maxBodyLogSize {
		body = body[:maxBodyLogSize]
	}

	// 尝试解析JSON并过滤敏感字段
	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err == nil {
		filterSensitiveFields(jsonData)
		sanitized, _ := json.Marshal(jsonData)
		return string(sanitized) + "..."
	}

	// 非JSON数据直接返回截断内容
	return string(body) + "..."
}

// NewCredentialBuilder 创建凭证构建器实例
func NewCredentialBuilder() *CredentialBuilder {
	return &CredentialBuilder{
		config: AuthConfig{},
	}
}

// WithAk 设置Access Key
func (b *CredentialBuilder) WithAk(ak string) *CredentialBuilder {
	b.config.Ak = ak
	return b
}

// WithSk 设置Secret Key
func (b *CredentialBuilder) WithSk(sk string) *CredentialBuilder {
	b.config.Sk = sk
	return b
}

// Build 构建认证信息实例（校验AK/SK非空）
func (b *CredentialBuilder) Build() *Auth {
	if b.config.Ak == "" || b.config.Sk == "" {
		panic("ak and sk cannot be empty")
	}
	return &Auth{config: b.config}
}

// NewClientBuilder 创建客户端构建器实例（默认配置初始化）
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		config: Config{
			Timeout:         defaultTimeout,
			SkipTlsVerify:   false,
			UrlPath:         defaultCheckActionPath,
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}
}

// WithEndpoint 设置服务端点（如"https://api.example.com"）
func (b *ClientBuilder) WithEndpoint(endpoint string) *ClientBuilder {
	b.config.Endpoint = endpoint
	return b
}

// WithUrlPath 设置请求API路径（覆盖默认路径）
func (b *ClientBuilder) WithUrlPath(path string) *ClientBuilder {
	b.config.UrlPath = path
	return b
}

// WithCredential 设置认证信息（关联AK/SK）
func (b *ClientBuilder) WithCredential(auth *Auth) *ClientBuilder {
	if auth != nil {
		b.auth = auth
	}
	return b
}

// WithSkipTlsVerify 设置是否跳过TLS证书验证（生产环境不推荐）
func (b *ClientBuilder) WithSkipTlsVerify(skip bool) *ClientBuilder {
	b.config.SkipTlsVerify = skip
	return b
}

// WithCACertPath 设置CA证书文件路径（用于自定义证书校验）
func (b *ClientBuilder) WithCACertPath(path string) *ClientBuilder {
	b.config.CACertPath = path
	return b
}

// WithTimeout 设置请求超时时间（覆盖默认10秒）
func (b *ClientBuilder) WithTimeout(timeout time.Duration) *ClientBuilder {
	b.config.Timeout = timeout
	return b
}

// WithConnectionPool 设置HTTP连接池参数
func (b *ClientBuilder) WithConnectionPool(maxIdleConns int, idleTimeout time.Duration) *ClientBuilder {
	b.config.MaxIdleConns = maxIdleConns
	b.config.IdleConnTimeout = idleTimeout
	return b
}

// Build 构建客户端实例（恐慌式初始化，适合明确配置场景）
func (b *ClientBuilder) Build() *Client {
	client, err := SafeBuild(b.config, b.auth)
	if err != nil {
		panic(fmt.Sprintf("failed to build UIAS client: %v", err))
	}
	return client
}

// SafeBuild 安全构建客户端实例（返回错误，适合需处理初始化失败的场景）
func SafeBuild(config Config, auth *Auth, opts ...ClientOption) (*Client, error) {
	// 校验服务端点非空
	if config.Endpoint == "" {
		return nil, ErrInvalidEndpoint
	}

	// 创建HTTP客户端（包含TLS配置、连接池）
	httpClient, err := createHttpClient(config.SkipTlsVerify, config.CACertPath, config.MaxIdleConns, config.IdleConnTimeout)
	if err != nil {
		return nil, wrapError(err, "failed to create http client")
	}

	// 设置HTTP请求超时
	httpClient.Timeout = config.Timeout

	// 初始化客户端基础结构
	client := &Client{
		auth:       auth,
		config:     config,
		httpClient: httpClient,
		logger:     &DefaultLogger{}, // 默认日志
		metrics:    &NoopMetrics{},   // 默认无指标
	}

	// 应用客户端选项（如自定义日志、指标）
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// WithLogger 设置自定义日志实例
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		if logger != nil {
			c.logger = logger
		}
	}
}

// WithMetrics 设置自定义指标实例
func WithMetrics(metrics Metrics) ClientOption {
	return func(c *Client) {
		if metrics != nil {
			c.metrics = metrics
		}
	}
}

// createTLSConfig 创建TLS配置（支持跳过验证、自定义CA证书）
func createTLSConfig(skipTlsVerify bool, caCertPath string) (*tls.Config, error) {
	// 跳过TLS验证（仅测试环境使用）
	if skipTlsVerify {
		return &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12, // 强制TLS 1.2及以上
		}, nil
	}

	// 基础TLS配置（安全优先）
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{ // 优先使用安全曲线
			tls.CurveP256,
			tls.X25519,
		},
	}

	// 加载自定义CA证书（如需校验私有证书）
	if caCertPath != "" {
		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			return nil, wrapError(err, "failed to read CA certificate file")
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, ErrInvalidCACert
		}
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}

// createHttpClient 创建HTTP客户端（包含TLS配置和连接池）
func createHttpClient(skipTlsVerify bool, caCertPath string, maxIdleConns int, idleConnTimeout time.Duration) (*http.Client, error) {
	// 创建TLS配置
	tlsConfig, err := createTLSConfig(skipTlsVerify, caCertPath)
	if err != nil {
		return nil, err
	}

	// 创建传输层配置（包含连接池）
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 拨号超时
			KeepAlive: 30 * time.Second, // 长连接保持时间
		}).DialContext,
		MaxIdleConns:          maxIdleConns,     // 最大空闲连接数
		IdleConnTimeout:       idleConnTimeout,  // 空闲连接超时时间
		TLSHandshakeTimeout:   10 * time.Second, // TLS握手超时
		ExpectContinueTimeout: 1 * time.Second,  // 100-continue响应超时
		ForceAttemptHTTP2:     true,             // 优先使用HTTP/2
		TLSClientConfig:       tlsConfig,        // 关联TLS配置
	}

	return &http.Client{Transport: transport}, nil
}

// sendHttpRequest 发送HTTP请求（内部通用实现）
func sendHttpRequest(
	ctx context.Context,
	client *http.Client,
	method, url string,
	body []byte,
	headers map[string]string,
) (*Response, error) {
	// 创建请求体读取器
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, wrapError(err, "failed to create request")
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 执行请求并记录耗时
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	// 记录请求耗时指标（如果客户端传输层支持）
	if metrics, ok := client.Transport.(interface {
		ObserveRequestDuration(method, url string, duration time.Duration)
	}); ok {
		metrics.ObserveRequestDuration(method, url, duration)
	}

	// 处理请求错误
	if err != nil {
		// 区分超时错误
		if os.IsTimeout(err) || errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrRequestTimeout
		}
		return nil, wrapError(ErrRequestFailed, err.Error())
	}
	defer resp.Body.Close()

	// 记录请求状态码指标（如果客户端传输层支持）
	if metrics, ok := client.Transport.(interface {
		IncRequestCounter(method, url string, statusCode int)
	}); ok {
		metrics.IncRequestCounter(method, url, resp.StatusCode)
	}

	// 读取响应体（带大小限制）
	limitedReader := &io.LimitedReader{R: resp.Body, N: maxBodySize}
	responseBody, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, wrapError(ErrInvalidResponse, "failed to read response body")
	}

	// 检查响应体是否超过大小限制
	if limitedReader.N <= 0 {
		return nil, ErrResponseTooLarge
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       responseBody,
		Header:     resp.Header,
	}, nil
}

// doRequest 客户端封装的请求方法（带日志）
func (c *Client) doRequest(ctx context.Context, method, path string, body []byte, headers map[string]string) (*Response, error) {
	url := c.config.Endpoint + path
	c.logger.Debugf("Making request: %s %s, headers: %v, body: %s",
		method, url, headers, sanitizeBodyForLog(body))

	// 发送HTTP请求
	response, err := sendHttpRequest(ctx, c.httpClient, method, url, body, headers)
	if err != nil {
		c.logger.Errorf("Request failed: %s %s: %v", method, url, err)
		return nil, err
	}

	c.logger.Debugf("Response received: %s %s -> status: %d, body: %s",
		method, url, response.StatusCode, sanitizeBodyForLog(response.Body))
	return response, nil
}

// parseResponse 解析响应体到指定结构体
func (c *Client) parseResponse(response *Response, result interface{}) error {
	if err := json.Unmarshal(response.Body, result); err != nil {
		c.logger.Errorf("Failed to parse response: invalid JSON - %s, body: %s",
			err.Error(), sanitizeBodyForLog(response.Body))
		return wrapError(ErrInvalidResponse, "failed to parse response body")
	}
	return nil
}

// VerifyAction 验证操作权限
func (c *Client) VerifyAction(ctx context.Context, requestId, token string, rawBody []byte) (*WrapResponse[RespCheckActionData], error) {
	// 构建请求头
	headers := map[string]string{
		"Content-Type": contentTypeJSON,
		xRequestIdKey:  requestId,
		xAuthTokenKey:  token,
	}

	// 发送请求
	response, err := c.doRequest(ctx, "POST", c.config.UrlPath, rawBody, headers)
	if err != nil {
		// 构建错误响应
		wrapResp := &WrapResponse[RespCheckActionData]{
			Error: err.Error(),
		}
		// 填充已知错误的状态码
		switch {
		case errors.Is(err, ErrRequestTimeout):
			wrapResp.StatusCode = http.StatusRequestTimeout
		case errors.Is(err, ErrPermissionDenied):
			wrapResp.StatusCode = http.StatusForbidden
		case errors.Is(err, ErrInvalidCredentials):
			wrapResp.StatusCode = http.StatusUnauthorized
		}
		return wrapResp, err
	}

	// 解析业务数据
	var respData RespCheckActionData
	if err := c.parseResponse(response, &respData); err != nil {
		return &WrapResponse[RespCheckActionData]{
			StatusCode: response.StatusCode,
			Header:     response.Header,
			Error:      err.Error(),
		}, err
	}

	// 构建成功响应
	return &WrapResponse[RespCheckActionData]{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Data:       &respData,
	}, nil
}

// buildTokenRequest 构建创建令牌的请求体
func (c *Client) buildTokenRequest() ([]byte, error) {
	type rawReq struct {
		Auth struct {
			Credential struct {
				ID     string `json:"id"`
				Secret string `json:"secret"`
			} `json:"credential"`
		} `json:"auth"`
	}

	var raw rawReq
	raw.Auth.Credential.ID = c.auth.config.Ak
	raw.Auth.Credential.Secret = c.auth.config.Sk

	rawBody, err := json.Marshal(raw)
	if err != nil {
		return nil, wrapError(err, "failed to marshal token request")
	}

	return rawBody, nil
}

// CreateToken 创建令牌
func (c *Client) CreateToken(ctx context.Context, requestId string) (*WrapResponse[RespToken], error) {
	// 检查认证配置
	if c.auth == nil {
		err := ErrAuthNotConfigured
		return &WrapResponse[RespToken]{
			Error: err.Error(),
		}, err
	}

	// 构建请求头
	headers := map[string]string{
		"Content-Type": contentTypeJSON,
		xRequestIdKey:  requestId,
	}

	// 构建请求体
	requestBody, err := c.buildTokenRequest()
	if err != nil {
		return &WrapResponse[RespToken]{
			Error: err.Error(),
		}, err
	}

	// 发送请求
	response, err := c.doRequest(ctx, "POST", defaultTokenPath, requestBody, headers)
	if err != nil {
		wrapResp := &WrapResponse[RespToken]{
			Error: err.Error(),
		}
		if errors.Is(err, ErrRequestTimeout) {
			wrapResp.StatusCode = http.StatusRequestTimeout
		}
		return wrapResp, err
	}

	// 解析响应体业务数据
	var respData RespDataToken
	if err := json.Unmarshal(response.Body, &respData); err != nil {
		c.logger.Errorf("Invalid token response JSON: %s", sanitizeBodyForLog(response.Body))
		return &WrapResponse[RespToken]{
			StatusCode: response.StatusCode,
			Header:     response.Header,
			Error:      wrapError(ErrInvalidResponse, "failed to parse token response").Error(),
		}, err
	}

	// 提取响应头中的令牌
	token := response.Header.Get(xSubjectTokenKey)
	respToken := &RespToken{
		Token: token,
		Body:  respData,
	}

	// 构建成功响应
	return &WrapResponse[RespToken]{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Data:       respToken,
	}, nil
}

// VerifyToken 验证令牌
func (c *Client) VerifyToken(ctx context.Context, requestId, token string) (*WrapResponse[RespVerifyTokenData], error) {
	// 构建请求头
	headers := map[string]string{
		"Content-Type": contentTypeJSON,
		xRequestIdKey:  requestId,
		xAuthTokenKey:  token,
	}

	// 发送请求
	response, err := c.doRequest(ctx, "POST", defaultVerifyTokenPath, nil, headers)
	if err != nil {
		wrapResp := &WrapResponse[RespVerifyTokenData]{
			Error: err.Error(),
		}
		switch {
		case errors.Is(err, ErrRequestTimeout):
			wrapResp.StatusCode = http.StatusRequestTimeout
		case errors.Is(err, ErrPermissionDenied):
			wrapResp.StatusCode = http.StatusForbidden
		}
		return wrapResp, err
	}

	// 解析业务数据
	var respData RespVerifyTokenData
	if err := c.parseResponse(response, &respData); err != nil {
		return &WrapResponse[RespVerifyTokenData]{
			StatusCode: response.StatusCode,
			Header:     response.Header,
			Error:      err.Error(),
		}, err
	}

	// 构建成功响应
	return &WrapResponse[RespVerifyTokenData]{
		StatusCode: response.StatusCode,
		Header:     response.Header,
		Data:       &respData,
	}, nil
}

// Close 释放客户端资源（关闭空闲连接）
func (c *Client) Close() {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}
