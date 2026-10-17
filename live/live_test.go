package live

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vieolo/mansil"
)

// Records each write separately
type writes struct {
	mu   sync.Mutex
	list []string
}

func (w *writes) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.list = append(w.list, string(p))
	return len(p), nil
}

func (w *writes) all() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]string(nil), w.list...)
}

func size(width, height int) func() (int, int) {
	return func() (int, int) { return width, height }
}

// Returns a render function that draws whatever lines are set, and the setter
func lines(initial ...string) (Render, func(...string)) {
	var mu sync.Mutex
	current := initial
	render := func(int, int, bool) []string {
		mu.Lock()
		defer mu.Unlock()
		return current
	}
	set := func(l ...string) {
		mu.Lock()
		defer mu.Unlock()
		current = l
	}
	return render, set
}

func TestFirstDrawDoesNotMoveAndLaterOnesMoveUp(t *testing.T) {
	var w writes
	render, set := lines("a", "b")
	r := New(&w, render, Options{Size: size(20, 0)})

	r.Start()
	set("x", "y", "z")
	r.Refresh()
	set("1")
	r.Refresh()

	want := []string{
		mansil.ClearScreenDown + "a\nb\n",
		mansil.CursorPrevLine(2) + mansil.ClearScreenDown + "x\ny\nz\n",
		mansil.CursorPrevLine(3) + mansil.ClearScreenDown + "1\n",
	}
	got := w.all()
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("writes =\n%q\nwant one write per frame:\n%q", got, want)
	}
}

func TestShorterFrameLeavesNothingOfTheLastOne(t *testing.T) {
	s := newScreen(20, 0)
	s.write("$ sirup test\n")
	render, set := lines("├─ app    running", "├─ core   running", "└─ cli    waiting")
	r := New(s, render, Options{Size: size(20, 0)})

	r.Start()
	set("all done")
	r.Refresh()

	want := []string{"$ sirup test", "all done"}
	if got := s.visible(); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("screen = %q, want %q", got, want)
	}
}

func TestWideColouredLinesAreCutKeepingTheirColours(t *testing.T) {
	var b bytes.Buffer
	render, _ := lines(mansil.RedFG + "abcdefghij" + mansil.Reset)
	r := New(&b, render, Options{Size: size(5, 0)})

	r.Start()

	want := mansil.ClearScreenDown + mansil.RedFG + "abcd…" + mansil.Reset + "\n"
	if b.String() != want {
		t.Fatalf("got %q, want %q", b.String(), want)
	}
}

func TestWideLinesDoNotBreakTheRedraw(t *testing.T) {
	s := newScreen(5, 0)
	render, set := lines("aaaaaaaaaa", "bbbbbbbbbb")
	r := New(s, render, Options{Size: size(5, 0)})

	r.Start()
	set("c")
	r.Refresh()

	if got := s.visible(); fmt.Sprint(got) != fmt.Sprint([]string{"c"}) {
		t.Fatalf("screen = %q, want [c]", got)
	}
}

func TestTallFrameIsCutUnlessFinal(t *testing.T) {
	var ten []string
	for i := range 10 {
		ten = append(ten, strconv.Itoa(i))
	}

	tests := []struct {
		height int
		final  bool
		want   []string
	}{
		{5, false, []string{"0", "1", "2", "… 7 more"}},
		{5, true, ten},
		{11, false, ten},
		{0, false, ten},
		{2, false, []string{"… 10 more"}},
		{1, false, nil},
	}
	for _, tt := range tests {
		got := fit(ten, 20, tt.height, tt.final)
		if fmt.Sprint(got) != fmt.Sprint(tt.want) {
			t.Errorf("height %d, final %v: got %q, want %q", tt.height, tt.final, got, tt.want)
		}
	}
}

func TestTallFramesRedrawWithinTheScreen(t *testing.T) {
	s := newScreen(20, 5)
	s.write("$ sirup test\n")
	render, set := lines()
	r := New(s, render, Options{Size: size(20, 5)})

	r.Start()
	for i := range 3 {
		var frame []string
		for j := range 8 {
			frame = append(frame, fmt.Sprintf("frame %d line %d", i, j))
		}
		set(frame...)
		r.Refresh()
	}

	// The 4 lines of the frame and the cursor's line below take the 5 rows,
	// so the prompt has scrolled off
	want := []string{"frame 2 line 0", "frame 2 line 1", "frame 2 line 2", "… 5 more"}
	if got := s.visible(); fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("screen = %q, want %q", got, want)
	}
}

