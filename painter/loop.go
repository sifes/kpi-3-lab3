package painter

import (
    "image"
    "sync"

    "golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
    Receiver Receiver
    next     screen.Texture
    prev     screen.Texture
    Mq       MessageQueue
    stop     chan struct{}
    stopReq  bool
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
    l.next, _ = s.NewTexture(size)
    l.prev, _ = s.NewTexture(size)
    l.stop = make(chan struct{})
    go l.eventLoop()
}

// eventLoop обробляє операції з черги повідомлень
func (l *Loop) eventLoop() {
    for {
        select {
        case <-l.stop:
            return
        default:
            op := l.Mq.pull()
            if op == nil {
                continue
            }
            if update := op.Do(l.next); update && l.Receiver != nil {
                l.Receiver.Update(l.next)
                l.next, l.prev = l.prev, l.next
            }
        }
    }
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
    l.Mq.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
    l.stopReq = true
    close(l.stop)
}

// MessageQueue — експортована структура черги повідомлень
type MessageQueue struct {
    Queue   []Operation
    Mu      sync.Mutex
    Blocked chan struct{}
}

func (mq *MessageQueue) push(op Operation) {
    mq.Mu.Lock()
    defer mq.Mu.Unlock()
    mq.queue = append(mq.queue, op)
    if mq.Blocked != nil {
        close(mq.blocked)
        mq.blocked = nil
    }
}

func (mq *MessageQueue) pull() Operation {
    mq.Mu.Lock()
    defer mq.Mu.Unlock()

    for len(mq.Queue) == 0 {
        mq.blocked = make(chan struct{})
        mq.Mu.Unlock()
        <-mq.blocked
        mq.Mu.Lock()
    }

    op := mq.queue[0]
    mq.queue = mq.queue[1:]
    return op
}

func (mq *MessageQueue) empty() bool {
    mq.Mu.Lock()
    defer mq.Mu.Unlock()
    return len(mq.queue) == 0
}
