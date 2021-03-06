package media

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/tympanix/supper/parse"
	"github.com/tympanix/supper/types"
)

// Movie represents a movie file
type Movie struct {
	Metadata
	TypeNone
	NameX string
	YearX int
	tags  string
}

// MarshalJSON return the JSON representation of a movie
func (m *Movie) MarshalJSON() (b []byte, err error) {
	type jsonMovie struct {
		Meta Metadata `json:"metadata"`
		Name string   `json:"name"`
		Year int      `json:"year"`
	}

	return json.Marshal(jsonMovie{
		m.Metadata,
		m.NameX,
		m.YearX,
	})
}

var movieRegexp = regexp.MustCompile(`^(.+)[\W_]+(19\d\d|20\d\d)[\W_]+(.*)$`)

// NewMovie parses media info from a filename (without extension). The filename
// must describe the movie adequately (e.g. must contain the release year)
func NewMovie(filename string) (*Movie, error) {
	groups := movieRegexp.FindStringSubmatch(filename)

	if groups == nil {
		return nil, errors.New("could not parse movie")
	}

	name := groups[1]
	year, err := strconv.Atoi(groups[2])
	tags := groups[3]

	if err != nil {
		return nil, err
	}

	return &Movie{
		Metadata: ParseMetadata(tags),
		NameX:    parse.CleanName(name),
		tags:     tags,
		YearX:    year,
	}, nil
}

func (m *Movie) String() string {
	return fmt.Sprintf("%s (%v)", m.MovieName(), m.Year())
}

// Merge merges metadata from another movie into this one
func (m *Movie) Merge(other types.Media) error {
	if movie, ok := other.TypeMovie(); ok {
		if m.Year() != movie.Year() {
			return errors.New("invalid media merge year does not match")
		}
		m.NameX = movie.MovieName()
		return nil
	}
	return errors.New("invalid media merge not same media type")
}

// Meta returnes the metadata interface for a movie
func (m *Movie) Meta() types.Metadata {
	return m.Metadata
}

// TypeMovie returns true, since a movie is a movie
func (m *Movie) TypeMovie() (types.Movie, bool) {
	return m, true
}

// MovieName is the name of the movie
func (m *Movie) MovieName() string {
	return m.NameX
}

// Year is the release year of the movie
func (m *Movie) Year() int {
	return m.YearX
}