func TestLinesWithNewlinesCountAsSeveralLines(t *testing.T) {
	var w writes
	render, _ := lines("a\nb", "c\n")
	r := New(&w, render, Options{Size: size(20, 0)})

	r.Start()
	r.Refresh()

	got := w.all()
	if got[0] != mansil.ClearScreenDown+"a\nb\nc\n" {
		t.Fatalf("first frame = %q", got[0])
	}
	if !strings.HasPrefix(got[1], mansil.CursorPrevLine(3)) {
		t.Fatalf("second frame = %q, want it to move up 3 lines", got[1])
	}
}

func TestStopDrawsTheFinalFrameAndStopsTheTimer(t *testing.T) {
	var w writes
	var renders, finals atomic.Int32
	r := New(&w, func(_, _ int, final bool) []string {
		renders.Add(1)
		if final {
			finals.Add(1)
		}
		return []string{"line"}
	}, Options{Size: size(20, 0), Interval: time.Millisecond})

	r.Start()
	deadline := time.Now().Add(2 * time.Second)
	for renders.Load() < 5 {
		if time.Now().After(deadline) {
			t.Fatalf("the timer redrew %d times in 2s", renders.Load())
		}
		time.Sleep(time.Millisecond)
	}
	r.Stop()
	afterStop := len(w.all())
	time.Sleep(20 * time.Millisecond)

	if finals.Load() != 1 {
		t.Fatalf("final frames = %d, want 1", finals.Load())
	}
	if n := len(w.all()); n != afterStop {
		t.Fatalf("%d writes after Stop", n-afterStop)
	}
}

func TestStopTwiceDrawsOneFinalFrame(t *testing.T) {
	var w writes
	render, _ := lines("a")
	r := New(&w, render, Options{Size: size(20, 0)})

	r.Start()
	r.Stop()
	r.Stop()
	r.Refresh()

	if got := w.all(); len(got) != 2 {
		t.Fatalf("writes = %q, want the first and the final frame", got)
	}
}

func TestNothingIsDrawnBeforeStartOrWithoutIt(t *testing.T) {
	var w writes
	called := false
	r := New(&w, func(int, int, bool) []string { called = true; return []string{"a"} }, Options{Size: size(20, 0)})

	r.Refresh()
	r.Stop()
	r.Start()
	r.Refresh()

	if got := w.all(); len(got) != 0 || called {
		t.Fatalf("writes = %q, render called = %v, want nothing", got, called)
	}
}

func TestConcurrentRefreshNeverRendersTwiceAtOnce(t *testing.T) {
	var b bytes.Buffer
	var rendering, overlaps atomic.Int32
	r := New(&b, func(int, int, bool) []string {
		if rendering.Add(1) > 1 {
			overlaps.Add(1)
		}
		defer rendering.Add(-1)
		return []string{"a", "b"}
	}, Options{Size: size(20, 0), Interval: time.Millisecond})

	r.Start()
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 50 {
				r.Refresh()
			}
		})
	}
	wg.Wait()
	r.Stop()

	if overlaps.Load() > 0 {
		t.Fatalf("render ran %d times while another call was running", overlaps.Load())
	}
}

func TestHideCursor(t *testing.T) {
	for _, hide := range []bool{true, false} {
		var w writes
		render, _ := lines("a")
		r := New(&w, render, Options{Size: size(20, 0), HideCursor: hide})

		r.Start()
		r.Refresh()
		r.Stop()

		got := w.all()
		first, last := got[0], got[len(got)-1]
		if strings.HasPrefix(first, mansil.HideCursor) != hide || strings.HasSuffix(last, mansil.ShowCursor) != hide {
			t.Errorf("HideCursor %v: first = %q, last = %q", hide, first, last)
		}
		if strings.Contains(got[1], mansil.HideCursor) || strings.Contains(got[1], mansil.ShowCursor) {
			t.Errorf("HideCursor %v: a redraw changed the cursor: %q", hide, got[1])
		}
	}
}

