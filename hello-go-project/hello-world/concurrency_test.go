package main

// import (
// 	"reflect"
// 	"testing"
// )

// func mockWebsiteChecker(url string) bool {
// 	return url != "waat://furhurterwe.geds"
// }

// func TestCheckWebsites(t *testing.T) {
// 	websites := []string{
// 		"http://google.com",
// 		"http://blog.gypsydave5.com",
// 		"waat://furhurterwe.geds",
// 	}

// 	want := map[string]bool{
// 		"http://google.com":          true,
// 		"http://blog.gypsydave5.com": true,
// 		"waat://furhurterwe.geds":    false,
// 	}

// 	got := CheckWebsites(mockWebsiteChecker, websites)

// 	if !reflect.DeepEqual(want, got) {
// 		t.Fatalf("wanted %v, got %v", want, got)
// 	}
// }

// type WebsiteChecker func(string) bool

// type result struct {
// 	string
// 	bool
// }

// func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
// 	results := make(map[string]bool)
// 	resultChannel := make(chan result)

// 	for _, url := range urls {
// 		go func() {
// 			resultChannel <- result{url, wc(url)}
// 		}()
// 	}

// 	for i := 0; i < len(urls); i++ {
// 		r := <-resultChannel
// 		results[r.string] = r.bool
// 	}

// 	return results
// }
