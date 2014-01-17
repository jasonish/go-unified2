package main

import (
	"fmt"
	unified2 "github.com/jasonish/go-unified2"
	"os"
)

func main() {
	file, _ := os.Open(os.Args[1])

	event := new(unified2.Event)
	fmt.Println(event.SensorId)

	queue := unified2.NewQueue()

	for {
		record, err := unified2.ReadRecord(file)
		if err != nil {
			fmt.Println(err)
			break
		}

		event := queue.Append(record)
		if event != nil {
			//fmt.Println(event)
		}
	}
}
