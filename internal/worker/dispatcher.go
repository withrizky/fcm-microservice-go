package worker

import (
	"log"
	"sync"

	"fcm_microservice/internal/fcm"
	"fcm_microservice/internal/model"
)

type Dispatcher struct {
	JobQueue   chan model.FcmPayload
	MaxWorkers int
	FCM        *fcm.Client
	wg         sync.WaitGroup
}

func NewDispatcher(maxWorkers int, bufferSize int, fcmClient *fcm.Client) *Dispatcher {
	return &Dispatcher{
		JobQueue:   make(chan model.FcmPayload, bufferSize),
		MaxWorkers: maxWorkers,
		FCM:        fcmClient,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		d.wg.Add(1)
		go func(workerID int) {
			defer d.wg.Done()
			log.Printf("Worker FCM #%d siap.", workerID)

			for job := range d.JobQueue {
				err := d.FCM.Send(job)
				if err != nil {
					log.Printf("[Worker-%d] ❌ Gagal kirim ke %s: %v", workerID, job.Target, err)
				} else {
					log.Printf("[Worker-%d] ✅ Sukses kirim ke %s", workerID, job.Target)
				}
			}
			log.Printf("Worker FCM #%d berhenti.", workerID)
		}(i)
	}
}

func (d *Dispatcher) Stop() {
	close(d.JobQueue)
	d.wg.Wait()
}
