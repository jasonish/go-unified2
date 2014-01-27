package unified2

import "testing"
import "container/list"
import "os"
import "io"
import "log"

// Load unified2 records from a file, returning them as an array.
func LoadRecordsFromFile(filename string) ([]interface{}, error) {

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
	records := make([]interface{}, buffer.Len())
	record := buffer.Front()
	for key := range records {
		records[key] = record.Value.(interface{})
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
	record := *(records[0]).(*EventRecord)
	record.EventId++
	event := aggregator.Add(&record)
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
