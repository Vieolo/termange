package termange

import (
	"testing"

	"github.com/vieolo/mansil"
)

func TestPaintText(t *testing.T) {
	tests := []struct {
		color Color
		want  string
	}{
		{ColorBlack, mansil.BlackFG + "hi" + mansil.Reset},
		{ColorRed, mansil.RedFG + "hi" + mansil.Reset},
		{ColorGreen, mansil.GreenFG + "hi" + mansil.Reset},
		{ColorYellow, mansil.YellowFG + "hi" + mansil.Reset},
		{ColorBlue, mansil.BlueFG + "hi" + mansil.Reset},
		{ColorMagenta, mansil.MagentaFG + "hi" + mansil.Reset},
		{ColorCyan, mansil.CyanFG + "hi" + mansil.Reset},
		{ColorWhite, mansil.Bold + mansil.WhiteFG + "hi" + mansil.Reset},
		{"RED", mansil.RedFG + "hi" + mansil.Reset},
		{"orange", "hi"},
		{"", "hi"},
	}
	for _, tt := range tests {
		if got := PaintText("hi", tt.color); got != tt.want {
			t.Errorf("PaintText(%q) = %q, want %q", tt.color, got, tt.want)
		}
	}
}
