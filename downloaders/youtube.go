package downloaders

import (
	"github.com/kkdai/youtube/v2"
)

func GetYoutubeURL(videoURL string) map[string]string {
	client := youtube.Client{}

	video, err := client.GetVideo(videoURL)

	if err != nil {
		panic(err)
	}
	results := map[string]string{}
	for i := range video.Formats.WithAudioChannels() {
		data := video.Formats.WithAudioChannels()[i]
		key := data.QualityLabel
		if len(key) < 1 {
			key = data.Quality
		}
		results[key] = data.URL
	}
	return results
}
