package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GalassoX/social-downloader-api/downloaders"
)

type Map map[string]interface{}

func SendJSON(w http.ResponseWriter, body interface{}) {
	data, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	val := r.URL.Query().Get("val")
	fmt.Println(val)
	SendJSON(w, Map{"message": "Hello World", "val": val})
	// w.Write([]byte("{\"message\":\"Hello World!\"}"))
}

func FBHandler(w http.ResponseWriter, r *http.Request) {
	videoUrl := r.URL.Query().Get("url")

	video, _ := downloaders.GetFBVideo(videoUrl)
	SendJSON(w, video)
}

func YoutubeHandler(w http.ResponseWriter, r *http.Request) {
	videoURL := r.URL.Query().Get("url")
	videos := downloaders.GetYoutubeURL(videoURL)
	SendJSON(w, videos)
}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/facebook", FBHandler)
	http.HandleFunc("/youtube", YoutubeHandler)
	http.ListenAndServe(":8080", nil)
}
