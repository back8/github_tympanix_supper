package list

import (
	"encoding/json"
	"time"

	"github.com/tympanix/supper/types"
)

type LocalMedia struct {
	media []types.LocalMedia
}

func NewLocalMedia(media ...types.LocalMedia) *LocalMedia {
	return &LocalMedia{
		media,
	}
}

func (l *LocalMedia) Add(m types.LocalMedia) {
	l.media = append(l.media, m)
}

func (l *LocalMedia) Len() int {
	return len(l.media)
}

func (l *LocalMedia) List() []types.LocalMedia {
	return l.media
}

func (l *LocalMedia) FilterModified(d time.Duration) types.LocalMediaList {
	t := time.Now().Local().Add(-1 * d)
	media := make([]types.LocalMedia, 0)
	for _, m := range l.List() {
		if m.ModTime().After(t) {
			media = append(media, m)
		}
	}
	return NewLocalMedia(media...)
}

func (l *LocalMedia) FilterVideo() types.VideoList {
	video := make([]types.Video, 0)
	for _, media := range l.List() {
		if v, ok := media.(types.Video); ok {
			video = append(video, v)
		}
	}
	return NewVideo(video...)
}

func (l *LocalMedia) MarshalJSON() (b []byte, err error) {
	return json.Marshal(l.List())
}
