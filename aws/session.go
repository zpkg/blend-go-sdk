package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// MustNewSession creates a new aws session from a config.
func MustNewSession(cfg *Config) *session.Session {
	if cfg.IsZero() {
		return session.Must(session.NewSession())
	}

	awsConfig := &aws.Config{
		Region:      aws.String(cfg.GetRegion()),
		Credentials: credentials.NewStaticCredentials(cfg.GetAccessKeyID(), cfg.GetSecretAccessKey(), cfg.GetToken()),
	}
	return session.Must(session.NewSession(awsConfig))
}
