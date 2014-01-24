/* Copyright (c) 2013 Jason Ish
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

/*

Package unified2 provides a decoder for unified v2 log files
produced by Snort and Suricata.

Example usage:

    func main() {

    	file, err := os.Open(os.Args[1])
    	if err != nil {
    		log.Fatal(err)
    	}

    	for {
    		recordHolder, err := unified2.ReadRecord(file)
    		if err != nil {
    			if err == io.EOF {
    				break
    			}
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

*/
package unified2

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

// Unified2 record types.
const (
	UNIFIED2_PACKET           = 2
	UNIFIED2_IDS_EVENT        = 7
	UNIFIED2_IDS_EVENT_IP6    = 72
	UNIFIED2_IDS_EVENT_V2     = 104
	UNIFIED2_IDS_EVENT_IP6_V2 = 105
	UNIFIED2_EXTRA_DATA       = 110
)

// RawHeader is the raw unified2 record header.
type RawHeader struct {
	Type uint32
	Len  uint32
}

// RawRecord is a holder type for a raw un-decoded record.
type RawRecord struct {
	Type uint32
	Data []byte
}

// EventRecord is a struct representing a decoded event record.
//
// This struct is used to represent the decoded form of all the event
// types.  The difference between an IPv4 and IPv6 event will be the
// length of the IP address IpSource and IpDestination.
type EventRecord struct {
	SensorId          uint32
	EventId           uint32
	EventSecond       uint32
	EventMicrosecond  uint32
	SignatureId       uint32
	GeneratorId       uint32
	SignatureRevision uint32
	ClassificationId  uint32
	Priority          uint32
	IpSource          []byte
	IpDestination     []byte
	SportItype        uint16
	DportIcode        uint16
	Protocol          uint8
	ImpactFlag        uint8
	Impact            uint8
	Blocked           uint8
	MplsLabel         uint32
	VlanId            uint16
}

// PacketRecord is a struct representing a decoded packet record.
type PacketRecord struct {
	SensorId          uint32
	EventId           uint32
	EventSecond       uint32
	PacketSecond      uint32
	PacketMicrosecond uint32
	LinkType          uint32
	Length            uint32
	Data              []byte
}

// The length of a PacketRecord before variable length data.
const PACKET_RECORD_HDR_LEN = 28

// ExtraDataRecord is a struct representing a decoded extra data record.
type ExtraDataRecord struct {
	EventType   uint32
	EventLength uint32
	SensorId    uint32
	EventId     uint32
	EventSecond uint32
	Type        uint32
	DataType    uint32
	DataLength  uint32
	Data        []byte
}

// The length of an ExtraDataRecord before variable length data.
const EXTRA_DATA_RECORD_HDR_LEN = 32

// RecordContainer is a container struct for decoded records.
type RecordContainer struct {

	//The record type.
	Type uint32

	// The decoded record. One of EventRecord, PacketRecord or
	// ExtraDataRecord.
	Record interface{}
}

// Helper function for reading binary data as all reads are big
// endian.
func read(reader io.Reader, data interface{}) error {
	return binary.Read(reader, binary.BigEndian, data)
}

// IsEventType checks if a record type is an event type or not.  It
// returns true if the recordType is an event type, otherwise false
// will be returned.
func IsEventType(recordType uint32) bool {
	switch recordType {
	case UNIFIED2_IDS_EVENT,
		UNIFIED2_IDS_EVENT_IP6,
		UNIFIED2_IDS_EVENT_V2,
		UNIFIED2_IDS_EVENT_IP6_V2:
		return true
	default:
		return false
	}
}

// DecodeEventRecord decodes a raw record into an EventRecord.
//
// This function will decode any of the event record types.
func DecodeEventRecord(
	eventType uint32, data []byte) (event *EventRecord, err error) {

	event = &EventRecord{}

	reader := bytes.NewBuffer(data)

	// SensorId
	if err = read(reader, &event.SensorId); err != nil {
		return nil, err
	}
	if err = read(reader, &event.EventId); err != nil {
		return nil, err
	}
	if err = read(reader, &event.EventSecond); err != nil {
		return nil, err
	}
	if err = read(reader, &event.EventMicrosecond); err != nil {
		return nil, err
	}

	/* SignatureId */
	if err = read(reader, &event.SignatureId); err != nil {
		return nil, err
	}

	/* GeneratorId */
	if err = read(reader, &event.GeneratorId); err != nil {
		return nil, err
	}

	/* SignatureRevision */
	if err = read(reader, &event.SignatureRevision); err != nil {
		return nil, err
	}

	/* ClassificationId */
	if err = read(reader, &event.ClassificationId); err != nil {
		return nil, err
	}

	/* Priority */
	if err = read(reader, &event.Priority); err != nil {
		return nil, err
	}

	/* Source and destination IP addresses. */
	switch eventType {

	case UNIFIED2_IDS_EVENT, UNIFIED2_IDS_EVENT_V2:
		event.IpSource = make([]byte, 4)
		if err = read(reader, &event.IpSource); err != nil {
			return nil, err
		}
		event.IpDestination = make([]byte, 4)
		if err = read(reader, &event.IpDestination); err != nil {
			return nil, err
		}

	case UNIFIED2_IDS_EVENT_IP6, UNIFIED2_IDS_EVENT_IP6_V2:
		event.IpSource = make([]byte, 16)
		if err = read(reader, &event.IpSource); err != nil {
			return nil, err
		}
		event.IpDestination = make([]byte, 16)
		if err = read(reader, &event.IpDestination); err != nil {
			return nil, err
		}
	}

	/* Source port/ICMP type. */
	if err = read(reader, &event.SportItype); err != nil {
		return nil, err
	}

	/* Destination port/ICMP code. */
	if err = read(reader, &event.DportIcode); err != nil {
		return nil, err
	}

	/* Protocol. */
	if err = read(reader, &event.Protocol); err != nil {
		return nil, err
	}

	/* Impact flag. */
	if err = read(reader, &event.ImpactFlag); err != nil {
		return nil, err
	}

	/* Impact. */
	if err = read(reader, &event.Impact); err != nil {
		return nil, err
	}

	/* Blocked. */
	if err = read(reader, &event.Blocked); err != nil {
		return nil, err
	}

	switch eventType {
	case UNIFIED2_IDS_EVENT_V2, UNIFIED2_IDS_EVENT_IP6_V2:

		/* MplsLabel. */
		if err = read(reader, &event.MplsLabel); err != nil {
			return nil, err
		}

		/* VlanId. */
		if err = read(reader, &event.VlanId); err != nil {
			return nil, err
		}

	}

	return event, nil
}

