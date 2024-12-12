package mailing

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
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

func SaveState() error {
	appStateDir := "./app_state"
	qFileName := "queue_state"
	err := os.Mkdir(appStateDir, 0750)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating directory %s: %w", appStateDir, err)
	}
	f, err := os.Create(fmt.Sprintf("%s/%s", appStateDir, qFileName))
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", appStateDir, err)
	}
	defer f.Close()
	ordersQ.mu.Lock()
	defer ordersQ.mu.Unlock()
	for k, v := range ordersQ.queue {
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("error encoding %+v as json: %w", v, err)
		}
		s := fmt.Sprintf("%s|%s\n", k, b)
		_, err = f.WriteString(s)
		if err != nil {
			return fmt.Errorf("error writing orders queue data to %s: %w", qFileName, err)
		}
	}
	return nil
}

func RecoverState() error {
	f, err := os.Open("./app_state/queue_state")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("error opening queue state file: %w", err)
	}
	defer f.Close()
	ordersQ.mu.Lock()
	defer ordersQ.mu.Unlock()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		keyVal := strings.SplitN(line, "|", 2)
		if len(keyVal) != 2 {
			return fmt.Errorf("error queue state file data format is invalid")
		}
		rawKey := keyVal[0]
		rawValue := keyVal[1]
		var value models.QueuedOrder
		err = json.Unmarshal([]byte(rawValue), &value)
		if err != nil {
			return fmt.Errorf("error parsing queue value from state file: %w", err)
		}
		ordersQ.queue[rawKey] = value
	}
	return nil
}
