// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

import (
	"github.com/letscool/lc-go/lg"
	"io"
	"os"
	"strings"
)

const (
	DefaultWebServerPort    = 8000
	DefaultValidatorTagName = "validate"
)

type RuntimeEnv string

const (
	RtProd  RuntimeEnv = "prod"
	RtDev   RuntimeEnv = "dev"
	RtDebug RuntimeEnv = "debug"
	RtTest  RuntimeEnv = "test"
)

func (r RuntimeEnv) IsInOnlyEnv(onlyEnv string) bool {
	if onlyEnv = strings.TrimSpace(onlyEnv); len(onlyEnv) == 0 {
		return true
	}
	return lg.Contains(strings.Split(onlyEnv, "&"), string(r))
}

type WebConfig struct {
	Address          string // default is localhost
	Port             int    // if not setting, 8000 will be used.
	MaxConn          int
	BasePath         string // Default is empty
	ValidatorTagName string // Default is "validate", but go-gin preferred "binding".
	DependentRefs    map[string]any
	RtEnv            RuntimeEnv // Default is "dev"
	OpenApi          OpenApiConfig
	DefaultWriter    io.Writer
}

// NewWebConfig returns an instance with default values.
func NewWebConfig() *WebConfig {
	config := &WebConfig{}
	config.ApplyDefaultValues()
	return config
}

// ApplyDefaultValues if some properties are zero or empty, it will set the default values.
func (c *WebConfig) ApplyDefaultValues() {
	if c.Address == "" {
		c.Address = "localhost"
	}
	if c.Port == 0 {
		c.Port = DefaultWebServerPort
	}
	if c.ValidatorTagName == "" {
		c.ValidatorTagName = DefaultValidatorTagName
	}
	if c.RtEnv == "" {
		c.RtEnv = RtDev
	}
	if c.DependentRefs == nil {
		c.DependentRefs = map[string]any{}
	}
	if c.DefaultWriter == nil {
		c.DefaultWriter = os.Stdout
	}
	c.OpenApi.ApplyDefaultValues()
}

func (c *WebConfig) SetRtMode(mode RuntimeEnv) *WebConfig {
	c.RtEnv = mode
	return c
}

func (c *WebConfig) SetAddress(addr string) *WebConfig {
	c.Address = addr
	return c
}

func (c *WebConfig) SetPort(port int) *WebConfig {
	c.Port = port
	return c
}

func (c *WebConfig) SetBasePath(path string) *WebConfig {
	c.BasePath = path
	return c
}

func (c *WebConfig) SetOpenApi(f func(o *OpenApiConfig)) *WebConfig {
	f(&c.OpenApi)
	return c
}

func (c *WebConfig) SetDependentRef(key string, ref any) *WebConfig {
	if c.DependentRefs == nil {
		c.DependentRefs = map[string]any{}
	}
	c.DependentRefs[key] = ref
	return c
}

//func (c *WebConfig) SetLogFormatter(formatter gin.LogFormatter) *WebConfig {
//	return c.SetDependentRef("_.mdl.log.formatter", formatter)
//}

func (c *WebConfig) SetDefaultWriter(writer io.Writer) *WebConfig {
	c.DefaultWriter = writer
	return c
}

type OpenApiConfig struct {
	Enabled     bool // Default is false
	Title       string
	Description string
	Version     string
	Schemes     []string // ex: "http", "https". Default is "https".
	DocPath     string   // Default is "doc"
}

func (o *OpenApiConfig) ApplyDefaultValues() {
	if o.Schemes == nil || len(o.Schemes) == 0 {
		o.Schemes = []string{"https"}
	}
	if o.DocPath == "" {
		o.DocPath = "doc"
	}
}

func (o *OpenApiConfig) UseHttpOnly() *OpenApiConfig {
	o.Schemes = []string{"http"}
	return o
}

func (o *OpenApiConfig) UseHttpsOnly() *OpenApiConfig {
	o.Schemes = []string{"https"}
	return o
}

func (o *OpenApiConfig) UseHttpAndHttps() *OpenApiConfig {
	o.Schemes = []string{"https", "http"}
	return o
}

func (o *OpenApiConfig) Enable() *OpenApiConfig {
	o.Enabled = true
	return o
}

func (o *OpenApiConfig) SetEnabled(en bool) *OpenApiConfig {
	o.Enabled = en
	return o
}

func (o *OpenApiConfig) SetDocPath(path string) *OpenApiConfig {
	o.DocPath = path
	return o
}

func (o *OpenApiConfig) SetTitle(title string) *OpenApiConfig {
	o.Title = title
	return o
}

func (o *OpenApiConfig) SetDescription(description string) *OpenApiConfig {
	o.Description = description
	return o
}

func (o *OpenApiConfig) SetVersion(version string) *OpenApiConfig {
	o.Version = version
	return o
}
