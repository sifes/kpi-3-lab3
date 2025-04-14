package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

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
	if tr.LastTexture != nil {
		t.Fatal("Receiver got the texture too early")
	}

	l.Post(OperationFunc(GreenFill))
	l.Post(UpdateOp)

	l.Post(OperationFunc(GreenFill))
	l.Post(UpdateOp)

	var testOps []string

	l.Post(OperationFunc(func(screen.Texture) {
		testOps = append(testOps, "op1")

		l.Post(OperationFunc(func(screen.Texture) {
			testOps = append(testOps, "op2")
		}))

	}))
	l.Post(OperationFunc(func(screen.Texture) {
		testOps = append(testOps, "op3")
	}))

	l.StopAndWait()

	tx, ok := tr.LastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Receiver still has not texture")
	}
	if tx.FillCnt != 3 {
		t.Error("Unexpected number of Fill calls:", tx.FillCnt)
	}

	if !reflect.DeepEqual(testOps, []string{"op1", "op3", "op2"}) {
		t.Error("Bad order of operations:", testOps)
	}

}

type testReceiver struct {
	LastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.LastTexture = t
}

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

type mockTexture struct {
	FillCnt int
}

func (m *mockTexture) Release()                {}
func (m *mockTexture) Size() image.Point       { return size }
func (m *mockTexture) Bounds() image.Rectangle { return image.Rectangle{Max: size} }

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	panic("implement me")
}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.FillCnt++
}
