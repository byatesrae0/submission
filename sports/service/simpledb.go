package service

import (
	"sync"

	"google.golang.org/protobuf/proto"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

// simpleDB is just a super quick db implementation for listing/getting results (which are static).
type simpleDB struct {
	events []*sports.Event
	m      sync.Mutex
}

func (d *simpleDB) list() []*sports.Event {
	d.m.Lock()
	defer d.m.Unlock()

	eventsCopy := make([]*sports.Event, 0, len(d.events))

	for a := range d.events {
		eventsCopy = append(eventsCopy, proto.Clone(d.events[a]).(*sports.Event))
	}

	return eventsCopy
}

func (d *simpleDB) get(id int64) *sports.Event {
	d.m.Lock()
	defer d.m.Unlock()

	for a := range d.events {
		if d.events[a].Id == id {
			return proto.Clone(d.events[a]).(*sports.Event)
		}
	}

	return nil
}
