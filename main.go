package main

import "fmt"
import "strings"

type item struct{ id string }
type collection struct {
	items []item
}

func main() {
	s := "{asuh}.koe.jaran"
	if strings.HasPrefix(s, "{asuh}.") {
		fmt.Println("due")
	}

	fmt.Println(strings.TrimPrefix(s, "{asuh}."))
}

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func main() {
// 	var wg sync.WaitGroup
// 	for i := 0; i < 10; i++ {
// 		ii := i * 3
// 		wg.Add(1)
// 		go func(iii *int) {
// 			time.Sleep(time.Second * i)
// 			fmt.Printf("ii : %d\n", *iii)
// 			wg.Done()
// 		}(&ii)
// 	}
// 	wg.Wait()
// }
