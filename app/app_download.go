package application

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/fatih/set"
	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/types"
	"github.com/urfave/cli"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var green = color.New(color.FgGreen)
var red = color.New(color.FgRed)

func (a *Application) DownloadSubtitles(media types.LocalMediaList, lang set.Interface, out io.Writer) error {
	numsubs := 0

	dry := a.Context().GlobalBool("dry")
	hi := a.Context().GlobalBool("impaired")

	// Iterate all media files found in each path
	for i, item := range media.List() {
		cursubs, err := item.ExistingSubtitles()

		if err != nil {
			return cli.NewExitError(err, 2)
		}

		missingLangs := set.Difference(lang, cursubs.LanguageSet())

		if missingLangs.Size() == 0 {
			continue
		}

		fmt.Fprintf(out, "(%v/%v) - %s\n", i+1, media.Len(), item)

		subs := list.Subtitles()

		if !dry {
			subs, err = a.SearchSubtitles(item)
			if err != nil {
				return cli.NewExitError(err, 2)
			}
		}

		subs = subs.HearingImpaired(hi)

		// Download subtitle for each language
		for _, l := range missingLangs.List() {
			l, ok := l.(language.Tag)

			numsubs++

			if !ok {
				return cli.NewExitError(err, 3)
			}

			langsubs := subs.FilterLanguage(l)

			if langsubs.Len() == 0 && !dry {
				red.Fprintln(out, " - no subtitles found")
				continue
			}

			if !dry {
				err := item.SaveSubtitle(langsubs.Best())
				if err != nil {
					red.Fprintln(out, err.Error())
					continue
				}
			}

			green.Fprintf(out, " - %v\n", display.English.Languages().Name(l))
		}
	}
	return nil
}