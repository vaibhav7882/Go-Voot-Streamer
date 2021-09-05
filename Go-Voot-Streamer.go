package main

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"html/template"
	"io"
	"net/http"
	"strings"
)

func getStreamingUrl(url string) (string, string, string) {
	videoID := url[strings.LastIndex(url, "/")+1:]

	webURL := "https://wapi.voot.com/ws/ott/getMediaInfo.json?platform=Web&pId=2&mediaId=" + videoID
	httpClient := http.Client{}

	requests, _ := http.NewRequest("GET", webURL, nil)

	result, _ := httpClient.Do(requests)

	response, _ := io.ReadAll(result.Body)

	jsonParsed, _ := gabs.ParseJSON(response)

	if jsonParsed.Path("status.message").Data().(string) != "Data not found or Invalid MediaID" {

		streamURL := jsonParsed.Path("assets.Files.3.URL").Data().(string)
		title := jsonParsed.Path("assets.MediaName").Data().(string)
		thumbnail := jsonParsed.Path("assets.Pictures.0.URL").Data().(string)
		thumbnail = strings.ReplaceAll(thumbnail, "https://viacom18-res.cloudinary.com/image/upload/f_auto,q_auto:eco,fl_lossy/kimg", "https://kimg.voot.com")
		return title, thumbnail, streamURL
	} else {
		return "", "", ""
	}
}

func welcomeScreen() {
	fmt.Println("Go Voot Streamer\n" +
		"Developed By Henry Richard J")
	fmt.Println("Example: http://localhost:8080/player?url=https://www.voot.com/movies/petta/962823")
}

type streamingLink struct {
	Title     string
	LinkUrl   string
	PosterUrl string
}

func handlePlayer(w http.ResponseWriter, r *http.Request) {
	vootURL := r.URL.Query().Get("url")
	title, thumbNail, steamUrl := getStreamingUrl(vootURL)
	if title != "" {

		w.WriteHeader(http.StatusOK)
		s := streamingLink{LinkUrl: steamUrl,
			PosterUrl: thumbNail,
			Title:     title}

		t, _ := template.ParseFiles("templates/player.html")
		err := t.Execute(w, s)
		if err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "The given URL is Invalid")
	}

}

func main() {
	welcomeScreen()
	http.HandleFunc("/player", handlePlayer)
	http.ListenAndServe(":8080", nil)
}
