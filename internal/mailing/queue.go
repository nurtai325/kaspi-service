package mailing

import (
	"sync"

	"github.com/nurtai325/kaspi-service/internal/models"
)

type ordersQueueContainer struct {
	mu    sync.Mutex
	queue map[string]models.QueuedOrder
}

var (
	ordersQ ordersQueueContainer
)

func init() {
	ordersQ.queue = make(map[string]models.QueuedOrder)
}

func (q *ordersQueueContainer) add(orderId string, qOrder models.QueuedOrder) {
	q.mu.Lock()
	defer q.mu.Unlock()
	ordersQ.queue[orderId] = qOrder
}
