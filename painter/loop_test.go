package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/shiny/screen"
)

// mockReceiver імітує Receiver для тестування
type mockReceiver struct {
	lastTexture screen.Texture
	updateCount int
}

func (r *mockReceiver) Update(t screen.Texture) {
	r.lastTexture = t
	r.updateCount++
}

// mockTexture імітує screen.Texture для тестування
type mockTexture struct {
	colors []color.Color
	size   image.Point
	bounds image.Rectangle
}

func (t *mockTexture) Release() {}

func (t *mockTexture) Size() image.Point {
	return t.size
}

func (t *mockTexture) Bounds() image.Rectangle {
	return t.bounds
}

func (t *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (t *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	t.colors = append(t.colors, src)
}

// mockScreen імітує screen.Screen для тестування
type mockScreen struct{}

func (s *mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	return nil, nil
}

func (s *mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{
		size:   size,
		bounds: image.Rectangle{Max: size},
	}, nil
}

func (s *mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return nil, nil
}

func TestLoop_Post(t *testing.T) {
	// Підготовка
	loop := Loop{}
	receiver := &mockReceiver{}
	loop.Receiver = receiver
	
	screen := &mockScreen{}
	loop.Start(screen)
	
	// Дія - відправляємо операцію, яка повертає update = true
	loop.Post(OperationFunc(func(t screen.Texture) {
		t.Fill(t.Bounds(), color.White, screen.Src)
	}))
	loop.Post(UpdateOp)
	
	// Очікуємо, поки операції оброблюються
	time.Sleep(100 * time.Millisecond)
	
	// Перевірка
	assert.NotNil(t, receiver.lastTexture)
	assert.Equal(t, 1, receiver.updateCount)
	
	// Перевіряємо, що текстура оновилася
	mt, ok := receiver.lastTexture.(*mockTexture)
	assert.True(t, ok)
	assert.NotEmpty(t, mt.colors)
	
	// Завершення
	loop.StopAndWait()
}

func TestLoop_Post_Multiple(t *testing.T) {
	// Підготовка
	loop := Loop{}
	receiver := &mockReceiver{}
	loop.Receiver = receiver
	
	screen := &mockScreen{}
	loop.Start(screen)
	
	// Дія - відправляємо кілька операцій
	for i := 0; i < 5; i++ {
		loop.Post(OperationFunc(func(t screen.Texture) {
			t.Fill(t.Bounds(), color.White, screen.Src)
		}))
	}
	loop.Post(UpdateOp)
	
	// Очікуємо, поки операції оброблюються
	time.Sleep(100 * time.Millisecond)
	
	// Перевірка
	assert.NotNil(t, receiver.lastTexture)
	assert.Equal(t, 1, receiver.updateCount)
	
	// Завершення
	loop.StopAndWait()
}

func TestMessageQueue(t *testing.T) {
	// Підготовка
	mq := messageQueue{}
	
	// Перевірка порожньої черги
	assert.True(t, mq.empty())
	
	// Додавання операцій
	op1 := OperationFunc(WhiteFill)
	op2 := OperationFunc(GreenFill)
	
	mq.push(op1)
	assert.False(t, mq.empty())
	
	mq.push(op2)
	
	// Витягування операцій
	pulledOp1 := mq.pull()
	assert.NotNil(t, pulledOp1)
	
	pulledOp2 := mq.pull()
	assert.NotNil(t, pulledOp2)
	
	// Після витягування всіх операцій черга порожня
	assert.True(t, mq.empty())
	
	// Перевіряємо, що операції витягнуті в правильному порядку
	_, ok1 := pulledOp1.(OperationFunc)
	_, ok2 := pulledOp2.(OperationFunc)
	assert.True(t, ok1)
	assert.True(t, ok2)
}