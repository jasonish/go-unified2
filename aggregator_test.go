package unified2

import "testing"
import "container/list"
import "os"
import "io"
import "log"

// Load unified2 records from a file, returning them as an array.
func LoadRecordsFromFile(filename string) ([]*RecordContainer, error) {

	input, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	buffer := list.New()

	for {
		record, err := ReadRecord(input)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		buffer.PushBack(record)
	}

	input.Close()

	// Return as an array.
	records := make([]*RecordContainer, buffer.Len())
	record := buffer.Front()
	for key, _ := range records {
		records[key] = record.Value.(*RecordContainer)
		record = record.Next()
	}

	return records, nil
}

// Add an event, flush it.
func TestSimpleAddFlush(t *testing.T) {

	records, err := LoadRecordsFromFile("test/multi-record-event.log")
	if err != nil {
		t.Fail()
	}
	log.Printf("Loaded %d records.\n", len(records))
	if len(records) != 17 {
		t.Fatalf("Loaded %d records, expected 17.\n", len(records))
	}

	aggregator := NewEventAggregator()

	event := aggregator.Add(records[0])
	if event != nil {
		t.Fatalf("expected no event to be returned\n")
	}
	if aggregator.Len() != 1 {
		t.Fatal("expected aggregator length to be 1")
	}

	event = aggregator.Flush()
	if len(event) != 1 {
		t.Fatal("Expected only one record in event.")
	}
	if aggregator.Len() != 0 {
		t.Fatal("Expected aggregator to be empty.")
	}
}

func TestCompleteEvent(t *testing.T) {

	// Load the test records.
	records, err := LoadRecordsFromFile("test/multi-record-event.log")
	if err != nil {
		t.Fail()
	}
	log.Printf("Loaded %d records.\n", len(records))
	if len(records) != 17 {
		t.Fatalf("Loaded %d records, expected 17.\n", len(records))
	}

	aggregator := NewEventAggregator()

	// Load all the records, as they make up one event, an event
	// should not be generated.
	for _, record := range records {
		event := aggregator.Add(record)
		if event != nil {
			t.Fatalf("did not expect an event to be returned")
		}
	}

	// To signal the aggregator to flush, we'll mock up an event with
	// a new event-id.
	record := *(records[0]).Record.(*EventRecord)
	record.EventId++
	event := aggregator.Add(&RecordContainer{UNIFIED2_IDS_EVENT_V2, &record})
	if event == nil {
		t.Fatalf("expected event records to be returned")
	}
	if len(event) != 17 {
		t.Fatalf("expected 17 records in event")
	}
	if aggregator.Len() != 1 {
		t.Fatalf("expected 1 record in aggregator")
	}
}

// EventAggregator example.
func ExampleEventAggregator() {

	// Create the aggregator.
	aggregator := NewEventAggregator()

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
		recordHolder, err := ReadRecord(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			// Unexpected error.
			log.Fatal(err)
		}

		event := aggregator.Add(recordHolder)
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
