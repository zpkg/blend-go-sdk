package airbrake

// Config is the airbrake config.
type Config struct {
	ProjectID   string `json:"projectID" yaml:"projectID" env:"AIRBRAKE_PROJECT_ID"`
	ProjectKey  string `json:"projectKey" yaml:"projectKey" env:"AIRBRAKE_PROJECT_KEY"`
	Environment string `json:"environment" yaml:"environment" env:"SERVICE_ENV"`
}

// IsZero returns if the config is set or not.
func (c Config) IsZero() bool {
	return len(c.ProjectKey) == 0 || len(c.ProjectID) == 0
}
