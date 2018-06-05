package audi_vw

import "sync"

type ResponseQueue struct {
	Results		map[uint][]byte
	mtx			sync.Mutex
}

func (q *ResponseQueue) update(serial uint, msg []byte) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	q.Results[serial] = msg
}

func (q *ResponseQueue) remove(serial uint) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	delete(q.Results, serial)
}

func (q *ResponseQueue) get(serial uint) ([]byte) {
	defer q.mtx.Unlock()

	q.mtx.Lock()
	return q.Results[serial]
}