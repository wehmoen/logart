package client

import (
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"os"
)

const host = "https://api.logart.app"

type Client struct {
	*zap.Logger
	clientLogger *zap.Logger
	apiKey       string
	host         string
	project      string
	level        zapcore.Level
}

type HTTPWriteSyncer struct {
	url       string
	userAgent string
	apiKey    string
	project   string
}

func (h *HTTPWriteSyncer) Write(p []byte) (n int, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", h.url+"/log", bytes.NewReader(p))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", h.userAgent) // Adding the User-Agent header
	req.Header.Set("X-API-Key", h.apiKey)
	req.Header.Set("X-Project", h.project)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return resp.StatusCode, nil
}
func (h *HTTPWriteSyncer) Sync() error {
	return nil
}

func New(apiKey string, project string, level zapcore.Level) *Client {

	client := &Client{}

	client.SetHost(host)
	client.SetProject(project)
	client.SetApiKey(apiKey)
	client.SetLevel(level)

	return client
}

func NewWithModule(apiKey string, project string, level zapcore.Level, module string) *Client {

	client := &Client{}

	client.SetHost(host)
	client.SetProject(project)
	client.SetApiKey(apiKey)
	client.SetLevel(level)

	client.With(zap.String("module", module))

	return client
}

func (c *Client) Local() *zap.Logger {
	return c.clientLogger
}

func (c *Client) Project() string {
	return c.project
}

func (c *Client) SetProject(project string) {
	c.project = project
}

func (c *Client) SetHost(host string) {
	c.host = host
	config := zap.NewProductionEncoderConfig()
	config.EncodeCaller = zapcore.FullCallerEncoder
	jsonEncoder := zapcore.NewJSONEncoder(config)

	// HTTP core
	httpWriter := &HTTPWriteSyncer{
		url:       host,
		userAgent: c.userAgent(),
		apiKey:    c.apiKey,
		project:   c.project,
	}
	httpCore := zapcore.NewCore(jsonEncoder, httpWriter, zapcore.DebugLevel)

	// Console core
	consoleCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(zapcore.Lock(os.Stdout)), c.level)

	c.clientLogger = zap.New(consoleCore, zap.AddCaller())

	// Combine cores
	combinedCore := zapcore.NewTee(httpCore, consoleCore)

	logger := zap.New(combinedCore, zap.AddCaller())
	c.Logger = logger
}

func (c *Client) Host() string {
	return c.host
}

func (c *Client) SetApiKey(apiKey string) {
	c.apiKey = apiKey
}

func (c *Client) SetLevel(level zapcore.Level) {
	c.level = level
}

func (c *Client) Version() string {
	return "0.0.1"
}

func (c *Client) userAgent() string {
	return "logart-go/" + c.Version()
}
