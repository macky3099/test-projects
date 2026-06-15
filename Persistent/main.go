package main

import (
	"fmt"
	"sync"
)

func main() {
	// list of integer slices
	a := []int{5, 2, 9, 1}
	b := []int{8, 3, 7, 4}
	c := []int{10, 6, 12, 11}

	lists := [][]int{a, b, c}

	var wg sync.WaitGroup

	ch := make(chan []int)

	for _, list := range lists {
		wg.Add(1)

		go sortList(list, ch, &wg)
	}

	// close channel once all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// receive sorted lists from channel
	for sortedList := range ch {
		fmt.Println(sortedList)
	}
}

func sortList(list []int, ch chan []int, wg *sync.WaitGroup) {
	defer wg.Done()

	// simple bubble sort
	for i := 0; i < len(list); i++ {
		for j := 0; j < len(list)-i-1; j++ {
			if list[j] > list[j+1] {
				list[j], list[j+1] = list[j+1], list[j]
			}
		}
	}

	// send sorted list back
	ch <- list
}
