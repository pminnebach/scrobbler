package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Upload route
	http.HandleFunc("/plex", uploadHandler)

	//Listen on port 8080
	http.ListenAndServe(":8090", nil)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {

	log.Println("request")
	for name, headers := range r.Header {
		for _, h := range headers {
			log.Printf("KEY: %v VALUE: %v\n", name, h)
		}
	}

	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	p := r.FormValue("payload")
	log.Printf("PAYLOAD: %v\n", p)

	var wh Webhook
	json.Unmarshal([]byte(p), &wh)
	log.Printf("EVENT: %v\n", wh.Event)
	log.Printf("gpTitle: %v\n", wh.Metadata.GrandparentTitle)
	log.Printf("pTitle: %v\n", wh.Metadata.ParentTitle)
	log.Printf("Title: %v\n", wh.Metadata.Title)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("thumb")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create(handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		uploadFile(w, r)
	}
}

// Webhook contains the complete payload that plex sends out.
// https://support.plex.tv/articles/115002267687-webhooks/
type Webhook struct {
	Event   string `json:"event"`
	User    bool   `json:"user"`
	Owner   bool   `json:"owner"`
	Account struct {
		ID    int    `json:"id"`
		Thumb string `json:"thumb"`
		Title string `json:"title"`
	} `json:"Account"`
	Server struct {
		Title string `json:"title"`
		UUID  string `json:"uuid"`
	} `json:"Server"`
	Player struct {
		Local         bool   `json:"local"`
		PublicAddress string `json:"publicAddress"`
		Title         string `json:"title"`
		UUID          string `json:"uuid"`
	} `json:"Player"`
	Metadata struct {
		LibrarySectionType    string  `json:"librarySectionType"`
		RatingKey             string  `json:"ratingKey"`
		Key                   string  `json:"key"`
		ParentRatingKey       string  `json:"parentRatingKey"`
		GrandparentRatingKey  string  `json:"grandparentRatingKey"`
		GUID                  string  `json:"guid"`
		ParentGUID            string  `json:"parentGuid"`
		GrandparentGUID       string  `json:"grandparentGuid"`
		Type                  string  `json:"type"`
		Title                 string  `json:"title"`
		GrandparentTitle      string  `json:"grandparentTitle"`
		ParentTitle           string  `json:"parentTitle"`
		ContentRating         string  `json:"contentRating"`
		Summary               string  `json:"summary"`
		Index                 int     `json:"index"`
		ParentIndex           int     `json:"parentIndex"`
		Rating                float64 `json:"rating"`
		Year                  int     `json:"year"`
		Thumb                 string  `json:"thumb"`
		Art                   string  `json:"art"`
		ParentThumb           string  `json:"parentThumb"`
		GrandparentThumb      string  `json:"grandparentThumb"`
		GrandparentArt        string  `json:"grandparentArt"`
		GrandparentTheme      string  `json:"grandparentTheme"`
		OriginallyAvailableAt string  `json:"originallyAvailableAt"`
		AddedAt               int     `json:"addedAt"`
		UpdatedAt             int     `json:"updatedAt"`
		ChapterSource         string  `json:"chapterSource"`
		Director              []struct {
			ID  int    `json:"id"`
			Tag string `json:"tag"`
		} `json:"Director"`
		Writer []struct {
			ID  int    `json:"id"`
			Tag string `json:"tag"`
		} `json:"Writer"`
	} `json:"Metadata"`
}
