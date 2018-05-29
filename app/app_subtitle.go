package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/fatih/set"
	"github.com/tympanix/supper/list"
	"github.com/tympanix/supper/types"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// DownloadSubtitles downloads subtitles for a whole list of mediafiles for every
// langauge in the language set
func (a *Application) DownloadSubtitles(media types.LocalMediaList, lang set.Interface) ([]types.LocalSubtitle, error) {
	var result []types.LocalSubtitle

	if media == nil {
		return nil, errors.New("no media supplied for subtitles")
	}

	if lang == nil {
		return nil, errors.New("no languages supplied for subtitles")
	}

	video := media.FilterVideo()

	if video.Len() == 0 {
		return nil, errors.New("no video media found in path")
	}

	video, err := video.FilterMissingSubs(lang)

	if err != nil {
		return nil, err
	}

	// Iterate all media files in the list
	for i, item := range video.List() {
		ctx := log.WithFields(log.Fields{
			"media": item,
			"item":  fmt.Sprintf("%v/%v", i+1, video.Len()),
		})

		cursubs, err := item.ExistingSubtitles()

		if err != nil {
			return nil, err
		}

		missingLangs := set.Difference(lang, cursubs.LanguageSet())

		if missingLangs.Size() == 0 {
			continue
		}

		var subs = list.Subtitles()

		if !a.Config().Dry() {
			var search []types.OnlineSubtitle
			search, err = a.SearchSubtitles(item)
			if err != nil {
				ctx.WithError(err).Error("Subtitle failed")
				if a.Config().Strict() {
					return nil, err
				}
				continue
			}
			subs, err = list.NewSubtitlesFromInterface(search)
			if err != nil {
				ctx.WithError(err).Fatal("Subtitle error")
			}
		}

		subs = subs.HearingImpaired(a.Config().Impaired())

		// Download subtitle for each language
		for _, l := range missingLangs.List() {
			ctx = ctx.WithField("lang", display.English.Languages().Name(l))

			if a.Config().Delay() > 0 {
				time.Sleep(a.Config().Delay())
			}

			l, ok := l.(language.Tag)

			if !ok {
				return nil, err
			}

			langsubs := subs.FilterLanguage(l)

			if langsubs.Len() == 0 && !a.Config().Dry() {
				ctx.Warn("No subtitle available")
				continue
			}

			if !a.Config().Dry() {
				sub, err := a.downloadBestSubtitle(ctx, item, langsubs)
				if err != nil {
					if a.Config().Strict() {
						return nil, err
					}
					continue
				}
				result = append(result, sub)
			} else {
				ctx.WithField("reason", "dry-run").Info("Skip download")
			}
		}
	}
	return result, nil
}

func (a *Application) downloadBestSubtitle(ctx log.Interface, m types.Video, l types.SubtitleList) (types.LocalSubtitle, error) {
	rated := l.RateByMedia(m, a.Config().Evaluator())
	if rated.Len() == 0 {
		m := "no subtitles satisfied media"
		ctx.Warn(m)
		return nil, errors.New(m)
	}
	sub := rated.Best()
	if sub.Score() < (float32(a.Config().Score()) / 100.0) {
		m := fmt.Sprintf("Score too low %.0f%%", sub.Score()*100.0)
		ctx.Warnf(m)
		return nil, errors.New(m)
	}
	onl, ok := sub.Subtitle().(types.OnlineSubtitle)
	if !ok {
		ctx.Fatal("Subtitle could not be cast to online subtitle")
	}
	srt, err := onl.Download()
	if err != nil {
		ctx.WithError(err).Error("Could not download subtitle")
		return nil, err
	}
	defer srt.Close()
	saved, err := m.SaveSubtitle(srt, onl.Language())
	if err != nil {
		ctx.WithError(err).Error("Subtitle error")
		return nil, err
	}

	var strscore string
	if sub.Score() == 0.0 {
		strscore = "N/A"
	} else {
		strscore = fmt.Sprintf("%.0f%%", sub.Score()*100.0)
	}

	ctx.WithField("score", strscore).Info("Subtitle downloaded")

	if err := a.execPluginsOnSubtitle(ctx, saved); err != nil {
		return nil, err
	}
	return saved, nil
}

func (a *Application) execPluginsOnSubtitle(ctx log.Interface, s types.LocalSubtitle) error {
	for _, plugin := range a.Config().Plugins() {
		err := plugin.Run(s)
		if err != nil {
			ctx.WithField("plugin", plugin.Name()).Error("Plugin failed")
			return err
		}
		ctx.WithField("plugin", plugin.Name()).Info("Plugin finished")
	}
	return nil
}
