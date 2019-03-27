package db

// Option is an option for database connections.
type Option func(c *Connection) error

// OptConfig sets the config on a connection.
func OptConfig(cfg *Config) Option {
	return func(c *Connection) error {
		c.Config = cfg
		return nil
	}
}

// OptConfigFromEnv sets the config on a connection from the environment.
func OptConfigFromEnv() Option {
	return func(c *Connection) error {
		cfg, err := NewConfigFromEnv()
		if err != nil {
			return err
		}
		c.Config = cfg
		return nil
	}
}
