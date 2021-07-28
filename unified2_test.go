package unified2

import (
	"errors"
	"os"
	"testing"
)

// Check that we get EOF at the end of a file.
func TestReadRecordEOF(t *testing.T) {

	// Use test/multi-record-event.log, its a complete file and should
	// finish up with an EOF.
	input, err := os.Open("test/multi-record-event.log")
	if err != nil {
		t.Fatal(err)
	}

	for {
		_, err := ReadRecord(input)
		if err != nil {
			if e := (&ErrBufferTooSmall{}); !errors.As(err, &e) {
				t.Fatalf("expected err of BufferTooSmall, got %s", err)
			}
			break
		}
	}

}

func TestShortReadOnHeader(t *testing.T) {

	input, err := os.Open("test/short-read-on-header.log")
	if err != nil {
		t.Fatal(err)
	}

	_, err = ReadRecord(input)
	if err == nil {
		t.Fatalf("expected non-nil err")
	}

	if e := (&ErrBufferTooSmall{}); !errors.As(err, &e) || e.MissingBytes != 2 {
		t.Fatalf("expected err == io.BufferTooSmall with MissingBytes = 2, got %s", err)
	}
	offset, err := input.Seek(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if offset != 0 {
		t.Fatalf("expected file offset to be at 0, was at %d", offset)
	}

	input.Close()
}

func TestInvalidHeader(t *testing.T) {

	input, err := os.Open("test/invalid-header.log")
	if err != nil {
		t.Fatal(err)
	}

	_, err = ReadRecord(input)
	if err == nil {
		t.Fatalf("expected non-nil err")
	}

	if !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("expected err == BufferTooSmallError with MissingBytes == 56, got %s", err)
	}
	offset, err := input.Seek(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if offset != 0 {
		t.Fatalf("expected file offset to be at 0, was at %d", offset)
	}

	input.Close()

}

func TestShortReadOnBody(t *testing.T) {

	input, err := os.Open("test/short-read-on-body.log")
	if err != nil {
		t.Fatal(err)
	}

	_, err = ReadRecord(input)
	if err == nil {
		t.Fatalf("expected non-nil err")
	}
	if e := (&ErrBufferTooSmall{}); !errors.As(err, &e) || e.MissingBytes != 56 {
		t.Fatalf("expected err == BufferTooSmallError with MissingBytes == 56, got %s", err)
	}
	offset, err := input.Seek(0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if offset != 0 {
		t.Fatalf("expected file offset to be at 0, was at %d", offset)
	}

	input.Close()

}

func TestDecodeError(t *testing.T) {

	data := []byte("this should fail")

	_, err := DecodeEventRecord(UNIFIED2_EVENT_V2, data)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestRecordReaderWithOffset(t *testing.T) {
	test_filename := "test/multi-record-event.log"

	// First open a known file at offset 0.
	reader, err := NewRecordReader(test_filename, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Read one record.
	record, err := reader.Next()
	if err != nil {
		t.Fatal(err)
	}
	if record == nil {
		t.Fatalf("unexpected nil record")
	}

	offset := reader.Offset()
	if offset == 0 {
		t.Fatalf("unexpected offset %d", offset)
	}

	// Close and reopen with offset, check offset and make sure the
	// first record returned is not an event record.
	reader.Close()
	reader, err = NewRecordReader(test_filename, offset)
	if err != nil {
		t.Fatal(err)
	}
	if offset != reader.Offset() {
		t.Fatalf("unexpected reader offset: expected %d; got %d", offset,
			reader.Offset())
	}
	record, err = reader.Next()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := record.(*EventRecord); ok {
		t.Fatal("did not expect Next() to return *EventRecord")
	}

}
