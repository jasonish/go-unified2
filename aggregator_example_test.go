package unified2_test

import "os"
import "log"
import "io"
import "github.com/jasonish/go-unified2"

// EventAggregator example.
func ExampleEventAggregator() {

	// Create the aggregator.
	aggregator := unified2.NewEventAggregator()

	// Open a file.  Note that the aggregator is meant to span the
	// input of multiple files, as the records that make up a single
	// event may span multiple files.
	file, err := os.Open("merged.log")
	if err != nil {
		log.Fatal(err)
	}

	// Submit records to the aggregator, it will return non-nil when a
	// complete event has been seen.
	for {
		record, err := unified2.ReadRecord(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			// Unexpected error.
			log.Fatal(err)
		}

		event := aggregator.Add(record)
		if event != nil {
			log.Printf("We have an event consisting of %d records.\n",
				len(event))
		}
	}

	// Since we hit EOF we may not have triggered the last event to be
	// flushed, so check.
	event := aggregator.Flush()
	if event != nil {
		log.Printf("Final event flushed\n")
	} else {
		// Unlikely to happen.
		log.Printf("No remaining events.")
	}

}
