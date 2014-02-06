package unified2

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

// Utility function to copy a file.
func copyFile(source string, dest string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

// Test that only files with the provided prefix are returned.
func TestRecordSpoolReader_getFiles(t *testing.T) {

	test_filename := "test/multi-record-event.log"

	tmpdir, err := ioutil.TempDir("", "unified2-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	copyFile(test_filename, fmt.Sprintf("%s/merged.log.005", tmpdir))
	copyFile(test_filename, fmt.Sprintf("%s/merged.log.004", tmpdir))
	copyFile(test_filename, fmt.Sprintf("%s/merged.log.003", tmpdir))
	copyFile(test_filename, fmt.Sprintf("%s/asdf.log.002", tmpdir))

	reader := NewSpoolRecordReader(tmpdir, "merged.log")
	if reader == nil {
		t.Fatal("reader should not be nil")
	}

	files, err := reader.getFiles()
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "merged.log") {
			t.Fatalf("unexpected filename: %s", file.Name())
		}
	}
}

// Basic test for RecordSpoolReader.
func TestRecordSpoolReader(t *testing.T) {

	test_filename := "test/multi-record-event.log"

	tmpdir, err := ioutil.TempDir("", "unified2-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	closeHookCount := 0

	reader := NewSpoolRecordReader(tmpdir, "merged.log")
	if reader == nil {
		t.Fatal("reader should not be nil")
	}
	reader.Logger(log.New(os.Stderr, "RecordSpoolReader: ", 0))
	reader.CloseHook = func(filename string) {
		closeHookCount++
	}

	// Offset should return an empty string and 0.
	filename, offset := reader.Offset()
	if filename != "" {
		t.Fatal(filename)
	}
	if offset != 0 {
		t.Fatal(offset)
	}

	copyFile(test_filename, fmt.Sprintf("%s/merged.log.1382627900", tmpdir))

	files, err := reader.getFiles()
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		log.Println(file.Name())
	}

	// Read the first record and check the offset.
	record, err := reader.Next()
	if err != nil {
		t.Fatal(err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}
	filename, offset = reader.Offset()
	if filename != "merged.log.1382627900" {
		t.Fatalf("got %s, expected %s", filename, "merged.log.1382627900")
	}

	// Offset known from previous testing.
	if offset != 68 {
		t.Fatal("bad offset")
	}

	// We know the input file has 17 records, so read 16 and make sure
	// we get back a record for each call.
	for i := 0; i < 16; i++ {
		record, err := reader.Next()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if record == nil {
			t.Fatalf("unexpected nil record")
		}
	}

	// On the next call, record should be nul and we should have an
	// error of EOF.
	record, err = reader.Next()
	if record != nil || err != io.EOF {
		t.Fatalf("unexpected results: record not nil, err not EOF")
	}

	// Copy in another file that should be picked up by the spool
	// reader.
	copyFile(test_filename, fmt.Sprintf("%s/merged.log.1382627901", tmpdir))

	// We should now read records again.
	record, err = reader.Next()
	if record == nil {
		t.Fatalf("expected non-nil record: err=%s", err)
	}

	if closeHookCount != 1 {
		t.Fatalf("bad closeHookCount: expected 1, got %d", closeHookCount)
	}
}
