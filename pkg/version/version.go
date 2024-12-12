package version

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

var (
	gitTag       string                                      // tag of the release
	gitCommit    string                                      // commit hash
	gitBranch    string                                      // branch name
	gitTreeState string                                      // tree state (clean or dirty)
	version      = "dev"                                     // version of the release
	buildDate    = time.Now().Format("2006-01-02T15:04:05Z") // build date
)

// Info contains version information.
type Info struct {
	GitTag       string `json:"gitTag"`
	GitCommit    string `json:"gitCommit"`
	GitBranch    string `json:"gitBranch"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
	Version      string `json:"version"`
}

// GetInfo returns version information.
func GetInfo() Info {
	return Info{
		GitTag:       gitTag,
		GitCommit:    gitCommit,
		GitBranch:    gitBranch,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Version:      version,
	}
}

// GetInfoString returns version information as a string.
func GetInfoString() string {
	v := GetInfo()
	marshalled, err := json.MarshalIndent(&v, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	return string(marshalled)
}
