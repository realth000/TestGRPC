package common

import (
	"flag"
	"log"
	"reflect"
)

// ConfMap use command line name as map key and a pointer points to its value as map value.
// Used to store configs from a config file.
type ConfMap map[string]interface{}

func MakeConfMap(conf *ConfMap, c interface{}) {
	tc := reflect.TypeOf(c)
	vc := reflect.ValueOf(c)
	for f := 0; f < tc.NumField(); f++ {
		ff := vc.FieldByName(tc.Field(f).Name)
		switch ff.Kind().String() {
		case "int":
			tmp := ff.Int()
			(*conf)[tc.Field(f).Tag.Get("name")] = &tmp
		case "uint":
			tmp := ff.Uint()
			(*conf)[tc.Field(f).Tag.Get("name")] = &tmp
		case "string":
			tmp := ff.String()
			(*conf)[tc.Field(f).Tag.Get("name")] = &tmp
		case "bool":
			tmp := ff.Bool()
			(*conf)[tc.Field(f).Tag.Get("name")] = &tmp
		}
	}
}

type Config struct {
	Variable     interface{}
	Name         string
	DefaultValue interface{}
	Usage        string
	CmdValue     interface{}
	FileValue    interface{}
	Override     bool
	registered   bool
}

// RegisterFlag is not usable, FIXME: Types always not match
func (c *Config) RegisterFlag() {
	if c.registered {
		return
	}
	var vt string
	switch t := c.Variable.(type) {
	case *int:
		vt = "int"
		p, _ := c.Variable.(*int)
		dp, ok := c.DefaultValue.(int)
		if !ok {
			goto notMatchType
		}
		flag.IntVar(p, c.Name, dp, c.Usage)
	case *uint:
		vt = "uint"
		p, _ := c.Variable.(*uint)
		dp, ok := c.DefaultValue.(uint)
		if !ok {
			goto notMatchType
		}
		flag.UintVar(p, c.Name, dp, c.Usage)
	case *string:
		vt = "string"
		p, _ := c.Variable.(*string)
		dp, ok := c.DefaultValue.(string)
		if !ok {
			goto notMatchType
		}
		flag.StringVar(p, c.Name, dp, c.Usage)
	case *bool:
		vt = "bool"
		p, _ := c.Variable.(*bool)
		dp, ok := c.DefaultValue.(bool)
		if !ok {
			goto notMatchType
		}
		flag.BoolVar(p, c.Name, dp, c.Usage)
	default:
		log.Fatalf("Unsupported flag type :%s", t)
	}

notMatchType:
	log.Fatalf("error defining flag %s, defaultValue type is not %s", c.Name, vt)
}

func NewFlag(variable interface{}, name string, defaultValue interface{}, usage string) *Config {
	var (
		c  = new(Config)
		vt string
	)
	switch t := variable.(type) {
	case *int:
		vt = "int"
		p, _ := variable.(*int)
		dp, ok := defaultValue.(int)
		if !ok {
			goto notMatchType
		}
		flag.IntVar(p, name, dp, usage)
	case *uint:
		vt = "uint"
		p, _ := variable.(*uint)
		dp, ok := defaultValue.(uint)
		if !ok {
			goto notMatchType
		}
		flag.UintVar(p, name, dp, usage)
	case *string:
		vt = "string"
		p, _ := variable.(*string)
		dp, ok := defaultValue.(string)
		if !ok {
			goto notMatchType
		}
		flag.StringVar(p, name, dp, usage)
	case *bool:
		vt = "bool"
		p, _ := variable.(*bool)
		dp, ok := defaultValue.(bool)
		if !ok {
			goto notMatchType
		}
		flag.BoolVar(p, name, dp, usage)
	default:
		log.Fatalf("Unsupported flag type :%s", t)
	}
	c.Variable = variable
	c.Name = name
	c.DefaultValue = defaultValue
	c.Usage = usage
	c.Override = false
	c.registered = true
	return c

notMatchType:
	log.Fatalf("error defining flag %s, defaultValue type is not %s", name, vt)
	return nil
}
