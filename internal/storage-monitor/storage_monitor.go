package storagemonitor

import (
	"fmt"
	"log"
	"sync"

	"github.com/shirou/gopsutil/disk"
)

type StorageLimitEvent struct {
	Message string
}

type EventBroker struct {
	subscribers []chan StorageLimitEvent
	mu          sync.Mutex
}

var broker *EventBroker

func init() {
	broker = NewEventBroker()
}

func NewEventBroker() *EventBroker {
	return &EventBroker{}
}

func (b *EventBroker) Subscribe() chan StorageLimitEvent {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan StorageLimitEvent, 1) // Buffered channel
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *EventBroker) Broadcast(event StorageLimitEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, subscriber := range b.subscribers {
		select {
		case subscriber <- event:
		default:
			fmt.Println("Warning: subscriber channel is full. Event not sent.")
		}
	}
}

func checkDiskUsage() {
	const diskUsageThreshold = 80.0
	usage, err := disk.Usage("/")
	if err != nil {
		log.Fatalf("Error getting disk usage: %v", err)
	}

	currentUsagePercent := usage.UsedPercent
	fmt.Printf("Current disk usage: %.2f%%\n", currentUsagePercent)

	if currentUsagePercent > diskUsageThreshold {
		broker.Broadcast(StorageLimitEvent{Message: "Disk usage exceeds threshold"})
	}
}

func startLoggerSubscriber(broker *EventBroker) {
	logSub := broker.Subscribe()
	go func() {
		for event := range logSub {
			log.Printf("Logger: %s\n", event.Message)
		}
	}()
}

func startAlertSystemSubscriber(broker *EventBroker) {
	alertSub := broker.Subscribe()
	go func() {
		for event := range alertSub {
			fmt.Printf("Alert sent: %s\n", event.Message)
		}
	}()
}
