package cursor

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/vieolo/mansil"
)

// Records each write separately
type writes []string

func (w *writes) Write(p []byte) (int, error) {
	*w = append(*w, string(p))
	return len(p), nil
}

func TestPrintWritesOnceToW(t *testing.T) {
	var w writes
	Cursor{W: &w}.Print(AnsiSaveCursor, Ansi(mansil.CursorUp(2)), AnsiRestoreCursor)

	want := AnsiSaveCursor + "\033[2A" + AnsiRestoreCursor
	if len(w) != 1 || w[0] != want {
		t.Fatalf("writes = %q, want one write of %q", w, want)
	}
}

func TestZeroValueWritesToStdout(t *testing.T) {
	r, pw, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = pw
	Cursor{}.Print(mansil.ClearLine)
	os.Stdout = stdout
	pw.Close()

	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != mansil.ClearLine {
		t.Fatalf("stdout = %q, want %q", got, mansil.ClearLine)
	}
}

func TestClearPreviousLine(t *testing.T) {
	var b bytes.Buffer
	Cursor{W: &b}.ClearPreviousLine(2)

	want := AnsiSaveCursor + "\033[2A" + mansil.ClearLine + AnsiRestoreCursor
	if b.String() != want {
		t.Fatalf("got %q, want %q", b.String(), want)
	}
}

func TestReplaceLine(t *testing.T) {
	var b bytes.Buffer
	Cursor{W: &b}.ReplaceLine(1, "done")

	want := AnsiSaveCursor + "\033[1A" + mansil.ClearLine + mansil.CursorColumn1 + "done" + AnsiRestoreCursor
	if b.String() != want {
		t.Fatalf("got %q, want %q", b.String(), want)
	}
}

func TestReplaceLineZeroStaysOnTheCurrentLine(t *testing.T) {
	var b bytes.Buffer
	Cursor{W: &b}.ReplaceLine(0, "done")

	want := AnsiSaveCursor + mansil.ClearLine + mansil.CursorColumn1 + "done" + AnsiRestoreCursor
	if b.String() != want {
		t.Fatalf("got %q, want %q", b.String(), want)
	}
}

func TestMovementCodes(t *testing.T) {
	c := Cursor{}
	tests := []struct {
		name string
		got  Ansi
		want string
	}{
		{"up 3", c.ANSIToMoveUp(3), "\033[3A"},
		{"down 3", c.ANSIToMoveDown(3), "\033[3B"},
		{"right 3", c.MoveRight(3), "\033[3C"},
		{"left 3", c.MoveLeft(3), "\033[3D"},
		{"up 0", c.ANSIToMoveUp(0), ""},
		{"down 0", c.ANSIToMoveDown(0), ""},
		{"right -1", c.MoveRight(-1), ""},
		{"left 0", c.MoveLeft(0), ""},
	}
	for _, tt := range tests {
		if string(tt.got) != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}
