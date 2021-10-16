// This queue pins a dag to Infura every 15 seconds if the queue has objects

package utils

import (
	"sync"
	"time"

	"github.com/beevik/timerqueue"
)

var q *timerqueue.Queue
var lasttime time.Time
var objects = map[int]*object{}

type object struct {
	data         []byte
	expectedHash string
	attempts     int
}

func (o *object) OnTimer(t time.Time) {
	go uploadToInfura(o.data, o.expectedHash, o.attempts)
}

func CreateInfuraQueue() {
	q = timerqueue.New()
	go initTicker()
}

func initTicker() {
	t := time.NewTicker(2 * time.Second)
	for range t.C {
		q.Advance(time.Now())
	}
}

var l = &sync.Mutex{}

func AddToInfuraDagPinQueue(data []byte, expectedHash string, attempts int) {
	l.Lock()
	defer l.Unlock()
	var t time.Time
	var next int

	timer, _ := q.PeekFirst()

	if timer == nil {
		t = time.Now()
		next = 0
	} else {
		t = lasttime.Add(2 * time.Second)
		next = len(objects) - 1
	}

	lasttime = t

	objects[next] = &object{
		data:         data,
		expectedHash: expectedHash,
		attempts:     attempts,
	}

	q.Schedule(objects[next], t)
}
