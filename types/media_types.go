package types

import (
	"io"
	"os"

	"github.com/tympanix/supper/meta/codec"
	"github.com/tympanix/supper/meta/quality"
	"github.com/tympanix/supper/meta/source"
	"golang.org/x/text/language"
)

// Provider interfaces with subtitles websites to provide subtitles
type Provider interface {
	SearchSubtitles(LocalMedia) ([]OnlineSubtitle, error)
	ResolveSubtitle(Linker) (Downloadable, error)
}

// Scraper interfaces with 3rd party APIs to scrape meta data
type Scraper interface {
	Scrape(Media) (Media, error)
}

// Downloadable is an interface for media that can be downloaded from the internet
type Downloadable interface {
	Download() (io.ReadCloser, error)
}

// Evaluator determines how well two media types are alike
type Evaluator interface {
	Evaluate(Media, Media) float32
}

// Media is an interface for movies and TV shows
type Media interface {
	Meta() Metadata
	Merge(Media) error
	String() string
	TypeMovie() (Movie, bool)
	TypeEpisode() (Episode, bool)
}

// Metadata is an interface metadata information
type Metadata interface {
	Group() string
	Codec() codec.Tag
	Quality() quality.Tag
	Source() source.Tag
	AllTags() []string
}

// MediaReadCloser is an interface for media which need to be closed
type MediaReadCloser interface {
	Media
	io.ReadCloser
}

// MediaArchive is an interface for media sourced from archives
type MediaArchive interface {
	Next() (MediaReadCloser, error)
	io.Closer
}

// Opener is a interface for devices which can be opened for reading
type Opener interface {
	Open() (io.ReadCloser, error)
}

// Local is an interface for media content which is stored on disk
type Local interface {
	os.FileInfo
	Pather
	Opener
}

// Pather is an interface for resources which can be accessed in a local path
type Pather interface {
	Path() string
}

// LocalMedia is an interface for media found locally on disk
type LocalMedia interface {
	Local
	Media
	ExistingSubtitles() (SubtitleList, error)
	SaveSubtitle(Downloadable, language.Tag) (LocalSubtitle, error)
}

// Movie interface is for movie type media material
type Movie interface {
	Metadata
	MovieName() string
	Year() int
}

// Episode interface is for TV show type material
type Episode interface {
	Metadata
	TVShow() string
	EpisodeName() string
	Episode() int
	Season() int
}

// Linker is an object which can be fetched from the internet
type Linker interface {
	Link() string
}

// Subtitle can be downloaded
type Subtitle interface {
	Media
	Language() language.Tag
	IsLang(language.Tag) bool
	IsHI() bool
}

// LocalSubtitle is an subtitle which is stored on disk
type LocalSubtitle interface {
	Local
	Subtitle
}

// OnlineSubtitle is a subtitle obtained from the internet
// and can be downloaded and stored on disk
type OnlineSubtitle interface {
	Linker
	Downloadable
	Subtitle
}
