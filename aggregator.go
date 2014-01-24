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

import "container/list"

// EventAggregator is used to aggregate records into events.
type EventAggregator struct {
	buffer *list.List
}

// Pop returns and removes the first record from the aggregator.
func (ea EventAggregator) pop() *RecordContainer {
	record := ea.buffer.Front()
	ea.buffer.Remove(record)
	return record.Value.(*RecordContainer)
}

// NewEventAggregator creates a new EventAggregator.
func NewEventAggregator() *EventAggregator {
	aggregator := new(EventAggregator)
	aggregator.buffer = list.New()
	return aggregator
}

// Add adds a record to the event aggregated returning an array of
// records comprising a single event if the new record is the start of
// a new event.
func (ea EventAggregator) Add(record *RecordContainer) []*RecordContainer {

	var event []*RecordContainer = nil

	// Check if we need to flush.
	if IsEventType(record.Type) && ea.buffer.Len() > 0 {
		event = ea.Flush()
	}

	ea.buffer.PushBack(record)

	return event
}

// Len returns the number of records currently in the aggregator.
func (ea EventAggregator) Len() int {
	return ea.buffer.Len()
}

// Flush removes all records from the aggregator returning them as an
// array.
func (ea EventAggregator) Flush() []*RecordContainer {

	event := make([]*RecordContainer, ea.buffer.Len())

	for key, _ := range event {
		event[key] = ea.pop()
	}

	return event
}
