package secrets

// NewFromConfig is an alias to NewVaultClientFromConfig.
func NewFromConfig(cfg *Config) (*VaultClient, error) {
	return NewVaultClientFromConfig(cfg)
}

// NewFromEnv is an alias to NewVaultClientFromConfig with the config set by NewConfigFromEnv.
func NewFromEnv() (*VaultClient, error) {
	cfg, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewVaultClientFromConfig(cfg)
}
