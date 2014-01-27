package unified2_test

import "os"
import "log"
import "io"
import "github.com/jasonish/go-unified2"

func ExampleReadRecord() {

	// Open a file.
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Read records.
	for {
		recordHolder, err := unified2.ReadRecord(file)
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

		switch record := recordHolder.Record.(type) {
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
