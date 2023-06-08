package mlog_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dinifarb/mlog"
)

func TestINI(t *testing.T) {
	mlog.SetFormat(mlog.Fjson)
	mlog.Trace("Hello World, Trace!")
	mlog.Debug("Hello World, Debug!")
	mlog.Info("Hello World, Info!")
	mlog.Warn("Hello World, Warn!")
	mlog.Error("Hello World, Error!")
	mlog.SetFormat(mlog.Ftext)
}

func TestSetAppName(t *testing.T) {
	want := "[test]"
	mlog.SetAppName("test")
	lines := print()
	if !strings.Contains(lines[0], want) {
		t.Errorf("got %s, want %s", "no line contains [test] appname", want)
	}
}

func TestSetLevel(t *testing.T) {
	for i := 5; i <= 1; i-- {
		want := i
		mlog.SetLevel(mlog.Level(i))
		got := len(print())
		if got != want {
			t.Errorf("got %d lines, want %d", got, want)
		}
	}
}

func TestSetTimeFromat(t *testing.T) {
	want := "[0-9]{4}-[0-9]{2}-[0-9]{2}"
	mlog.SetTimeFormat("2006-01-02")
	lines := print()
	line := strings.Split(lines[0], "|")[0]
	line = strings.Replace(line, "[test]", "", 1)
	line = strings.TrimSpace(line)
	match, err := regexp.MatchString(want, line)
	if err != nil {
		t.Errorf("got %s, want %s", err.Error(), "pattern: "+want)
	}
	if !match {
		t.Errorf("got %s, want %s", line, "pattern: "+want)
	}
}

func TestSetFormat(t *testing.T) {
	mlog.SetFormat(mlog.Fjson)
	lines := print()
	for _, line := range lines {
		if !strings.HasPrefix(line, "{") {
			t.Errorf("got %s, want %s", "no line contains json", "{")
		}
	}
	mlog.SetFormat(mlog.Ftext)
}

func TestSetCustomFormat(t *testing.T) {
	mlog.SetCustomFormat(func(logline mlog.LogLine) string {
		return fmt.Sprintf("CUSTOM |%s|%s", logline.Level, logline.Message)
	})
	lines := print()
	for _, line := range lines {
		if !strings.HasPrefix(line, "CUSTOM") {
			t.Errorf("got %s, want %s", "no line contains CUSTOM", "CUSTOM")
		}
	}
}

func TestCustomOutput(t *testing.T) {
	mlog.SetLevel(mlog.Lerror)
	mlog.AddCustomOutput(func(logline mlog.LogLine) bool {
		fmt.Println(logline.AppName, logline.Message)
		return true
	})
	lines := print()
	for _, line := range lines {
		if !strings.HasPrefix(line, "[") {
			t.Errorf("got %s, want %s", "no line contains [", "[")
		}
	}
}

func TestCustomOutputPutBack(t *testing.T) {
	//start http service on port 1
	logs := []string{}
	failCount := 0
	go func() {
		count := 0
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			count++
			if count == 2 {
				//read body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error("failed at http req: " + err.Error())
				}
				logs = append(logs, string(body))
				w.WriteHeader(http.StatusOK)
				return
			}
			failCount++
			w.WriteHeader(http.StatusInternalServerError)
		})
		http.ListenAndServe(":1", mux)
	}()
	time.Sleep(300 * time.Millisecond)
	mlog.SetLevel(mlog.Lerror)
	mlog.AddCustomOutput(func(logline mlog.LogLine) bool {
		resp, err := http.Post("http://localhost:1", "text/plain", strings.NewReader(logline.Message))
		if err != nil {
			t.Error("failed at http req: " + err.Error())
		}
		return resp.StatusCode == http.StatusOK
	})
	mlog.Error("Hello World, Error!")
	time.Sleep(300 * time.Millisecond)
	if failCount != 1 {
		t.Errorf("got %d, want %d", failCount, 1)
	}
	if len(logs) != 1 {
		t.Errorf("got %d, want %d", len(logs), 1)
	}
}

func Test2CustomOutput(t *testing.T) {
	mlog.SetLevel(mlog.Lerror)
	mlog.SetFormat(mlog.Ftext)
	mlog.AddCustomOutput(func(logline mlog.LogLine) bool {
		fmt.Println("[CUSTOM1]")
		return true
	})
	mlog.AddCustomOutput(func(logline mlog.LogLine) bool {
		fmt.Println("[CUSTOM2]")
		return true
	})
	time.Sleep(2 * time.Second)
	lines := print()
	want := 3
	got := len(lines)
	if got != want {
		t.Errorf("got %d lines, want %d", got, want)
	}
}

func print() []string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	mlog.Trace("Testing Trace")
	mlog.Debug("Testing Debug")
	mlog.Info("Testing Info")
	mlog.Warn("Testing Warn")
	mlog.Error("Testing Error")
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	var resLines []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[") {
			resLines = append(resLines, line)
		}
		if strings.HasPrefix(line, "{") {
			resLines = append(resLines, line)
		}
	}
	return resLines
}
