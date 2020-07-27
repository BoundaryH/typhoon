package typhoon

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"
)

// Config represents Typhoo Configuration
type Config struct {
	Target    string
	NumCPU    int
	NumThread int
	Duration  time.Duration

	Method    string
	Header    http.Header
	Cookie    string
	UserAgent string
	Body      []byte

	NoCompression bool
	KeepAlive     bool
	SkipTLSVerify bool
	CertPool      *x509.CertPool
	ClientCert    *tls.Certificate
}

// DefaultConfig returns a configuration with GET Method
func DefaultConfig(target string) *Config {
	return &Config{
		Target:    target,
		NumCPU:    0,
		NumThread: 1,
		Duration:  time.Second,

		Method:    "GET",
		Header:    nil,
		Cookie:    "",
		UserAgent: "",
		Body:      nil,

		NoCompression: false,
		KeepAlive:     false,
		SkipTLSVerify: false,
		CertPool:      nil,
		ClientCert:    nil,
	}
}

// Typhoon returns a Typhoon followed the configuration
func (conf *Config) Typhoon() (*Typhoon, error) {
	return conf.TyphoonWithHandle(nil)
}

// TyphoonWithHandle returns a Typhoon with RespHandle
func (conf *Config) TyphoonWithHandle(handle RespHandle) (*Typhoon, error) {
	runtime.GOMAXPROCS(conf.NumThread)

	tlsConf := &tls.Config{
		InsecureSkipVerify: conf.SkipTLSVerify,
		RootCAs:            conf.CertPool,
	}
	if conf.ClientCert != nil {
		tlsConf.Certificates = []tls.Certificate{*conf.ClientCert}
	}

	cli := &http.Client{
		Transport: &http.Transport{
			DisableCompression: conf.NoCompression,
			DisableKeepAlives:  !conf.KeepAlive,
			TLSClientConfig:    tlsConf,
		},
	}

	header := make(http.Header)
	if conf.Header != nil {
		header = conf.Header.Clone()
	}
	if conf.UserAgent != "" {
		header.Add("User-Agent", conf.UserAgent)
	}
	if conf.Cookie != "" {
		header.Add("Cookie", conf.Cookie)
	}

	req, err := http.NewRequest(conf.Method, conf.Target, bytes.NewReader(conf.Body))
	if err != nil {
		return nil, err
	}
	req.Header = header
	return NewTyphoon(conf.NumThread, conf.Duration, cli, req, handle), nil
}

// ConfigJSON represents the configuration file format
type ConfigJSON struct {
	Target    string
	NumCPU    int
	NumThread int
	Duration  string

	Method    string
	Header    map[string]string
	Cookie    string
	UserAgent string
	BodyFile  string

	DisableCompression bool
	KeepAlive          bool
	SkipTLSVerify      bool
	ServerCertFile     string
	ClientCertFile     string
	ClientKeyFile      string
}

// ReadConfigJSON returns ConfigJSON which read from file
func ReadConfigJSON(file string) (*ConfigJSON, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cj ConfigJSON
	dec := json.NewDecoder(f)
	if err := dec.Decode(&cj); err != nil {
		return nil, err
	}
	return &cj, nil
}

// Config convert a ConfigJSON to Config
func (cj *ConfigJSON) Config() (*Config, error) {
	duration, err := time.ParseDuration(cj.Duration)
	if err != nil {
		return nil, err
	}
	header := make(http.Header)
	for k, v := range cj.Header {
		header.Add(k, v)
	}

	var body []byte
	if cj.BodyFile != "" {
		body, err = ioutil.ReadFile(cj.BodyFile)
		if err != nil {
			return nil, err
		}
	}

	var pool *x509.CertPool
	if cj.ServerCertFile != "" {
		serverCert, err := ioutil.ReadFile(cj.ServerCertFile)
		if err != nil {
			return nil, err
		}
		pool = x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(serverCert); !ok {
			return nil, fmt.Errorf("failed to append server certificate")
		}
	}

	var clientCert tls.Certificate
	if cj.ClientCertFile != "" && cj.ClientKeyFile != "" {
		cert, err := ioutil.ReadFile(cj.ClientCertFile)
		if err != nil {
			return nil, err
		}
		key, err := ioutil.ReadFile(cj.ClientKeyFile)
		if err != nil {
			return nil, err
		}
		clientCert, err = tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
	}

	conf := DefaultConfig(cj.Target)
	conf.Target = cj.Target
	conf.NumCPU = cj.NumCPU
	conf.NumThread = cj.NumThread
	conf.Duration = duration

	conf.Method = cj.Method
	conf.Header = header
	conf.Cookie = cj.Cookie
	conf.UserAgent = cj.UserAgent
	conf.Body = body

	conf.NoCompression = cj.DisableCompression
	conf.KeepAlive = cj.KeepAlive
	conf.SkipTLSVerify = cj.SkipTLSVerify
	conf.CertPool = pool
	conf.ClientCert = &clientCert
	return conf, nil
}

func (cj ConfigJSON) String() string {
	buf, _ := json.MarshalIndent(&cj, "", "  ")
	return string(buf)
}
