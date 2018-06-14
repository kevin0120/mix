package audi_vw

import "sync"

type ResponseQueue struct {
	Results map[uint32]string
	mtx     sync.Mutex
}

func (q *ResponseQueue) Add(serial uint32, msg string) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	q.Results[serial] = msg
}

func (q *ResponseQueue) update(serial uint32, msg string) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	_, e := q.Results[serial]
	if e {
		q.Results[serial] = msg
	}
}

func (q *ResponseQueue) HasResponse(serial uint32) bool {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	_, e := q.Results[serial]
	return e
}

func (q *ResponseQueue) remove(serial uint32) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	delete(q.Results, serial)
}

func (q *ResponseQueue) get(serial uint32) string {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	return q.Results[serial]
}
