package main

import (
	"fmt"
	"time"
)

func getCurrentTime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

func main() {
	fmt.Println("Hello")
	currentTime := getCurrentTime()
	fmt.Printf("Current time: %s\n", currentTime)
}
