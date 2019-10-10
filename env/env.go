package env

import (
	"os"
	"strings"
)

var properties = []*EnvProperty{
	{Key: "DATABASE_HOST", FallBackValue: "go-mongo"},
	{Key: "DATABASE_PORT", FallBackValue: "27017"},
	{Key: "FILE_SOURCE", FallBackValue: "q1_catalog.csv"},
	{Key: "FILE_SOURCE_SEPARATOR", FallBackValue: ";"},
}

type EnvProperty struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	FallBackValue string `json:"fallbackValue"`
}

type EnvProperties struct {
	Properties []*EnvProperty
}

func envProperties() *EnvProperties {
	objProperties := &EnvProperties{}
	objProperties.Properties = properties
	return objProperties
}

func DatabaseHost() string {
	return envProperties().envString("DATABASE_HOST")
}

func DatabasePort() string {
	return envProperties().envString("DATABASE_PORT")
}

func FileSource() string {
	return envProperties().envString("FILE_SOURCE")
}

func FileSourceSeparator() rune {
	return []rune(envProperties().envString("FILE_SOURCE_SEPARATOR"))[0]
}

func (e *EnvProperties) getEnvProperty(env string) *EnvProperty {
	for _, property := range e.Properties {
		if strings.ToUpper(property.Key) == strings.ToUpper(env) {
			return property
		}
	}
	return nil
}

func (e *EnvProperties) envString(env string) string {
	envProperty := e.getEnvProperty(env)
	if envProperty != nil {
		val := os.Getenv(env)
		if val == "" {
			return envProperty.FallBackValue
		}
		return val
	}
	return ""
}