package Snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	serialBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	serialMax   int64 = -1 ^ (-1 << serialBits)
	timeShit    uint8 = workerBits + serialBits
	workerShift uint8 = serialBits
	epoch       int64 = 1606203979693
)

type Worker struct {
	mulock    sync.Mutex
	timestamp int64
	workerId  int64
	serial    int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("Worker ID is over maximum : " + strconv.FormatInt(workerMax, 10))
	}
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		serial:    0,
	}, nil
}

func (w *Worker) Generate() int64 {
	w.mulock.Lock()
	defer w.mulock.Unlock()

	now := time.Now().UnixNano() / 1000000
	if w.timestamp == now {
		w.serial++
		if w.serial > serialMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		w.serial = 0
		w.timestamp = now
	}
	ID := int64((now-epoch)<<timeShit | (w.workerId << workerShift) | w.serial)
	return ID
}
