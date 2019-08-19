package main

import (
	"encoding/json"
	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"sync"
)

type StatePicture struct{
	Url			string		`json:"url"`
	State 		bool		`json:"success"`
}
type MetaPicture struct{
	Url			[]string	`json:"url"`
}
type Picture struct {
	Url			string 		`json:"url"`// picture metadata
	Picture 	image.Image	`json:"-"`//actual picture
	Thumbnail 	image.Image	`json:"-"`
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

func processMetaPicture(metaPicture MetaPicture) ([]Picture,[]StatePicture) {
	var states []StatePicture
	var pictures []Picture
	for _, url := range metaPicture.Url {
		body, err := downloadMetaPicture(url)
		if err != nil {
			states = append(states, StatePicture{url,false})
			continue
		}
		thumbnail := createTumbnail(body)
		pictures = append(pictures, Picture{url,body,thumbnail})
		err = imaging.Save(thumbnail, "C:\\Users\\Миха С\\Desktop\\ctf\\go\\thumbnail.jpg")
		if err != nil {
			log.Print(err.Error())
		}
		err = imaging.Save(body, "C:\\Users\\Миха С\\Desktop\\ctf\\go\\body.jpg")
		if err != nil {
			log.Print(err.Error())
		}
		states = append(states, StatePicture{url,true})
	}
	return pictures,states
}
func createTumbnail(source image.Image) image.Image {
	return imaging.Thumbnail(source, 100, 100, imaging.CatmullRom)
}

func downloadMetaPicture(url string) (image.Image, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil,err
	}
	data, errstr, err := image.Decode(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Print(errstr)
		return nil,err
	}
	return data,nil
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/add", postAddMetaPicture).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
