package openprotocol

import "sync"

type ResponseQueue struct {
	Results map[string]string
	mtx     sync.Mutex
}

func (q *ResponseQueue) Add(mid string, msg string) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	q.Results[mid] = msg
}

func (q *ResponseQueue) update(mid string, msg string) {
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

func (q *ResponseQueue) get(mid string) string {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	return q.Results[mid]
}
