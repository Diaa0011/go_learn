package main

import "fmt"

//	func Hello() string {
//		return "Hello, world"
//	}
const ENGLISHPREFIX = "Hello, "

func HelloName(name string) string {
	if name == "" {
		name = "World"
	}

	return ENGLISHPREFIX + name
}

func main() {
	// fmt.Println(Hello())

	fmt.Println(HelloName("Ahmed"))
}
