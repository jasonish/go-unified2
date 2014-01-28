package unified2_test

import (
	"github.com/jasonish/go-unified2"
	"io"
	"log"
	"os"
)

func ExampleReadRecord() {

	// Open a file.
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Read records.
	for {
		record, err := unified2.ReadRecord(file)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				// End of file is reached.  You may want to break here
				// or sleep and try again if you are expected more
				// data to be written to the input file.
				//
				// Lets break for the purpose of this example.
				break
			} else if err == unified2.DecodingError {
				// Error decoding a record, probably corrupt.
				log.Fatal(err)
			}
			// Some other error.
			log.Fatal(err)
		}

		switch record := record.(type) {
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

// RecordReader example.
func ExampleRecordReader() {

	// Create a reader starting at offset 0 of the provided file.
	reader, err := unified2.NewRecordReader("test/multi-record-event.log", 0)
	if err != nil {
		log.Fatal(err)
	}

	for {
		record, err := reader.Next()
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				// End of file is reached.  You may want to break here
				// or sleep and try again if you are expected more
				// data to be written to the input file.
				//
				// Lets break for the purpose of this example.
				break
			} else if err == unified2.DecodingError {
				// Error decoding a record, probably corrupt.
				log.Fatal(err)
			}
			// Some other error.
			log.Fatal(err)
		}

		switch record := record.(type) {
		case *unified2.EventRecord:
			log.Printf("Event: EventId=%d\n", record.EventId)
		case *unified2.ExtraDataRecord:
			log.Printf("- Extra Data: EventId=%d\n", record.EventId)
		case *unified2.PacketRecord:
			log.Printf("- Packet: EventId=%d\n", record.EventId)
		}
	}
}
