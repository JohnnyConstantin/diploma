package version

import (
	"fmt"
	"runtime"
)

var (
	// Version версия приложения
	Version = "0.0.1"
	// BuildDate дата сборки
	BuildDate = "unknown"
	// GitCommit хеш коммита
	GitCommit = "unknown"
	// GoVersion версия Go
	GoVersion = runtime.Version()
)

// Info содержит информацию о версии
type Info struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Get возвращает информацию о версии
func Get() Info {
	return Info{
		Version:   Version,
		BuildDate: BuildDate,
		GitCommit: GitCommit,
		GoVersion: GoVersion,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
