package unified2_test

import "fmt"
import "io"
import "log"
import "github.com/jasonish/go-unified2"

// EventAggregator example.
func ExampleEventAggregator() {

	// Create the aggregator.
	aggregator := unified2.NewEventAggregator()

	// Use a RecordReader to read records from a file.
	reader, err := unified2.NewRecordReader("test/multi-record-event-x2.log", 0)
	if err != nil {
		log.Fatal(err)
	}

	// Submit records to the aggregator, it will return non-nil when a
	// complete event has been seen.
	for {
		record, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			// Unexpected error.
			log.Fatal(err)
		}

		// Add the record to the aggregator.  If the records signals
		// that start of a new event, the previous event will be
		// returned as an array of records.
		event := aggregator.Add(record)
		if event != nil {
			fmt.Printf("Got event consisting of %d records.\n",
				len(event))
		}
	}

	// Since we hit EOF we may not have triggered the last event to be
	// flushed, so check.
	event := aggregator.Flush()
	if event != nil {
		fmt.Printf("Flushed pending event of %d records.\n", len(event))
	} else {
		// Unlikely to happen.
		log.Printf("No remaining events.\n")
	}

	// Output:
	// Got event consisting of 17 records.
	// Flushed pending event of 17 records.
}
