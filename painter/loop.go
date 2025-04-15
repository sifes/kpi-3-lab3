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

  next screen.Texture // текстура, яка зараз формується
  prev screen.Texture // текстура, яка була відправлення останнього разу у Receiver

  Mq MessageQueue

  stop    chan struct{}
  stopReq bool
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
  l.next, _ = s.NewTexture(size)
  l.prev, _ = s.NewTexture(size)

  l.stop = make(chan struct{})
  l.stopReq = false

  // Запускаємо цикл подій у горутині
  go l.eventLoop()
}

// eventLoop обробляє операції з черги повідомлень
func (l *Loop) eventLoop() {
  for !l.stopReq {
    op := l.Mq.pull()
    if op == nil {
      continue
    }

    if update := op.Do(l.next); update {
      l.Receiver.Update(l.next)
      l.next, l.prev = l.prev, l.next
    }
  }
  close(l.stop)
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
  l.Mq.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
  l.Post(OperationFunc(func(screen.Texture) {
    l.stopReq = true
  }))
  <-l.stop
}

// MessageQueue — експортована структура черги повідомлень
type MessageQueue struct {
  queue   []Operation
  mu      sync.Mutex
  blocked chan struct{}
}

func (mq *MessageQueue) push(op Operation) {
  mq.mu.Lock()
  defer mq.mu.Unlock()

  mq.queue = append(mq.queue, op)
  if mq.blocked != nil {
    close(mq.blocked)
    mq.blocked = nil
  }
}

func (mq *MessageQueue) pull() Operation {
  mq.mu.Lock()
  defer mq.mu.Unlock()

  for len(mq.queue) == 0 {
    mq.blocked = make(chan struct{})
    mq.mu.Unlock()
    <-mq.blocked
    mq.mu.Lock()
  }

  op := mq.queue[0]
  mq.queue[0] = nil
  mq.queue = mq.queue[1:]
  return op
}

func (mq *MessageQueue) empty() bool {
  mq.mu.Lock()
  defer mq.mu.Unlock()
  return len(mq.queue) == 0
}
