package commander

import (
	"time"
)

// Context stores the argument parameters for the command and additional values middlewares might add
type Context struct {
	Name   string
	Args   []string
	params map[string]interface{}
}

// Arg retrieves a parameter value. If it's not found, returns nil.
func (c Context) Arg(param string) interface{} {
	return c.params[param]
}

// ArgString retrieves a parameter as string. If not found, returns empty string.
func (c Context) ArgString(param string) string {
	intf, ok := c.params[param]
	if !ok {
		return ""
	}

	return intf.(string)
}

// ArgInt retrieves a parameter as integer. If not found, returns 0
func (c Context) ArgInt(param string) int {
	intf, ok := c.params[param]
	if !ok {
		return 0
	}

	return intf.(int)
}

// ArgFloat retrieves a parameter as float64. If not found, returns 0
func (c Context) ArgFloat(param string) float64 {
	intf, ok := c.params[param]
	if !ok {
		return 0
	}

	return intf.(float64)
}

// ArgDuration returns a parameter as time.Duration. If not found, returns 0
func (c Context) ArgDuration(param string) time.Duration {
	intf, ok := c.params[param]
	if !ok {
		return 0
	}

	return intf.(time.Duration)
}

// AddArg Adds a new key-value to the context. Useful for middlewares.
func (c *Context) AddArg(key string, value interface{}) {
	c.params[key] = value
}
