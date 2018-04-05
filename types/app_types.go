package types

import (
	"net/http"

	"github.com/fatih/set"
	"github.com/urfave/cli"
)

// App is the interface for the top level capabilities of the application.
// It is an HTTP handler, a provider (for subtitles) and a CLI application.
// It means App can both be used as a HTTP server and a CLI application.
type App interface {
	Provider
	http.Handler
	Plugins() []Plugin
	Context() *cli.Context
	FindMedia(...string) (LocalMediaList, error)
	Languages() set.Interface
	DownloadSubtitles(LocalMediaList, set.Interface) (int, error)
}

// Plugin is an interface for external functionality
type Plugin interface {
	Name() string
	Run(LocalSubtitle) error
}
