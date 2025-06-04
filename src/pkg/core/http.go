package core

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// CreateHttpClient initializes an HTTP client with optional TLS configuration.
// It supports custom CA certificates and TLS verification skipping.
// Note: skipTlsVerify is ignored when a custom CA certificate is provided.
func CreateHttpClient(klog *logrus.Entry, skipTlsVerify bool, caCertPath string) (*http.Client, error) {
	var tlsConfig *tls.Config
	switch {
	case caCertPath != "":
		caCert, err := os.ReadFile(caCertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		// Create certificate pool (includes system certificates if available)
		caCertPool, _ := x509.SystemCertPool()
		if caCertPool == nil {
			caCertPool = x509.NewCertPool()
		}

		// Add custom CA certificate to the pool
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, errors.New("failed to parse CA certificate")
		}

		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}

		if skipTlsVerify {
			klog.Warnf("warning: skipTlsVerify=true is ignored when using custom CA: %s", caCertPath)
		}

	case skipTlsVerify:
		// Skip TLS verification (INSECURE - for testing only)
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if tlsConfig != nil {
		transport.TLSClientConfig = tlsConfig
	}

	return &http.Client{Transport: transport}, nil
}

// Response encapsulates the results of an HTTP request.
type Response struct {
	StatusCode int         // e.g. 200
	Body       []byte      // Response body content
	Header     http.Header // Response headers
}

// SendHttpRequest sends an HTTP request and handles the response safely.
// Note: The returned responseBody contains the fully read content,
// so there's no need to interact with the original response body afterward.
func SendHttpRequest(klog *logrus.Entry, ctx context.Context, client *http.Client, method, url string, body []byte, headers map[string]string, timeout time.Duration) (*Response, error) {
	// Set default timeout if not provided
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	// Create context with timeout and ensure cancellation
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Prepare request body
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	// Initialize response structure
	_resp := &Response{}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return _resp, err
	}

	// Apply headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return _resp, err
	}

	// Ensure response body is properly closed after reading
	defer func() {
		if resp != nil && resp.Body != nil {
			// Safely discard remaining response body and close
			if _, err := io.Copy(io.Discard, resp.Body); err != nil {
				klog.Errorf("Error discarding response body: %v", err)
			}
			if err := resp.Body.Close(); err != nil {
				klog.Errorf("Error closing response body: %v", err)
			}
		}
	}()

	// Read full response body into memory
	_resp.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		return _resp, err
	}

	// Populate response metadata
	_resp.Header = resp.Header
	_resp.StatusCode = resp.StatusCode

	return _resp, nil
}
