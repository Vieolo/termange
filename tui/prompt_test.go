package tui

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vieolo/termange/cursor"
)

func TestTextInputReadsOneLinePerCall(t *testing.T) {
	// A strings.Reader hands over all of its content in one read, so a
	// buffered reader in the first call would take the second line too
	in := strings.NewReader("alice\nbob\n")
	var out bytes.Buffer

	first := TextInput(TextInputOptions{Prompt: "Name?", In: in, Out: &out})
	second := TextInput(TextInputOptions{Prompt: "Name?", In: in, Out: &out})

	if first != "alice" || second != "bob" {
		t.Fatalf("got %q and %q, want alice and bob", first, second)
	}
}

func TestTextInputLineEndings(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"newline", "alice\n", "alice"},
		{"windows", "alice\r\n", "alice"},
		{"no line ending at the end of input", "alice", "alice"},
		{"empty line", "\n", ""},
		{"no input", "", ""},
		{"spaces are kept", "  alice \n", "  alice "},
	}
	for _, tt := range tests {
		got := TextInput(TextInputOptions{Prompt: "Name?", In: strings.NewReader(tt.in), Out: &bytes.Buffer{}})
		if got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestTextInputWritesThePromptToOut(t *testing.T) {
	var out bytes.Buffer
	TextInput(TextInputOptions{Prompt: "Name?", In: strings.NewReader("alice\n"), Out: &out})

	// The input is not a terminal, so nothing echoes the Enter
	// and TextInput moves to the next line itself
	if want := "Name? \n"; out.String() != want {
		t.Fatalf("out = %q, want %q", out.String(), want)
	}
}

func TestTextInputPostValuePrint(t *testing.T) {
	var out bytes.Buffer
	TextInput(TextInputOptions{
		Prompt:         "Name?",
		In:             strings.NewReader("alice\n"),
		Out:            &out,
		PostValuePrint: func(s string) string { return "Name: " + s },
	})

	var replaced bytes.Buffer
	cursor.Cursor{W: &replaced}.ReplaceLine(1, "Name: alice")
	if want := "Name? \n" + replaced.String(); out.String() != want {
		t.Fatalf("out = %q, want %q", out.String(), want)
	}
}

func TestTextInputPostValuePrintSkippedWithoutInput(t *testing.T) {
	called := false
	TextInput(TextInputOptions{
		Prompt:         "Name?",
		In:             strings.NewReader(""),
		Out:            &bytes.Buffer{},
		PostValuePrint: func(s string) string { called = true; return s },
	})
	if called {
		t.Fatal("PostValuePrint was called although nothing was entered")
	}
}

func TestConfirmAnswers(t *testing.T) {
	tests := []struct {
		in   string
		def  bool
		want bool
	}{
		{"\n", true, true},
		{"\n", false, false},
		{"y\n", false, true},
		{"Y\n", false, true},
		{"yes\n", false, true},
		{"YES\n", false, true},
		{" y \n", false, true},
		{"y\r\n", false, true},
		{"y", false, true},
		{"n\n", true, false},
		{"N\n", true, false},
		{"no\n", true, false},
		{"No\n", true, false},
	}
	for _, tt := range tests {
		got := Confirm(ConfirmOptions{Prompt: "Write?", Default: tt.def, In: strings.NewReader(tt.in), Out: &bytes.Buffer{}})
		if got != tt.want {
			t.Errorf("answer %q with default %v: got %v, want %v", tt.in, tt.def, got, tt.want)
		}
	}
}

func TestConfirmPromptShowsTheDefault(t *testing.T) {
	tests := []struct {
		prompt string
		def    bool
		want   string
	}{
		{"Write?", true, "Write? [Y/n] \n"},
		{"Write?", false, "Write? [y/N] \n"},
		{"Write? ", true, "Write? [Y/n] \n"},
		// A hint already in the prompt is replaced with the one of the default
		{"Write? [Y/n]", false, "Write? [y/N] \n"},
		{"Write? [y/N] ", true, "Write? [Y/n] \n"},
	}
	for _, tt := range tests {
		var out bytes.Buffer
		Confirm(ConfirmOptions{Prompt: tt.prompt, Default: tt.def, In: strings.NewReader("\n"), Out: &out})
		if out.String() != tt.want {
			t.Errorf("prompt %q with default %v: out = %q, want %q", tt.prompt, tt.def, out.String(), tt.want)
		}
	}
}

func TestConfirmAsksAgainAfterAnUnknownAnswer(t *testing.T) {
	var out bytes.Buffer
	got := Confirm(ConfirmOptions{Prompt: "Write?", In: strings.NewReader("maybe\ny\n"), Out: &out})

	if !got {
		t.Fatal("got false, want true from the second answer")
	}
	if want := "Write? [y/N] \nWrite? [y/N] \n"; out.String() != want {
		t.Fatalf("out = %q, want %q", out.String(), want)
	}
}

func TestConfirmIsFalseWhenTheInputEnds(t *testing.T) {
	for _, in := range []string{"", "maybe\n", "maybe"} {
		got := Confirm(ConfirmOptions{Prompt: "Write?", Default: true, In: strings.NewReader(in), Out: &bytes.Buffer{}})
		if got {
			t.Errorf("input %q: got true, want false since nobody answered", in)
		}
	}
}

func TestConfirmPostConfirmationPrint(t *testing.T) {
	var out bytes.Buffer
	var printed []bool
	Confirm(ConfirmOptions{
		Prompt:  "Write?",
		Default: true,
		In:      strings.NewReader("\n"),
		Out:     &out,
		PostConfirmationPrint: func(v bool) string {
			printed = append(printed, v)
			return "Writing"
		},
	})

	if len(printed) != 1 || !printed[0] {
		t.Fatalf("PostConfirmationPrint got %v, want [true]", printed)
	}
	if !strings.HasSuffix(out.String(), "Writing"+cursor.AnsiRestoreCursor) {
		t.Fatalf("out = %q, want it to end with the replaced line", out.String())
	}
}
