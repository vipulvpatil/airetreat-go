package config

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type Config struct {
	EnableTls        bool
	RedisUrl         string
	TestDbUrl        string
	DbUrl            string
	CaCertBase64     string
	ServerCertBase64 string
	ServerKeyBase64  string
	AllowUnauthed    bool
}

func envVarLoaderBool(envVarName string, required bool, errorCollector *[]error) bool {
	value, ok := os.LookupEnv(envVarName)
	if !ok && required {
		*errorCollector = append(*errorCollector, errors.Errorf("%s is a required Env var", envVarName))
		return false
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		*errorCollector = append(*errorCollector, errors.Errorf("Env var %s is expected to be a boolean", envVarName))
		return false
	}
	return boolValue
}

func envVarLoaderString(envVarName string, required bool, errorCollector *[]error) string {
	value, ok := os.LookupEnv(envVarName)
	if !ok && required {
		*errorCollector = append(*errorCollector, errors.Errorf("%s is a required Env var", envVarName))
		return ""
	}
	return value
}

func NewConfigFromEnvVars() (*Config, []error) {
	c := Config{}

	errs := []error{}

	c.EnableTls = envVarLoaderBool("ENABLE_TLS", true, &errs)
	c.RedisUrl = envVarLoaderString("REDIS_URL", true, &errs)
	c.TestDbUrl = envVarLoaderString("TEST_DB_URL", false, &errs)
	c.DbUrl = envVarLoaderString("DB_URL", true, &errs)
	c.CaCertBase64 = envVarLoaderString("CA_CERT_BASE64", true, &errs)
	c.ServerCertBase64 = envVarLoaderString("SERVER_CERT_BASE64", true, &errs)
	c.ServerKeyBase64 = envVarLoaderString("SERVER_KEY_BASE64", true, &errs)
	c.AllowUnauthed = envVarLoaderBool("ALLOW_UNAUTHED", true, &errs)

	return &c, errs
}
