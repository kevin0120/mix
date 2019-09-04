package openprotocol

import (
	"context"
	"sync"
	"time"
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

//从返回结果队列中取数据，当取到非nil数据时立即返回，否则在规定时间到达时返回nil
func (q *ResponseQueue) Get(mid string, ctx context.Context) interface{} {
	ch := make(chan interface{})

	go func() {
		for {
			in := q.get(mid)
			if in != nil {
				ch <- in
				break
			}
			time.Sleep(REPLY_TIMEOUT)
		}
	}()

	select {
	case resp := <-ch:
		return resp
	case <-ctx.Done():
		return nil
	}

}
