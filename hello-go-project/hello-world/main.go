package main

import (
	"fmt"
	"io"
)

// type InMemoryPlayerStore struct{}

// func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
// 	return 123
// }

//	func main() {
//		server := &PlayerServer{&InMemoryPlayerStore{}}
//		log.Fatal(http.ListenAndServe(":5000", server))
//	}
// func main() {
// 	Greet(os.Stdout, "Elodie")
// }

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

// func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
// 	Greet(w, "world")
// }

// func main() {
// 	// log.Fatal(http.ListenAndServe(":5001", http.HandlerFunc(MyGreeterHandler)))
// 	Countdown(os.Stdout)

// }

// func Countdown(out io.Writer) {
// 	fmt.Fprint(out, "3 /n")
// }

// func main() {
// 	sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
// 	Countdown(os.Stdout, sleeper)
// }
