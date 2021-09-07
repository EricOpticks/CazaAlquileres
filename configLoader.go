package main

import (
	"flag"
	"fmt"
	"github.com/juju/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"strings"
)

type AConfig struct {
	v          *viper.Viper
	properties map[string]*Property
	withFile   bool
	loaded     bool
}

type Property struct {
	name     string
	alias    string
	env      string
	required bool
}

func NewConfig() *AConfig {
	newConf := AConfig{}
	newConf.v = viper.New()
	newConf.properties = map[string]*Property{}

	return &newConf
}

func (c *AConfig) Load(print bool) (Provider, error) {

	if c.loaded {
		return nil, errors.New("config was already loaded")
	} else {
		for name, p := range c.properties {
			// register alias
			if p.alias != "" {
				c.v.RegisterAlias(p.alias, name)
			}

			// register env name
			if p.env != "" {
				_ = c.v.BindEnv(name, p.env)
			}
		}

		if c.withFile {
			if err := c.v.ReadInConfig(); err != nil {
				return nil, errors.New(fmt.Sprintf("error reading config file - %s", err))
			}
		}

		c.v.SetEnvPrefix("") // will be uppercased automatically
		c.v.AutomaticEnv()

		for name, _ := range c.properties {
			flag.String(name, c.v.GetString(name), name+" is required")
		}

		pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
		pflag.Parse()

		if err := c.v.BindPFlags(pflag.CommandLine); err != nil {
			return nil, errors.New(fmt.Sprintf("error binding flags - %s", err))
		}

		for name, p := range c.properties {
			if p.required {
				if c.v.GetString(name) == "" {
					return nil, errors.New(name + " config property is required")
				}
			}
		}
	}

	if print {
		c.Print()
	}

	return c.v, nil
}

func (c *AConfig) Print() {

	builder := strings.Builder{}

	builder.WriteString("LOADED CONFIG:\n")
	for name, p := range c.properties {
		builder.WriteString("        ")
		propName := strings.ToUpper(name)
		if p.env != "" {
			propName = p.env
		}
		propName = strings.Replace(propName, "-", "_", -1)
		propName = strings.Replace(propName, " ", "_", -1)

		for len(propName) < 26 {
			propName += " "
		}

		builder.WriteString(propName)
		builder.WriteString("= ")
		builder.WriteString(c.v.GetString(name))
		builder.WriteString("\n")
	}

	fmt.Println(builder.String())
}

func (c *AConfig) WithProperty(name string, required bool) *Property {
	prop := Property{}
	prop.name = name
	prop.required = required
	c.properties[name] = &prop

	return c.properties[name]
}

func (c *AConfig) SetFilePath(path string) {
	c.v.AddConfigPath(path)

}

func (c *AConfig) SetFileName(name string) {
	c.v.SetConfigName(name)
	c.withFile = true
}

func (p *Property) Alias(alias string) *Property {
	p.alias = alias
	return p
}

func (p *Property) EnvAlias(env string) *Property {
	p.env = env
	return p
}
