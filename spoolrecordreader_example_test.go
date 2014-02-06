package unified2_test

import (
	"io"
	"log"
	"time"

	"github.com/jasonish/go-unified2"
)

func ExampleSpoolRecordReader() {

	// Create a SpoolRecordReader.
	reader := unified2.NewSpoolRecordReader("/var/log/snort", "unified2.log")

	for {
		record, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				// EOF is returned when the end of the last spool file
				// is reached and there is nothing else to read.  For
				// the purposes of the example, just sleep for a
				// moment and try again.
				time.Sleep(time.Millisecond)
			} else {
				// Unexpected error.
				log.Fatal(err)
			}
		}

		if record == nil {
			// The record and err are nil when there are no files at
			// all to be read.  This will happen if the Next() is
			// called before any files exist in the spool
			// directory. For now, sleep.
			time.Sleep(time.Millisecond)
			continue
		}

		switch record := record.(type) {
		case *unified2.EventRecord:
			log.Printf("Event: EventId=%d\n", record.EventId)
		case *unified2.ExtraDataRecord:
			log.Printf("- Extra Data: EventId=%d\n", record.EventId)
		case *unified2.PacketRecord:
			log.Printf("- Packet: EventId=%d\n", record.EventId)
		}

		filename, offset := reader.Offset()
		log.Printf("Current position: filename=%s; offset=%d", filename, offset)
	}

}
