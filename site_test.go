package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestPostAdd struct {
	Data     []byte
	Expected []byte
}
type TestUrls struct {
	url string
	err bool
}

var postAddQueryes = []TestPostAdd{
	{
		[]byte(`{"url":[""]}`),
		[]byte(`{"url":[""]}`),
	},
	{
		[]byte(`{
	"url": [
			"https://blog.golang.org/go-image-package_image-package-01.png"]
	}`),
		[]byte(`[
	{
		"url":"https://blog.golang.org/go-image-package_image-package-01.png",
		"success": true
	}]`),
	},
	{
		[]byte(`{
    "url": [
        "https://blog.golang.org/go-image-package_image-package-01.png",
        "1",
        "https://miro.medium.com/max/870/1*b3XJkfO_e6b251CWnZ8g7A.png"
    ]
}`),
		[]byte(`[
    {
        "url": "https://blog.golang.org/go-image-package_image-package-01.png",
        "success": true
    },
    {
        "url": "1",
        "success": false
    },
    {
        "url": "https://miro.medium.com/max/870/1*b3XJkfO_e6b251CWnZ8g7A.png",
        "success": true
    }
]`),
	},
}
var urls = []TestUrls{
	{"123", true},
	{"https://blog.golang.org/go-image-package_image-package-01.png", false},
	{"https://blog.golang.org/", true},
	{"https://camo.githubusercontent.com/de5c9030adaa419ed202a1b2ed0939d228e46804/687474703a2f2f7777772e676f72696c6c61746f6f6c6b69742e6f72672f7374617469632f696d616765732f676f72696c6c612d69636f6e2d36342e706e67", false},
}

func TestDownloadMetaPicture(t *testing.T) {
	for _, url := range urls {
		_, err := downloadMetaPicture(url.url)
		if (err != nil) != url.err {
			t.Error("For link ", url.url, "Error expected ", url.err, "got", (err != nil))
		}
	}
}

func TestPostAddMetaPicture(t *testing.T) {
	for _, testCase := range postAddQueryes {
		req, err := http.NewRequest("POST", "/images", bytes.NewBuffer(testCase.Data))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(postAddMetaPicture)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		var data string
		var expected string
		json.NewDecoder(rr.Body).Decode(&data)
		json.NewDecoder(bytes.NewBuffer(testCase.Expected)).Decode(&expected)
		if data != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				data, expected)
		}
	}
}
