Problem Statement: Concurrent Sorting Using Goroutines and Channels

Write a Go program that takes multiple lists of integers and sorts each list concurrently.

You need to:

Create multiple integer slices.
Start one goroutine for each slice.
Each goroutine should sort its assigned slice.
After sorting, each goroutine should send the sorted slice back through a channel.
Use a sync.WaitGroup to wait for all sorting goroutines to finish.
Close the channel after all goroutines are done.
Print all sorted slices received from the channel.

Example input inside code:

a := []int{5, 2, 9, 1}
b := []int{8, 3, 7, 4}
c := []int{10, 6, 12, 11}