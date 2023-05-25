package mlog_test

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

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
	mlog.SetCustomFormat(func(logine mlog.LogLine) {
		fmt.Println("CUSTOM")
	})
	lines := print()
	for _, line := range lines {
		if !strings.HasPrefix(line, "CUSTOM") {
			t.Errorf("got %s, want %s", "no line contains CUSTOM", "CUSTOM")
		}
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
