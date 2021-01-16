/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package copyright

// Option is a function that modifies a config.
type Option func(*Copyright)

// OptVerbose sets if we should show verbose output.
func OptVerbose(verbose bool) Option {
	return func(p *Copyright) {
		p.Config.Verbose = &verbose
	}
}

// OptDebug sets if we should show debug output.
func OptDebug(debug bool) Option {
	return func(p *Copyright) {
		p.Config.Debug = &debug
	}
}

// OptExitFirst sets if we should stop after the first failure.
func OptExitFirst(exitFirst bool) Option {
	return func(p *Copyright) {
		p.Config.ExitFirst = &exitFirst
	}
}

// OptRoot sets the root directory to start the profanity check.
func OptRoot(root string) Option {
	return func(p *Copyright) {
		p.Config.Root = root
	}
}

// OptIncludeFiles sets the include glob filter for files.
func OptIncludeFiles(includeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.IncludeFiles = includeGlobs
	}
}

// OptExcludeFiles sets the exclude glob filter for files.
func OptExcludeFiles(excludeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.ExcludeFiles = excludeGlobs
	}
}

// OptIncludeDirs sets the include glob filter for files.
func OptIncludeDirs(includeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.IncludeDirs = includeGlobs
	}
}

// OptExcludeDirs sets the exclude glob filter for directories.
func OptExcludeDirs(excludeGlobs ...string) Option {
	return func(p *Copyright) {
		p.Config.ExcludeDirs = excludeGlobs
	}
}

// OptConfig sets the config in its entirety.
func OptConfig(cfg Config) Option {
	return func(p *Copyright) {
		p.Config = cfg
	}
}
