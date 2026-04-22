package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	args := os.Args[1:]
	path := "/bin/sh"

	start := time.Now()

	cmd := exec.Command(path, args...)

	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr

	var err error

	if currErr := cmd.Start(); currErr != nil {
		err = errors.Join(err, currErr)
	}

	if currErr := cmd.Wait(); currErr != nil {
		err = errors.Join(err, currErr)
	}

	if err != nil {
		panic(err)
	}

	output := strings.Join([]string{out.String(), outErr.String()}, "\n")

	output = strings.TrimPrefix(output, "\n")
	output = strings.TrimSuffix(output, "\n")

	output = `<process_open path="` + path + `">` + "\n" +
		`<args>` + "\n" +
		strings.Join(args, "\n") + "\n" +
		`</args>` + "\n" +
		`</process_open>` + "\n" +
		`<process_output exec_time="` + formatDuration(time.Since(start)) + `">` + output + `</process_output>`

	fmt.Println(output)

}

// formatDuration returns a human-readable duration like "1ms", "42ms", etc.
func formatDuration(d time.Duration) string {
	ms := d.Milliseconds()
	if ms > 0 {
		return fmt.Sprintf("%dms", ms)
	}
	us := d.Microseconds()
	if us > 0 {
		return fmt.Sprintf("%dus", us)
	}
	return "0ms"
}
