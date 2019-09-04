package openprotocol

import (
	"sync"
)

type ResponseQueue struct {
	Results map[string]interface{}
	mtx     sync.Mutex
}

func (q *ResponseQueue) Add(mid string, msg interface{}) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	q.Results[mid] = msg
}

func (q *ResponseQueue) update(mid string, msg interface{}) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	_, e := q.Results[mid]
	if e {
		q.Results[mid] = msg
	}
}

func (q *ResponseQueue) remove(mid string) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	delete(q.Results, mid)
}

func (q *ResponseQueue) get(mid string) interface{} {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	m, e := q.Results[mid]
	if e {
		return m
	} else {
		return nil
	}
}
