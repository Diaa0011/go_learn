package main

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"
// )

// // func TestRacer(t *testing.T) {
// // 	slowURL := "http://www.facebook.com"
// // 	fastURL := "http://www.quii.dev"

// // 	want := fastURL
// // 	got := Racer(slowURL, fastURL)

// // 	if got != want {
// // 		t.Errorf("got %q, want %q", got, want)
// // 	}
// // }

// // func TestRacerV2(t *testing.T) {
// // 	slowServer := makeDelayedServer(20 * time.Millisecond)

// // 	fastServer := makeDelayedServer(0 * time.Millisecond)

// // 	slowURL := slowServer.URL
// // 	fastURL := fastServer.URL

// // 	want := fastURL
// // 	got := Racer(slowURL, fastURL)

// // 	fmt.Printf("slowServer a: %q\n", slowURL)
// // 	fmt.Printf("fastServer a: %q\n", fastURL)

// // 	if got != want {
// // 		t.Errorf("got %q, want %q", got, want)
// // 	}

// // 	slowServer.Close()
// // 	fastServer.Close()
// // }

// func TestRacer3(t *testing.T) {
// 	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
// 		slowServer := makeDelayedServer(20 * time.Millisecond)
// 		fastServer := makeDelayedServer(0 * time.Millisecond)

// 		slowURL := slowServer.URL
// 		fastURL := fastServer.URL

// 		want := fastURL
// 		got, err := Racer(slowURL, fastURL)

// 		if err != nil {
// 			t.Fatalf("did not expect an error but goe one %v", err)
// 		}

// 		if got != want {
// 			t.Errorf("got %q, want %q", got, want)
// 		}
// 	})

// 	t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
// 		server := makeDelayedServer(11 * time.Second)

// 		defer server.Close()

// 		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

// 		if err == nil {
// 			t.Error("expected an error but didn't get one")
// 		}
// 	})
// }

// var tenSecondTimeout = 10 * time.Second

// func Racer(a, b string) (winner string, err error) {

// 	// aDuration := measureResponseTime(a)
// 	// bDuration := measureResponseTime(b)

// 	// fmt.Printf("duration a: %q\n", aDuration)
// 	// fmt.Printf("duration b: %q\n", bDuration)

// 	// if aDuration < bDuration {
// 	// 	return a
// 	// }

// 	// return b

// 	// select {
// 	// case <-ping(a):
// 	// 	return a, nil
// 	// case <-ping(b):
// 	// 	return b, nil
// 	// case <-time.After(10 * time.Second):
// 	// 	return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
// 	// }

// 	return ConfigurableRacer(a, b, tenSecondTimeout)
// }

// func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
// 	select {
// 	case <-ping(a):
// 		return a, nil
// 	case <-ping(b):
// 		return b, nil
// 	case <-time.After(timeout):
// 		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
// 	}
// }

// func ping(url string) chan struct{} {
// 	ch := make(chan struct{})
// 	go func() {
// 		http.Get(url)
// 		close(ch)
// 	}()

// 	return ch
// }

// func makeDelayedServer(delay time.Duration) *httptest.Server {
// 	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		time.Sleep(delay)
// 		w.WriteHeader(http.StatusOK)
// 	}))
// }

// func measureResponseTime(url string) time.Duration {
// 	start := time.Now()
// 	http.Get(url)
// 	return time.Since(start)
// }
