// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

type WebConfig struct {
	Address  string // default is empty
	Port     int    // if not setting, 8000 will be used.
	MaxConn  int
	BasePath string   // Default is empty
	Schemes  []string // ex: "http", "https". Default is "https".
}

// NewWebConfig returns an instance with default values.
func NewWebConfig() *WebConfig {
	config := &WebConfig{}
	config.ApplyDefaultValues()
	return config
}

// ApplyDefaultValues if some properties are zero or empty, it will set the default values.
func (c *WebConfig) ApplyDefaultValues() {
	//if c.Address == "" {
	//	c.Address = "localhost"
	//}
	if c.Port == 0 {
		c.Port = DefaultWebServerPort
	}
	if c.Schemes == nil || len(c.Schemes) == 0 {
		c.Schemes = []string{"https"}
	}
}

func (c *WebConfig) UseHttpOnly() *WebConfig {
	c.Schemes = []string{"http"}
	return c
}

func (c *WebConfig) UseHttpsOnly() *WebConfig {
	c.Schemes = []string{"https"}
	return c
}

func (c *WebConfig) UseHttpAndHttps() *WebConfig {
	c.Schemes = []string{"https", "http"}
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
