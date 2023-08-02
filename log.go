package proliferate

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"encoding/json"
	"io/ioutil"
)

// Message struct used by Log and LogRaw function
type Message struct {
	Level    int    `json:"-"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
	Caller   string `json:"caller"`
	Text     string `json:"text"`
}

// Log emits message
func (node *Node) Log(message Message) {
	n := *node
	l := n.Config.Logging
	wd, _ := os.Getwd()

	if l.Enabled == false {
		return
	}

	pc := make([]uintptr, 15)
	c := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:c])
	frame, _ := frames.Next()
	caller, _ := frames.Next()

	message.File = frame.File
	message.Line = frame.Line
	message.Function = frame.Function
	message.Caller = caller.Function

	text := string(Prepare(message))

	if l.Console == true && message.Level <= l.Level {
		LogEmit(LabelSeverity(message.Level),
			strings.Replace(text, wd, ".", -1))
	}

	if l.File == true && message.Level <= l.Level {
		LogToFile(text, l.FileLocation)
	}
}

// DumpChain emits entire blockchain
func DumpChain(chain Chain) {
	text, err := json.MarshalIndent(chain, "", "    ")

	if err != nil {
		LogRaw(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

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
	text, _ := json.MarshalIndent(message, "", "    ")
	return string(text)
}

// LogEmit emits log to console
func LogEmit(label string, message string) {
	fmt.Printf("\"%v\":%s\n", label, message)
}

// LogRaw directly emits error to console
func LogRaw(message Message) {
	text := Prepare(message)
	LogEmit(LabelSeverity(message.Level), text)
}

func LogToFile(text, filename string) error {

    // Write the text to the file
    err := ioutil.WriteFile(filename, []byte(text), 0644)
    if err != nil {
        return err
    }
    return nil
}

