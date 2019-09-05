package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type TestUrls struct {
	url string
	err bool
}

var case1 = `{"url":[""]}`
var ans1 = `[{"url":"","success":false}]`
var case2 = `{"url":["https://blog.golang.org/go-image-package_image-package-01.png"]}`
var ans2 = `[{"url":"https://blog.golang.org/go-image-package_image-package-01.png","success":true}]`
var case3 = `{"url":["https://blog.golang.org/go-image-package_image-package-01.png","1","https://miro.medium.com/max/870/1*b3XJkfO_e6b251CWnZ8g7A.png"]}`
var ans3 = `[{"url":"https://blog.golang.org/go-image-package_image-package-01.png","success":true},{"url":"1","success":false},{"url":"https://miro.medium.com/max/870/1*b3XJkfO_e6b251CWnZ8g7A.png","success":true}]`
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
	test1 := []string{case1,ans1}
	test2 := []string{case2,ans2}
	test3 := []string{case3,ans3}
	tests := [][]string{test1,test2,test3}
	for _, testCase := range tests {
		req, err := http.NewRequest("POST", "/images", bytes.NewBuffer([]byte(testCase[0])))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(postAddMetaPicture)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned unexpected body: got\n%v\nwant\n%v\n",
				status, http.StatusOK)
		}
		if strings.TrimSuffix(rr.Body.String() , "\n") != testCase[1] {
			t.Errorf("handler returned unexpected body: got\n%v\nwant\n%v\n",
				strings.TrimSuffix(rr.Body.String() , "\n") , testCase[1])
		}
	}
}
// Test with multiply simultaneously goroutines for race conditions
func TestPostAddMetaPictureRaceCondition(t *testing.T) {
	test1 := []string{case1,ans1}
	test2 := []string{case2,ans2}
	test3 := []string{case3,ans3}
	tests := [][]string{test1,test2,test3}
	var wg sync.WaitGroup
	start := make(chan 	struct {})
	wg.Add(len(tests))
	for _, testCase := range tests {
		go func(testCase []string) {
			defer wg.Done()
			req, err := http.NewRequest("POST", "/images", bytes.NewBuffer([]byte(testCase[0])))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(postAddMetaPicture)
			<-start
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned unexpected body: got\n%v\nwant\n%v\n",
					status, http.StatusOK)
			}
			if strings.TrimSuffix(rr.Body.String() , "\n")  != testCase[1] {
				t.Errorf("handler returned unexpected body: got\n%v\nwant\n%v\n",
					strings.TrimSuffix(rr.Body.String() , "\n") , testCase[1])
			}
		}(testCase)
	}
	close(start)
	wg.Wait()
}
