// Runs termange's checks in a fixed order, and stops at the first failure:
// formatting, build, vet, and the tests with the race detector.
//
//	go run ./internal/check
//
// It uses only the standard library, so it still runs, and reports the
// failure, when termange itself does not build
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var commands = []struct {
	name string
	args []string
}{
	{"build", []string{"go", "build", "./..."}},
	{"vet", []string{"go", "vet", "./..."}},
	{"test", []string{"go", "test", "-race", "-count=1", "./..."}},
}

func main() {
	root, err := moduleRoot()
	if err != nil {
		fail("root", err)
	}

	fmt.Println("==> format: gofmt -l .")
	unformatted, err := unformattedFiles(root)
	if err != nil {
		fail("format", err)
	}
	if len(unformatted) > 0 {
		fail("format", fmt.Errorf("not formatted, run gofmt -w on:\n  %s", strings.Join(unformatted, "\n  ")))
	}

	for _, c := range commands {
		fmt.Printf("==> %s: %s\n", c.name, strings.Join(c.args, " "))
		cmd := exec.Command(c.args[0], c.args[1:]...)
		cmd.Dir = root
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			fail(c.name, err)
		}
	}
	fmt.Println("All checks passed")
}

func fail(step string, err error) {
	fmt.Fprintf(os.Stderr, "FAILED: %s: %v\n", step, err)
	os.Exit(1)
}

// Returns the directory of the go.mod of the current directory's module
func moduleRoot() (string, error) {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		return "", err
	}
	gomod := strings.TrimSpace(string(out))
	if gomod == "" || gomod == os.DevNull {
		return "", fmt.Errorf("run it from inside the termange module")
	}
	return filepath.Dir(gomod), nil
}

// Returns the Go files under root that gofmt would change, like `gofmt -l`.
// Directories starting with . or _ are skipped, as the go command does
func unformattedFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if path != root && (strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(name, ".go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		formatted, err := format.Source(src)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if !bytes.Equal(src, formatted) {
			rel, _ := filepath.Rel(root, path)
			files = append(files, rel)
		}
		return nil
	})
	return files, err
}
