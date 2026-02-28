package comms

import (
	"fmt"
	"sync"
)

// Message is a routed byte payload between two named participants.
type Message struct {
	SenderID string
	TargetID string
	Payload  []byte
}

// ReceiverFunc is called by Tick when a message arrives for the registered ID.
type ReceiverFunc func(msg Message)

// MessageBus queues outbound messages and delivers them synchronously on Tick.
type MessageBus struct {
	mu          sync.Mutex
	subscribers map[string]ReceiverFunc
	queue       []Message
}

func NewMessageBus() *MessageBus {
	return &MessageBus{
		subscribers: make(map[string]ReceiverFunc),
	}
}

// Subscribe registers a ReceiverFunc for the given ID.
func (b *MessageBus) Subscribe(id string, receiver ReceiverFunc) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[id] = receiver
}

// Send enqueues a message. The payload is copied defensively.
func (b *MessageBus) Send(senderID string, targetID string, payload []byte) {
	fmt.Printf("senderId: %s, targetId: %s\n", senderID, targetID)
	p := make([]byte, len(payload))
	copy(p, payload)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.queue = append(b.queue, Message{SenderID: senderID, TargetID: targetID, Payload: p})
}

// Tick drains the queue and delivers each message to its target subscriber.
// Messages with no registered target are silently dropped.
// TODO: dont silently drop, write to log
func (b *MessageBus) Tick() {
	b.mu.Lock()
	pending := b.queue
	b.queue = nil
	b.mu.Unlock()

	for _, msg := range pending {
		b.mu.Lock()
		receiver, ok := b.subscribers[msg.TargetID]
		b.mu.Unlock()
		if ok {
			receiver(msg)
		}
	}
}
