package version

import "fmt"

// This value is set at release build time
// ldflags "-X github.com/clover0/issue-agent/cli/command/version.version=1.0.0)"
var version = "development"

const VersionCommand = "version"

func Version() error {
	fmt.Printf("Version: %s\n", version)
	return nil
}
