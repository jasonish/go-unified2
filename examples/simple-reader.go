// Simplest record reader for a simple example.

package main

import "os"
import "log"
import "io"
import "github.com/jasonish/go-unified2"

func main() {

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	for {
		recordHolder, err := unified2.ReadRecord(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		switch record := recordHolder.Record.(type) {
		case *unified2.EventRecord:
			log.Printf("Event: EventId=%d\n", record.EventId)
		case *unified2.ExtraDataRecord:
			log.Printf("- Extra Data: EventId=%d\n", record.EventId)
		case *unified2.PacketRecord:
			log.Printf("- Packet: EventId=%d\n", record.EventId)
		}
	}

	file.Close()
}
