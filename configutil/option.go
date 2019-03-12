package configutil

// Option is a modification of config options.
type Option func(*ConfigOptions) error

// OptAddPaths adds paths to the options
func OptAddPaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.Paths = append(co.Paths, paths...)
		return nil
	}
}

// OptSetPaths adds paths to the options
func OptSetPaths(paths ...string) Option {
	return func(co *ConfigOptions) error {
		co.Paths = paths
		return nil
	}
}

// OptResolver sets an additional resolver for the config read.
func OptResolver(resolver func(interface{}) error) Option {
	return func(co *ConfigOptions) error {
		co.Resolver = resolver
		return nil
	}
}
