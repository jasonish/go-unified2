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

package unified2

import (
	"container/list"
	"fmt"
)

type Queue struct {
	queue *list.List
}

func NewQueue() *Queue {
	queue := new(Queue)
	queue.queue = list.New()
	return queue
}

func (q Queue) Append(record *Record) *Event {

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

func (q Queue) pop() *Record {
	element := q.queue.Front()
	q.queue.Remove(element)
	return element.Value.(*Record)
}

func (q Queue) Flush() *Event {

	record := q.pop()

	var event *Event
	var err error

	if IsEventType(record) {
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
