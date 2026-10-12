package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os/exec"
	"time"
)

type request struct {
	Code string `json:"code"`
}

type result struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
	TimedOut bool   `json:"timed_out"`
}

const maxOutput = 64 * 1024

type cappedWriter struct {
	buf bytes.Buffer
}

func (c *cappedWriter) Write(p []byte) (int, error) {
	room := maxOutput - c.buf.Len()
	if room > len(p) {
		room = len(p)
	}
	if room > 0 {
		c.buf.Write(p[:room])
	}
	return len(p), nil
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python", "-c", req.Code)
	var stdout, stderr cappedWriter
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run() // exit code is handled in step 1.6

	res := result{
		Stdout:   stdout.buf.String(),
		Stderr:   stderr.buf.String(),
		TimedOut: ctx.Err() == context.DeadlineExceeded,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/run", runHandler)
	http.ListenAndServe(":8080", nil)
}
