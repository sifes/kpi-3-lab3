package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/shiny/screen"
)

// testReceiver імітує Receiver для тестування
type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

// mockTexture імітує screen.Texture для тестування
type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}

// mockScreen імітує screen.Screen для тестування
type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	l.Start(mockScreen{})
	l.Post(OperationFunc(WhiteFill))
	l.Post(OperationFunc(GreenFill))
	l.Post(UpdateOp)

	// Очікуємо, поки операції оброблюються
	time.Sleep(100 * time.Millisecond)

	// Перевірка
	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture type:", tr.lastTexture)
	}
	if len(mt.Colors) != 2 {
		t.Error("Unexpected number of colors:", mt.Colors)
	}

	// Завершення
	l.StopAndWait()
}

func TestLoop_Post_Multiple(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	l.Start(mockScreen{})

	// Create a separate channel to track execution order
	done := make(chan struct{})
	executed := make([]string, 0, 3)

	l.Post(OperationFunc(func(tx screen.Texture) {
		executed = append(executed, "op 1")
	}))

	l.Post(OperationFunc(func(tx screen.Texture) {
		executed = append(executed, "op 2")
	}))

	l.Post(OperationFunc(func(tx screen.Texture) {
		executed = append(executed, "op 3")
		close(done)
	}))

	// Wait for all operations to complete
	<-done

	// Check operations were executed in the right order
	expectedOrder := []string{"op 1", "op 2", "op 3"}
	if !reflect.DeepEqual(executed, expectedOrder) {
		t.Errorf("Expected execution order %v, got %v", expectedOrder, executed)
	}

	// Завершення
	l.StopAndWait()
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
