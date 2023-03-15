package errors

import (
	"bytes"
	"fmt"
	"github.com/iobrother/zoo/core/config"
	"runtime"
	"strings"
)

func (x *Error) Stack() string {
	return x.Metadata["_zoo_error_stack"]
}

var (
	// goRootForFilter is used for stack filtering in development environment purpose.
	goRootForFilter = runtime.GOROOT()
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.ReplaceAll(goRootForFilter, "\\", "/")
	}
}

type StackInfo struct {
	Service string `json:"service,omitempty"`
	Lines   []*StackLine
}

type StackLine struct {
	Func     string `json:"func,omitempty"`
	FileLine string `json:"line,omitempty"`
}

func formatStackInfo(info *StackInfo) string {
	var buffer = bytes.NewBuffer(nil)
	buffer.WriteString(fmt.Sprintf("service: %s\n", info.Service))
	space := "  "
	for i, line := range info.Lines {
		if i >= 9 {
			space = " "
		}
		buffer.WriteString(fmt.Sprintf(
			"  %d#%s%s\n      %s\n",
			i+1, space, line.Func, line.FileLine,
		))
	}

	return buffer.String()
}

func stacktrace() string {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	st := pcs[0:n]

	set := map[string]struct{}{}
	info := &StackInfo{Service: config.GetString("app.name")}
	for _, pc := range st {
		if fn := runtime.FuncForPC(pc - 1); fn != nil {
			file, line := fn.FileLine(pc - 1)
			// Avoid stack string like "`autogenerated`"
			if strings.Contains(file, "<") {
				continue
			}
			// Ignore GO ROOT paths.
			if goRootForFilter != "" &&
				len(file) >= len(goRootForFilter) &&
				file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}

			fileLine := fmt.Sprintf("%s:%d", file, line)
			if _, ok := set[fileLine]; !ok {
				set[fileLine] = struct{}{}
				info.Lines = append(info.Lines, &StackLine{
					Func:     fn.Name(),
					FileLine: fileLine,
				})
			}

		}
	}

	return formatStackInfo(info)
}