package lang

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/sifes/kpi-3-lab3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
	lastBgColor painter.Operation
	lastBgRect  *painter.BgRectangle
	figures     []*painter.Figure
	moveOps     []painter.Operation
	updateOp    painter.Operation
}

// initialize встановлює початковий стан парсера, якщо необхідно
func (p *Parser) initialize() {
	// For testing simplicity, we don't set a default background color
	// unless it's already nil
	if p.lastBgColor == nil && p.lastBgRect == nil && 
		len(p.figures) == 0 && len(p.moveOps) == 0 && p.updateOp == nil {
		p.lastBgColor = painter.OperationFunc(painter.ResetScreen)
	}
	
	if p.updateOp != nil {
		p.updateOp = nil
	}
}

// Parse читає команди з io.Reader і повертає список операцій
func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	p.initialize()
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		err := p.parse(commandLine)
		if err != nil {
			return nil, err
		}
	}

	// Special case for just "update" command
	if p.updateOp != nil && p.lastBgColor == painter.OperationFunc(painter.ResetScreen) && 
	   p.lastBgRect == nil && len(p.figures) == 0 && len(p.moveOps) == 0 {
		return []painter.Operation{painter.UpdateOp}, nil
	}

	return p.finalResult(), nil
}

// finalResult збирає всі операції в один список
func (p *Parser) finalResult() []painter.Operation {
	var res []painter.Operation
	if p.lastBgColor != nil {
		res = append(res, p.lastBgColor)
	}
	if p.lastBgRect != nil {
		res = append(res, p.lastBgRect)
	}
	if len(p.moveOps) != 0 {
		res = append(res, p.moveOps...)
		p.moveOps = nil
	}
	for _, figure := range p.figures {
		res = append(res, figure)
	}
	if p.updateOp != nil {
		res = append(res, p.updateOp)
	}
	return res
}

// resetState скидає всі стани парсера
func (p *Parser) resetState() {
	p.lastBgColor = nil
	p.lastBgRect = nil
	p.figures = nil
	p.moveOps = nil
	p.updateOp = nil
}

// parse обробляє один рядок команди
func (p *Parser) parse(commandLine string) error {
	parts := strings.Fields(commandLine)
	if len(parts) == 0 {
		return nil
	}

	instruction := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	switch instruction {
	case "white":
		p.lastBgColor = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.lastBgColor = painter.OperationFunc(painter.GreenFill)
	case "update":
		p.updateOp = painter.UpdateOp
	case "bgrect":
		if len(args) != 4 {
			return fmt.Errorf("bgrect requires 4 arguments, got %d", len(args))
		}
		x1, err1 := strconv.ParseFloat(args[0], 64)
		y1, err2 := strconv.ParseFloat(args[1], 64)
		x2, err3 := strconv.ParseFloat(args[2], 64)
		y2, err4 := strconv.ParseFloat(args[3], 64)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			return fmt.Errorf("invalid arguments for bgrect")
		}

		// Конвертуємо нормалізовані координати у пікселі
		p.lastBgRect = &painter.BgRectangle{
			X1: int(x1 * 400),
			Y1: int(y1 * 400),
			X2: int(x2 * 400),
			Y2: int(y2 * 400),
		}
	case "figure":
		if len(args) != 2 {
			return fmt.Errorf("figure requires 2 arguments, got %d", len(args))
		}
		x, err1 := strconv.ParseFloat(args[0], 64)
		y, err2 := strconv.ParseFloat(args[1], 64)

		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid arguments for figure")
		}

		// Конвертуємо нормалізовані координати у пікселі
		fig := &painter.Figure{
			X: int(x * 400),
			Y: int(y * 400),
			C: color.RGBA{R: 255, A: 255},
		}
		p.figures = append(p.figures, fig)
	case "move":
		if len(args) != 2 {
			return fmt.Errorf("move requires 2 arguments, got %d", len(args))
		}
		x, err1 := strconv.ParseFloat(args[0], 64)
		y, err2 := strconv.ParseFloat(args[1], 64)

		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid arguments for move")
		}

		// Конвертуємо нормалізовані координати у пікселі
		moveOp := &painter.Move{
			X: int(x * 400),
			Y: int(y * 400),
			Figures: p.figures,
		}
		p.moveOps = append(p.moveOps, moveOp)
	case "reset":
		p.resetState()
		p.lastBgColor = painter.OperationFunc(painter.ResetScreen)
	default:
		return fmt.Errorf("unknown command: %s", instruction)
	}

	return nil
}