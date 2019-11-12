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
		text, _ := json.Marshal(message)
		fmt.Printf("{\"%v\":%s}",
			LabelSeverity(message.Level),
			string(text))
	}
}

func DumpChain(chain Chain) {
	text, _ := json.Marshal(chain)
	fmt.Printf("{\"chain\":%s}", string(text))
}

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
	}
	return "message"
}
