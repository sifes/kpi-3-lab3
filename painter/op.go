package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, draw.Src)
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, draw.Src)
}

// BgRectangle малює чорний прямокутник на фоні.
type BgRectangle struct {
	X1, Y1, X2, Y2 int
}

func (op *BgRectangle) Do(t screen.Texture) bool {
	t.Fill(image.Rect(op.X1, op.Y1, op.X2, op.Y2), color.Black, draw.Src)
	return false
}

// Figure малює фігуру з центром у координатах (x, y).
type Figure struct {
	X, Y int
	C    color.RGBA
}

func (op *Figure) Do(t screen.Texture) bool {
	// Малюємо фігуру у вигляді перевернутої літери "Т"
	// Горизонтальна частина
	t.Fill(image.Rect(op.X-75, op.Y, op.X+75, op.Y-70), op.C, draw.Src)
	// Вертикальна частина
	t.Fill(image.Rect(op.X-30, op.Y-70, op.X+30, op.Y+70), op.C, draw.Src)
	return false
}

// Move переміщує всі фігури.
type Move struct {
	X, Y    int
	Figures []*Figure
}

func (op *Move) Do(t screen.Texture) bool {
	for i := range op.Figures {
		op.Figures[i].X += op.X
		op.Figures[i].Y += op.Y
	}
	return false
}

// ResetScreen очищає поточний стан текстури і заповнює її чорним кольором.
func ResetScreen(t screen.Texture) {
	t.Fill(t.Bounds(), color.Black, draw.Src)
}
