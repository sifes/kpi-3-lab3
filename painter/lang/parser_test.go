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
		name    string
		command string
		check   func(t *testing.T, ops []painter.Operation)
	}{
		{
			name:    "white fill",
			command: "white",
			check: func(t *testing.T, ops []painter.Operation) {
				assert.Equal(t, 1, len(ops), "Expected 1 operation")
				_, ok := ops[0].(painter.OperationFunc)
				assert.True(t, ok, "Expected OperationFunc")
			},
		},
		{
			name:    "green fill",
			command: "green",
			check: func(t *testing.T, ops []painter.Operation) {
				assert.Equal(t, 1, len(ops), "Expected 1 operation")
				_, ok := ops[0].(painter.OperationFunc)
				assert.True(t, ok, "Expected OperationFunc")
			},
		},
		{
			name:    "update",
			command: "update",
			check: func(t *testing.T, ops []painter.Operation) {
				assert.Equal(t, 1, len(ops), "Expected 1 operation")
				assert.Equal(t, painter.UpdateOp, ops[0], "Expected UpdateOp")
			},
		},
		{
			name:    "reset",
			command: "reset",
			check: func(t *testing.T, ops []painter.Operation) {
				assert.Equal(t, 1, len(ops), "Expected 1 operation")
				_, ok := ops[0].(painter.OperationFunc)
				assert.True(t, ok, "Expected OperationFunc")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			ops, err := parser.Parse(strings.NewReader(tc.command))
			
			assert.NoError(t, err)
			tc.check(t, ops)
		})
	}
}

func TestParser_Parse_StructCommands(t *testing.T) {
	tests := []struct {
		name    string
		command string
		check   func(t *testing.T, ops []painter.Operation)
	}{
		{
			name:    "background rectangle",
			command: "bgrect 0.1 0.1 0.9 0.9",
			check: func(t *testing.T, ops []painter.Operation) {
				assert.Equal(t, 2, len(ops), "Expected 2 operations")
				
				// First op is the background color (ResetScreen)
				_, ok := ops[0].(painter.OperationFunc)
				assert.True(t, ok, "First op should be OperationFunc")
				
				// Second op is the background rectangle
				bgRect, ok := ops[1].(*painter.BgRectangle)
				assert.True(t, ok, "Second op should be BgRectangle")
				if ok {
					assert.Equal(t, 40, bgRect.X1)
					assert.Equal(t, 40, bgRect.Y1)
					assert.Equal(t, 360, bgRect.X2)
					assert.Equal(t, 360, bgRect.Y2)
				}
			},
		},
		{
			name:    "figure",
			command: "figure 0.5 0.5",
			check: func(t *testing.T, ops []painter.Operation) {
				assert.Equal(t, 2, len(ops), "Expected 2 operations")
				
				// First op is the background color (ResetScreen)
				_, ok := ops[0].(painter.OperationFunc)
				assert.True(t, ok, "First op should be OperationFunc")
				
				// Second op is the figure
				figure, ok := ops[1].(*painter.Figure)
				assert.True(t, ok, "Second op should be Figure")
				if ok {
					assert.Equal(t, 200, figure.X)
					assert.Equal(t, 200, figure.Y)
					assert.Equal(t, color.RGBA{R: 255, A: 255}, figure.C)
				}
			},
		},
		{
			name:    "move",
			command: "move 0.1 0.1\nfigure 0.5 0.5", // Need a figure to move
			check: func(t *testing.T, ops []painter.Operation) {
				assert.GreaterOrEqual(t, len(ops), 3, "Expected at least 3 operations")
				
				// First op is the background color (ResetScreen)
				_, ok := ops[0].(painter.OperationFunc)
				assert.True(t, ok, "First op should be OperationFunc")
				
				// Look for the Move operation
				foundMove := false
				for _, op := range ops {
					if move, ok := op.(*painter.Move); ok {
						foundMove = true
						assert.Equal(t, 40, move.X)
						assert.Equal(t, 40, move.Y)
						break
					}
				}
				assert.True(t, foundMove, "Should find a Move operation")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &Parser{}
			ops, err := parser.Parse(strings.NewReader(tc.command))
			
			assert.NoError(t, err)
			tc.check(t, ops)
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