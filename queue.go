package mlog

import (
	"math"
	"sync"
)

type logQueue struct {
	mu           sync.Mutex
	queue        []LogLine
	output       func(LogLine) bool
	maxQueueSize int64
}

func newQueue(out func(LogLine) bool) *logQueue {
	queue := &logQueue{
		output:       out,
		queue:        make([]LogLine, 0),
		maxQueueSize: math.MaxInt64,
	}
	go queue.logLoop()
	return queue
}

func (q *logQueue) logLoop() {
	for {
		logline := q.dequeue()
		if logline == (LogLine{}) {
			continue
		}
		if !q.output(logline) {
			q.putback(logline)
		}
	}
}

func (q *logQueue) enqueue(logline LogLine) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, logline)
}

func (q *logQueue) dequeue() LogLine {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		return LogLine{}
	}
	logline := q.queue[0]
	q.queue = q.queue[1:]
	return logline
}

func (q *logQueue) putback(logline LogLine) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append([]LogLine{logline}, q.queue...)
}
