// Example of reading individual unified2 records.
//
// For benchmarking record's per second use -quiet. To benchmark just
// the record reading without decoding also use -nodecode.

package main

import "os"
import "fmt"
import "io"
import "flag"
import "time"
import "github.com/jasonish/go-unified2"

type Stats struct {
	Events    int
	Packets   int
	ExtraData int
}

func main() {

	flag.Parse()
	args := flag.Args()

	startTime := time.Now()
	var recordCount int
	var stats Stats
	for _, arg := range args {

		fmt.Println("Opening", arg)
		file, err := os.Open(arg)
		if err != nil {
			fmt.Println("error opening ", arg, ":", err)
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
			recordCount++

			if unified2.IsEventType(record.Type) {

				stats.Events++

			} else if record.Type == unified2.UNIFIED2_PACKET {

				stats.Packets++

			} else if record.Type == unified2.UNIFIED2_EXTRA_DATA {

				stats.ExtraData++

			}
		}

		file.Close()
	}

	elapsedTime := time.Now().Sub(startTime)
	perSecond := float64(recordCount) / elapsedTime.Seconds()

	fmt.Printf("Records: %d; Time: %s; Records/sec: %d\n",
		recordCount, elapsedTime, int(perSecond))
	fmt.Printf("  Events: %d; Packets: %d; ExtraData: %d\n",
		stats.Events, stats.Packets, stats.ExtraData)
}