// DecodePacketRecord decodes a raw unified2 record into a
// PacketRecord.
func DecodePacketRecord(data []byte) (packet *PacketRecord, err error) {

	packet = &PacketRecord{}

	reader := bytes.NewBuffer(data)

	err = read(reader, &packet.SensorId)
	if err != nil {
		return nil, err
	}

	err = read(reader, &packet.EventId)
	if err != nil {
		return nil, err
	}

	err = read(reader, &packet.EventSecond)
	if err != nil {
		return nil, err
	}

	err = read(reader, &packet.PacketSecond)
	if err != nil {
		return nil, err
	}

	err = read(reader, &packet.PacketMicrosecond)
	if err != nil {
		return nil, err
	}

	err = read(reader, &packet.LinkType)
	if err != nil {
		return nil, err
	}

	err = read(reader, &packet.Length)
	if err != nil {
		return nil, err
	}

	packet.Data = data[PACKET_RECORD_HDR_LEN:]

	return packet, nil
}

// DecodeExtraDataRecord decodes a raw extra data record into an
// ExtraDataRecord.
func DecodeExtraDataRecord(data []byte) (extra *ExtraDataRecord, err error) {

	extra = &ExtraDataRecord{}

	reader := bytes.NewBuffer(data)

	if err = read(reader, &extra.EventType); err != nil {
		return nil, err
	}

	if err = read(reader, &extra.EventLength); err != nil {
		return nil, err
	}

	if err = read(reader, &extra.SensorId); err != nil {
		return nil, err
	}

	if err = read(reader, &extra.EventId); err != nil {
		return nil, err
	}

	if err = read(reader, &extra.EventSecond); err != nil {
		return nil, err
	}

	if err = read(reader, &extra.Type); err != nil {
		return nil, err
	}

	if err = read(reader, &extra.DataType); err != nil {
		return nil, err
	}

	if err = read(reader, &extra.DataLength); err != nil {
		return nil, err
	}

	extra.Data = data[EXTRA_DATA_RECORD_HDR_LEN:]

	return extra, nil
}

// ReadRawRecord reads a raw record from the provided file.
func ReadRawRecord(file *os.File) (*RawRecord, error) {
	var header RawHeader

	/* Get the current offset so we can seek back to it. */
	offset, _ := file.Seek(0, 1)

	/* Now read in the header. */
	err := binary.Read(file, binary.BigEndian, &header)
	if err != nil {
		file.Seek(offset, 0)
		return nil, err
	}

	/* Create a buffer to hold the raw record data and read the
	/* record data into it */
	data := make([]byte, header.Len)
	len, err := file.Read(data)
	if uint32(len) < header.Len {
		/* Didn't read enough data. The error io.ErrShortBuffer seems
		/* to make sense here. */
		file.Seek(offset, 0)
		return nil, io.ErrShortBuffer
	} else if err != nil {
		file.Seek(offset, 0)
		return nil, err
	}

	return &RawRecord{header.Type, data}, nil
}

// ReadRecord reads and decodes a record from the provided file.
func ReadRecord(file *os.File) (*RecordContainer, error) {

	record, err := ReadRawRecord(file)
	if err != nil {
		return nil, err
	}

	var decoded interface{} = nil

	switch record.Type {
	case UNIFIED2_IDS_EVENT,
		UNIFIED2_IDS_EVENT_IP6,
		UNIFIED2_IDS_EVENT_V2,
		UNIFIED2_IDS_EVENT_IP6_V2:
		decoded, err = DecodeEventRecord(record.Type, record.Data)
	case UNIFIED2_PACKET:
		decoded, err = DecodePacketRecord(record.Data)
	case UNIFIED2_EXTRA_DATA:
		decoded, err = DecodeExtraDataRecord(record.Data)
	}

	if err != nil {
		return nil, err
	} else if decoded != nil {
		return &RecordContainer{record.Type, decoded}, nil
	} else {
		// Unknown record type.
		return nil, nil
	}
}
