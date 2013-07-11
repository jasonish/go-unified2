package unified2

import (
	"container/list"
	"fmt"
)

type Queue struct {
	queue *list.List
}

func NewQueue() (*Queue) {
	queue := new(Queue)
	queue.queue = list.New()
	return queue
}

func (q Queue) Append(record *Record) (*Event) {

	var event *Event = nil

	switch record.Type {

	case UNIFIED2_IDS_EVENT, UNIFIED2_IDS_EVENT_IP6, 
		UNIFIED2_IDS_EVENT_V2, UNIFIED2_IDS_EVENT_IP6_V2:
		if q.queue.Len() > 0 {
			event = q.Flush()
		}
		q.queue.PushBack(record)
	default:
		if q.queue.Len() == 0 {
			fmt.Println("Discarding non-event type while not in event context.")
		} else {
			q.queue.PushBack(record)
		}
	}

	return event
}

func (q Queue) pop() (*Record) {
	element := q.queue.Front()
	q.queue.Remove(element)
	return element.Value.(*Record)
}

func (q Queue) Flush() (*Event) {

	record := q.pop()

	var event *Event
	var err error

	if isEventType(record) {
		event, err = DecodeEvent(record)
	} else {
		panic("event type expected")
	}
	
	for {
		var packet *Packet
		var extra *ExtraData
		if q.queue.Len() == 0 {
			break
		}
		record := q.pop()
		switch record.Type {
		case UNIFIED2_PACKET:
			packet, err = DecodePacket(record)
			event.Packets.PushBack(packet)
		case UNIFIED2_EXTRA_DATA:
			extra, err = DecodeExtraData(record)
			event.ExtraData.PushBack(extra)
		default:
			panic("unexpted record type while in event context")
		}
	}

	_ = err

	return event
}

