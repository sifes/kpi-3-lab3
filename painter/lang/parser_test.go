package lang

import (
	"image/color"
	"strings"
	"testing"

	"github.com/sifes/kpi-3-lab3/painter"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse_BasicCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected []painter.Operation
	}{
		{
			name:    "white fill",
			command: "white",
			expected: []painter.Operation{
				painter.OperationFunc(painter.WhiteFill),
			},
		},
		{
			name:    "green fill",
			command: "green",
			expected: []painter.Operation{
				painter.OperationFunc(painter.GreenFill),
			},
		},
		{
			name:    "update",
			command: "update",
			expected: []painter.Operation{
				painter.UpdateOp,
			},
		},
		{
			name:    "reset",
			command: "reset",
			expected: []painter.Operation{
				painter.OperationFunc(painter.ResetScreen),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			ops, err := parser.Parse(strings.NewReader(tc.command))
			
			assert.NoError(t, err)
			assert.Len(t, ops, 1)
			// Перевіряємо тільки тип, бо не можемо порівнювати функції
			assert.IsType(t, tc.expected[0], ops[0])
		})
	}
}

func TestParser_Parse_StructCommands(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected painter.Operation
	}{
		{
			name:     "background rectangle",
			command:  "bgrect 0.1 0.1 0.9 0.9",
			expected: &painter.BgRectangle{X1: 40, Y1: 40, X2: 360, Y2: 360},
		},
		{
			name:     "figure",
			command:  "figure 0.5 0.5",
			expected: &painter.Figure{X: 200, Y: 200, C: color.RGBA{R: 255, A: 255}},
		},
		{
			name:     "move",
			command:  "move 0.1 0.1",
			expected: &painter.Move{X: 40, Y: 40},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			ops, err := parser.Parse(strings.NewReader(tc.command))
			
			assert.NoError(t, err)
			assert.NotEmpty(t, ops)

			// Структури можемо порівнювати за типом
			switch expected := tc.expected.(type) {
			case *painter.BgRectangle:
				bgRect, ok := ops[0].(*painter.BgRectangle)
				assert.True(t, ok)
				assert.Equal(t, expected.X1, bgRect.X1)
				assert.Equal(t, expected.Y1, bgRect.Y1)
				assert.Equal(t, expected.X2, bgRect.X2)
				assert.Equal(t, expected.Y2, bgRect.Y2)
			case *painter.Figure:
				figure, ok := ops[0].(*painter.Figure)
				assert.True(t, ok)
				assert.Equal(t, expected.X, figure.X)
				assert.Equal(t, expected.Y, figure.Y)
				assert.Equal(t, expected.C, figure.C)
			case *painter.Move:
				move, ok := ops[0].(*painter.Move)
				assert.True(t, ok)
				assert.Equal(t, expected.X, move.X)
				assert.Equal(t, expected.Y, move.Y)
			}
		})
	}
}

func TestParser_Parse_ComplexScript(t *testing.T) {
	script := `
white
bgrect 0.25 0.25 0.75 0.75
figure 0.5 0.5
green
figure 0.6 0.6
update
`
	parser := &Parser{}
	ops, err := parser.Parse(strings.NewReader(script))
	
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(ops), 5)
	
	// Перевіряємо, що останній є UpdateOp
	assert.Equal(t, painter.UpdateOp, ops[len(ops)-1])
}

func TestParser_Parse_InvalidCommands(t *testing.T) {
	tests := []struct {
		name    string
		command string
	}{
		{
			name:    "invalid command",
			command: "invalidcommand",
		},
		{
			name:    "bgrect with too few arguments",
			command: "bgrect 0.1 0.1",
		},
		{
			name:    "figure with too few arguments",
			command: "figure 0.5",
		},
		{
			name:    "move with too few arguments",
			command: "move 0.1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			_, err := parser.Parse(strings.NewReader(tc.command))
			
			assert.Error(t, err)
		})
	}
}