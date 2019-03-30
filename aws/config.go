package aws

import (
	"github.com/blend/go-sdk/configutil"
)

const (
	// DefaultAWSRegion is a default.
	DefaultAWSRegion = "us-east-1"
)

// Config is a config object.
type Config struct {
	Region          string `json:"region,omitempty" yaml:"region,omitempty" env:"AWS_REGION"`
	AccessKeyID     string `json:"accessKeyID,omitempty" yaml:"accessKeyID,omitempty" env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `json:"secretAccessKey,omitempty" yaml:"secretAccessKey,omitempty" env:"AWS_SECRET_ACCESS_KEY"`
	SecurityToken   string `json:"securityToken,omitempty" yaml:"securityToken,omitempty" env:"AWS_SECURITY_TOKEN"`
}

// Resolve adds extra resolution steps for the config.
func (a *Config) Resolve() error {
	return configutil.AnyError(
		configutil.SetString(&a.Region, configutil.String(a.Region), configutil.Env("AWS_REGION"), configutil.String(DefaultAWSRegion)),
		configutil.SetString(&a.AccessKeyID, configutil.String(a.AccessKeyID), configutil.Env("AWS_ACCESS_KEY_ID")),
		configutil.SetString(&a.SecretAccessKey, configutil.String(a.SecretAccessKey), configutil.Env("AWS_SECRET_ACCESS_KEY")),
		configutil.SetString(&a.SecurityToken, configutil.String(a.SecurityToken), configutil.Env("AWS_SECURITY_TOKEN")),
	)
}

// IsZero returns if the config is unset or not.
func (a Config) IsZero() bool {
	return len(a.AccessKeyID) == 0 || len(a.SecretAccessKey) == 0
}
