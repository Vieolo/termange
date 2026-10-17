// Package live draws a region: a block of lines at the end of what has been
// written to the terminal, redrawn in place, for progress, spinners and
// status trees.
//
// While a region is live, it owns the end of the output. Anything else
// written to the same terminal between Start and Stop, such as logs, breaks
// the next redraw. Show it in the region, or write it after Stop
package live

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/x/term"
	"github.com/vieolo/mansil"
)

// The width used when the terminal's width is unknown
const defaultWidth = 80

// Render returns the region's lines for a terminal of width columns and
// height rows. A height of 0 means unknown. final is true for the last
// draw, which may be taller than the terminal since nothing redraws it.
//
// Lines wider than the terminal are cut, and a line that contains `\n` counts
// as several lines. Tabs should be avoided, since their width depends on the
// terminal's tab stops
type Render func(width, height int, final bool) []string

type Options struct {
	// Interval redraws the region on a timer, to move spinners and clocks.
	// Zero redraws only on Refresh
	Interval time.Duration
	// Size returns the terminal's columns and rows. By default, it is the
	// size of the writer when the writer is a terminal, and otherwise 80
	// columns with no row limit. A width of 0 or less counts as 80, and a
	// height of 0 or less as no row limit
	Size func() (width, height int)
	// HideCursor hides the cursor while the region is live, and shows it
	// again on Stop.
	//
	// It is off by default because a process that exits without calling
	// Stop, such as when it is killed or on an unhandled Ctrl-C, leaves the
	// cursor hidden in the user's terminal. Turn it on only when Stop runs
	// on every exit, such as from a signal handler
	HideCursor bool
}

// Region is a block of lines that is redrawn in place
type Region struct {
	w      io.Writer
	render Render
	opts   Options

	mu       sync.Mutex
	started  bool
	stopped  bool
	drawn    int           // The number of lines of the last frame
	stopTick chan struct{} // Closed by Stop to end the timer
	tickDone chan struct{} // Closed by the timer once it has ended
	stopOnce sync.Once
}

// New returns a region that writes to w and gets its lines from render.
// Nothing is drawn until Start
func New(w io.Writer, render Render, opts Options) *Region {
	if opts.Size == nil {
		opts.Size = writerSize(w)
	}
	return &Region{w: w, render: render, opts: opts}
}

// Start draws the region for the first time and starts the timer.
//
// The region begins on the cursor's line, which should be empty, at the
// start of a line. Calling Start again, or after Stop, does nothing
func (r *Region) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.started || r.stopped {
		return
	}
	r.started = true

	prefix := ""
	if r.opts.HideCursor {
		prefix = mansil.HideCursor
	}
	r.draw(false, prefix, "")

	if r.opts.Interval > 0 {
		r.stopTick = make(chan struct{})
		r.tickDone = make(chan struct{})
		go r.tick(r.opts.Interval, r.stopTick, r.tickDone)
	}
}

// Refresh redraws the region now. It is safe to call from any goroutine,
// and does nothing before Start or after Stop.
//
// The region calls render while holding its own lock, so render is never
// called twice at once. render must not call Refresh or Stop, since that
// would deadlock
func (r *Region) Refresh() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.started || r.stopped {
		return
	}
	r.draw(false, "", "")
}

// Stop stops the timer, draws the final frame and leaves it on screen.
// What is written to w afterwards appears below it.
//
// Nothing is written to w once Stop returns. Calling Stop again does
// nothing, and a region that was never started draws nothing
func (r *Region) Stop() {
	r.stopOnce.Do(r.stop)
}

func (r *Region) stop() {
	r.mu.Lock()
	r.stopped = true
	started, stopTick, tickDone := r.started, r.stopTick, r.tickDone
	r.mu.Unlock()
	if !started {
		return
	}

	// The timer is waited for without the lock, since a redraw of the
	// timer may be waiting for the lock. It finds the region stopped
	if stopTick != nil {
		close(stopTick)
		<-tickDone
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	suffix := ""
	if r.opts.HideCursor {
		suffix = mansil.ShowCursor
	}
	r.draw(true, "", suffix)
}

func (r *Region) tick(interval time.Duration, stop <-chan struct{}, done chan<- struct{}) {
	defer close(done)
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			r.Refresh()
		}
	}
}

// Draws a frame over the last one, in a single write to avoid flicker.
// It moves up to the first line of the last frame, erases down, and prints
// the new lines. r.mu must be held
func (r *Region) draw(final bool, prefix, suffix string) {
	width, height := r.opts.Size()
	if width <= 0 {
		width = defaultWidth
	}
	if height < 0 {
		height = 0
	}
	lines := fit(r.render(width, height, final), width, height, final)

	var b strings.Builder
	b.WriteString(prefix)
	b.WriteString(mansil.CursorPrevLine(r.drawn))
	b.WriteString(mansil.ClearScreenDown)
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteString(suffix)
	io.WriteString(r.w, b.String())
	r.drawn = len(lines)
}

// Returns the lines of a frame, each taking exactly one row.
//
// Every line is cut to the width, keeping its colours, since a wrapped line
// takes two rows and throws off the next redraw. A frame that is not final
// keeps at most height-1 lines, so that the next redraw can move up to its
// first line: the line below the frame, where the cursor waits, takes the
// last row. The last line kept then says how many lines were left out
func fit(lines []string, width, height int, final bool) []string {
	var rows []string
	for _, line := range lines {
		rows = append(rows, strings.Split(strings.TrimSuffix(line, "\n"), "\n")...)
	}

	if !final && height > 0 && len(rows) > height-1 {
		if height-1 <= 0 {
			return nil
		}
		keep := height - 2
		rows = append(rows[:keep], fmt.Sprintf("… %d more", len(rows)-keep))
	}

	for i, row := range rows {
		rows[i] = mansil.Truncate(row, width, "…")
	}
	return rows
}

// Returns the size of w when it is a terminal, and otherwise 80 columns
// with no row limit
func writerSize(w io.Writer) func() (int, int) {
	f, ok := w.(interface{ Fd() uintptr })
	if !ok || !term.IsTerminal(f.Fd()) {
		return func() (int, int) { return defaultWidth, 0 }
	}
	fd := f.Fd()
	return func() (int, int) {
		width, height, err := term.GetSize(fd)
		if err != nil {
			return defaultWidth, 0
		}
		return width, height
	}
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner returns the frame of a spinner for the time t, such as "⠹".
//
// The frame changes every 100ms. Since it depends only on the time, every
// spinner on screen moves together, and a fixed time gives a fixed frame
func Spinner(t time.Time) string {
	n := int64(len(spinnerFrames))
	i := t.UnixMilli() / 100 % n
	if i < 0 {
		i += n
	}
	return spinnerFrames[i]
}
