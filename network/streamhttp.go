package network

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	. "github.com/Bensterriblescripts/Lib-Handlers/logging"
)

type FlushWriter struct {
	W io.Writer
	F http.Flusher
}

// Write forwards data to the underlying writer and flushes the response.
//
// Example:
//
//	_, _ = fw.Write([]byte("ping\n"))
func (fw FlushWriter) Write(p []byte) (int, error) {
	n, err := fw.W.Write(p)
	fw.F.Flush()
	return n, err
}

// CreateShellStream executes a command and streams output to the HTTP response.
//
// Example:
//
//	status := network.CreateShellStream(&w, "ping", "127.0.0.1")
func CreateShellStream(w *http.ResponseWriter, command string, arg string) string {
	f, ok := (*w).(http.Flusher)
	if !ok {
		return "Streaming not supported"
	}

	fw := FlushWriter{W: *w, F: f}
	cmd := exec.Command(command, arg)
	stdTraceOut := io.MultiWriter(os.Stdout, TraceLogFile)
	stdErrorOut := io.MultiWriter(os.Stderr, ErrorLogFile)
	cmd.Stdout = io.MultiWriter(stdTraceOut, fw)
	cmd.Stderr = io.MultiWriter(stdErrorOut, fw)

	if err := cmd.Start(); err != nil {
		return "Command start failed: " + err.Error()
	}

	if _, failed := ErrorExists(fmt.Fprintln(*w, "Started...")); failed {
		ErrorLog("Started print failed... Writer may be misconfigured.")
		return ""
	}
	f.Flush()

	if err := cmd.Wait(); err != nil {
		f.Flush()
	}

	f.Flush()

	return "Finished"
}
// CreateInternalStream returns trace and error writers that stream to the response.
//
// Example:
//
//	traceOut, errorOut := network.CreateInternalStream(&w)
func CreateInternalStream(w *http.ResponseWriter) (io.Writer, io.Writer) {
	f, ok := (*w).(http.Flusher)
	if !ok {
		return nil, nil
	}

	fw := FlushWriter{W: *w, F: f}
	traceOut := io.MultiWriter(TraceLogFile, fw)
	errorOut := io.MultiWriter(ErrorLogFile, fw)

	if _, failed := ErrorExists(fmt.Fprintln(*w, "Started...")); failed {
		ErrorLog("Started fprintln failed... Writer may be misconfigured.")
		return nil, nil
	}
	f.Flush()

	return traceOut, errorOut
}
