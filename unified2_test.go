package unified2

import "testing"
import "os"
import "io"

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
			if err != io.EOF {
				t.Fatalf("expected err of io.EOF, got %s", err)
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
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected err == io.ErrUnexpectedEOF, got %s", err)
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
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("expected err == io.ErrUnexpectedEOF, got %s", err)
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
