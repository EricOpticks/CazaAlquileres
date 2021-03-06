package main

import "github.com/spf13/viper"

// Provider provides the configuration settings for Hugo.
type Provider interface {
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	Get(key string) interface{}
	Set(key string, value interface{})
	IsSet(key string) bool
	Unmarshal(v interface{}, opts ...viper.DecoderConfigOption) error
}
