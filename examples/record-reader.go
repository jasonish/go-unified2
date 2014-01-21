// Example of reading individual unified2 records.

package main

import "os"
import "fmt"
import "io"
import "encoding/json"
import "github.com/jasonish/go-unified2"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("error: no filename specified")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("error opening ", os.Args[1], ":", err)
		os.Exit(1)
	}

	for {
		record, err := unified2.ReadRecord(file)
		if err != nil {
			if err != io.EOF {
				fmt.Println("failed to read record:", err)
			}
			break
		}

		// We now have a record.
		fmt.Printf("Record type: %d; length: %d\n", record.Type,
			len(record.Data))

		// Now that we have a record type, we can decode it, and print
		// out as JSON or whatever.
		if unified2.IsEventType(record) {

			// An event record.
			event, err := unified2.DecodeEvent(record)
			if err != nil {
				fmt.Println("error decoding event:", err)
				os.Exit(1)
			}

			// Print out as JSON.
			asJson, err := json.Marshal(*event)
			if err != nil {
				fmt.Println("failed to encode as json:", err)
				break
			}
			fmt.Println("  Decoded:", string(asJson))

		} else if record.Type == unified2.UNIFIED2_PACKET {

			// A packet record.
			packet, err := unified2.DecodePacket(record)
			if err != nil {
				fmt.Println("error decoding packet:", err)
				break
			}
			// Print out as JSON.
			asJson, err := json.Marshal(*packet)
			if err != nil {
				fmt.Println("failed to encode as json:", err)
				break
			}
			fmt.Println("  Decoded:", string(asJson))
		} else {
			// Add other types here.
		}
	}

}
