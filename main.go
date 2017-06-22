// package main

// import "fmt"
// import "time"

// func inspectRec(r interface{}) {
// 	fmt.Printf("s : %v", r)
// }
// func subscriber(c <-chan struct{}) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			inspectRec(r)
// 		} else {
// 			fmt.Println("subscriber : Not Recover")
// 		}

// 		fmt.Println("subscriber exit")
// 	}()

// 	for range c {
// 		fmt.Printf(".")
// 	}

// 	fmt.Println("subscriber exit gracefully")
// }

// func publisher(c chan<- struct{}) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			inspectRec(r)
// 		} else {
// 			fmt.Println("publisher : Not Recover")
// 		}

// 		fmt.Println("publisher exit")
// 	}()

// 	tick := time.Tick(1 * time.Second)
// 	for range tick {
// 		c <- struct{}{}
// 	}
// }

// func main() {
// 	c := make(chan struct{})
// 	go subscriber(c)
// 	go publisher(c)

// 	time.Sleep(5 * time.Second)
// 	close(c)
// 	c <- struct{}{}
// }

package main

import (
	"fmt"
)

func main() {

	w := func() error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()

		// panic(fmt.Errorf("Hello"))
		return fmt.Errorf("Hello")
	}

	r := w()
	fmt.Printf("r %T", r)
}
