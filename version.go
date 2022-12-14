package main

import "runtime"

var (
	semanticVersion string
	gitCommit       string
)

// BuildInfo describes compile time information.
type BuildInfo struct {
	// Version is the current semanticVersion.
	Version string

	// GitCommit is the git sha1.
	GitCommit string

	// GoVersion is the version of the Go compiler used.
	GoVersion string
}

func Get() BuildInfo {
	return BuildInfo{
		Version:   semanticVersion,
		GitCommit: gitCommit,
		GoVersion: runtime.Version(),
	}
}
