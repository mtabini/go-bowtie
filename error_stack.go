package bowtie

import (
	"bytes"
	"io/ioutil"
	"runtime"
)

// Stack trace code heavily borrowed from https://github.com/go-martini/martini/blob/master/recovery.go

// Struct StackFrame represents a frame of a stack trace
type StackFrame struct {
	Path   string `json:"path"`
	Line   int    `json:"line"`
	Func   string `json:"func"`
	Source string `json:"source"`
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []StackFrame {
	result := []StackFrame{}

	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string

	maxCount := 100

	for i := skip; i < skip+maxCount; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)

		if !ok {
			break
		}

		// Print this much at least.  If we can't find the source, it won't show.

		frame := StackFrame{
			Path: file,
			Line: line,
		}

		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		frame.Func = string(function(pc))
		frame.Source = string(source(lines, line))

		result = append(result, frame)
	}

	return result
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func (e *ErrorInstance) CaptureStackTrace() Error {
	e.stackTrace = stack(2)

	return e
}
