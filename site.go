package main

import (
	"context"
	"encoding/json"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type StatePicture struct {
	Url   string `json:"url"`
	State bool   `json:"success"`
}
type MetaPicture struct {
	Url []string `json:"url"`
}
type Picture struct {
	Url       string      `json:"url"` // picture metadata
	Picture   image.Image `json:"-"`   //actual picture
	Thumbnail image.Image `json:"-"`
}

var mutex sync.Mutex
var Pictures []Picture


func postAddMetaPicture(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var metaPicture MetaPicture
	_ = json.NewDecoder(r.Body).Decode(&metaPicture)
	pictures, states := processMetaPicture(metaPicture)
	mutex.Lock()
	Pictures = append(Pictures, pictures...)
	mutex.Unlock()
	json.NewEncoder(w).Encode(states)
}

func processMetaPicture(metaPicture MetaPicture) ([]Picture, []StatePicture) {
	var states []StatePicture
	var pictures []Picture
	for _, url := range metaPicture.Url {
		body, err := downloadMetaPicture(url)
		if err != nil {
			states = append(states, StatePicture{url, false})
			continue
		}
		thumbnail := createThumbnail(body)
		pictures = append(pictures, Picture{url, body, thumbnail})
		states = append(states, StatePicture{url, true})
	}
	return pictures, states
}
func createThumbnail(source image.Image) image.Image {
	return imaging.Thumbnail(source, 100, 100, imaging.CatmullRom)
}

func downloadMetaPicture(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, errstr, err := image.Decode(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Print(errstr)
		return nil, err
	}
	return data, nil
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/images", postAddMetaPicture).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	gravefulShutdown(srv)
}

func gravefulShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
