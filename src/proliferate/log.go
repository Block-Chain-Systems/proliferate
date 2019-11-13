package proliferate

import (
	"fmt"
	"runtime"

	"encoding/json"
)

type Message struct {
	Level    int    `json:"-"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
	Text     string `json:"text"`
}

// Log emits message
func (node *Node) Log(message Message) {
	n := *node
	l := n.Config.Logging

	if l.Enabled == false {
		return
	}

	pc := make([]uintptr, 15)
	c := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:c])
	frame, _ := frames.Next()

	message.File = frame.File
	message.Line = frame.Line
	message.Function = frame.Function

	if l.Console == true && message.Level <= l.Level {
		//text, _ := json.MarshalIndent(message, "", "  ")
		text := Prepare(message)
		LogEmit(LabelSeverity(message.Level), string(text))
	}
}

// DumpChain emits entire blockchain
func DumpChain(chain Chain) {
	text, _ := json.MarshalIndent(chain, "", "  ")
	LogEmit("chain", string(text))
}

// LabelSeverity converts Message.Level to string
func LabelSeverity(severity int) string {
	switch severity {
	case 0:
		return "fatal"
	case 1:
		return "error"
	case 2:
		return "warning"
	case 3:
		return "notice"
	case 4:
		return "verbose"
	case 5:
		return "noise"
	}
	return "message"
}

// Prepare returns json as indented string
func Prepare(message Message) string {
	text, _ := json.MarshalIndent(message, "", "  ")
	return string(text)
}

// LogEmit emits log to console
func LogEmit(label string, message string) {
	fmt.Printf("{\"%v\":%s}", label, message)
}

// LogError Directly emits error to console
func LogError(label string, message Message) {
	text := Prepare(message)
	LogEmit(LabelSeverity(1), text)
}