func TestDefaultSizeOfAWriterThatIsNotATerminal(t *testing.T) {
	var width, height int
	r := New(&bytes.Buffer{}, func(w, h int, _ bool) []string {
		width, height = w, h
		return nil
	}, Options{})

	r.Start()

	if width != 80 || height != 0 {
		t.Fatalf("size = %dx%d, want 80x0", width, height)
	}
}

func TestUnknownSizeFromOptions(t *testing.T) {
	var width, height int
	r := New(&bytes.Buffer{}, func(w, h int, _ bool) []string {
		width, height = w, h
		return nil
	}, Options{Size: size(0, -1)})

	r.Start()

	if width != 80 || height != 0 {
		t.Fatalf("size = %dx%d, want 80x0", width, height)
	}
}

func TestSpinner(t *testing.T) {
	tests := []struct {
		ms   int64
		want string
	}{
		{0, "⠋"},
		{99, "⠋"},
		{100, "⠙"},
		{250, "⠹"},
		{900, "⠏"},
		{1000, "⠋"},
		{-150, "⠏"},
	}
	for _, tt := range tests {
		if got := Spinner(time.UnixMilli(tt.ms)); got != tt.want {
			t.Errorf("Spinner(%dms) = %q, want %q", tt.ms, got, tt.want)
		}
	}

	seen := map[string]bool{}
	for ms := int64(0); ms < 1000; ms += 100 {
		seen[Spinner(time.UnixMilli(ms))] = true
	}
	if len(seen) != len(spinnerFrames) {
		t.Fatalf("one second shows %d frames, want %d", len(seen), len(spinnerFrames))
	}
}

// A minimal terminal, to check what a region leaves on screen.
//
// It understands text, `\n` as a new line at the first column (as a terminal
// shows output, which turns `\n` into `\r\n`), wrapping at the width,
// scrolling at the height, and the CSI sequences F (previous line) and
// J (erase down). Other CSI sequences, such as colours, are ignored, and
// every character takes one column. A height of 0 means no row limit
type screen struct {
	width, height int
	rows          [][]rune
	row, col      int
}

func newScreen(width, height int) *screen {
	return &screen{width: width, height: height, rows: [][]rune{nil}}
}

func (s *screen) Write(p []byte) (int, error) {
	s.write(string(p))
	return len(p), nil
}

func (s *screen) write(text string) {
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		switch r := runes[i]; {
		case r == '\033' && i+1 < len(runes) && runes[i+1] == '[':
			j := i + 2
			for j < len(runes) && (runes[j] < 0x40 || runes[j] > 0x7E) {
				j++
			}
			if j < len(runes) {
				s.csi(string(runes[i+2:j]), runes[j])
			}
			i = j
		case r == '\n':
			s.newLine()
		default:
			if s.col >= s.width {
				s.newLine()
			}
			for len(s.rows[s.row]) <= s.col {
				s.rows[s.row] = append(s.rows[s.row], ' ')
			}
			s.rows[s.row][s.col] = r
			s.col++
		}
	}
}

func (s *screen) newLine() {
	s.row, s.col = s.row+1, 0
	if s.row == len(s.rows) {
		s.rows = append(s.rows, nil)
	}
	if s.height > 0 && len(s.rows) > s.height {
		s.rows = s.rows[1:]
		s.row--
	}
}

func (s *screen) csi(params string, final rune) {
	switch final {
	case 'F':
		n, _ := strconv.Atoi(params)
		s.row, s.col = max(0, s.row-max(n, 1)), 0
	case 'J':
		if params == "" || params == "0" {
			if s.col < len(s.rows[s.row]) {
				s.rows[s.row] = s.rows[s.row][:s.col]
			}
			for i := s.row + 1; i < len(s.rows); i++ {
				s.rows[i] = nil
			}
		}
	}
}

// Returns the rows of the screen without trailing spaces or empty rows
func (s *screen) visible() []string {
	var out []string
	for _, row := range s.rows {
		out = append(out, strings.TrimRight(string(row), " "))
	}
	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return out
}
