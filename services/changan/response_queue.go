package changan

import "sync"

type ResponseQueue struct {
	results map[int]interface{}
	mtx     sync.Mutex
}

func (q *ResponseQueue) Init() {
	q.results = map[int]interface{}{}
	q.mtx = sync.Mutex{}
}

func (q *ResponseQueue) update(serial int, msg string) {
	defer q.mtx.Unlock()
	q.mtx.Lock()

	q.results[serial] = msg
}

func (q *ResponseQueue) getAndRemove(serial int) interface{} {
	defer q.mtx.Unlock()
	q.mtx.Lock()

	_, e := q.results[serial]
	if e {
		return q.results[serial]
	} else {
		return nil
	}
}
