package main

import (
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/rect"
	"github.com/aarzilli/nucular/style"
	"github.com/jviguy/brainfuck_interpreter/brainfuck_interpreter"
	"go.uber.org/atomic"
	"io/ioutil"
	"log"
	"os"
)

var content string

var filename string

func main() {

	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	if filename != "" {
		bs, _ := ioutil.ReadFile(filename)
		content = string(bs)
	} else {
		filename = "unknown file"
		content = "Welcome to BrainStorm a BrainStorm IDE!"
	}

	wnd := nucular.NewMasterWindow(0, "BrainStorm", updateFn())
	wnd.SetStyle(style.FromTheme(style.RedTheme, 2.0))
	wnd.Main()
}

func updateFn() func(window *nucular.Window) {
	var multilineTextEditor nucular.TextEditor
	multilineTextEditor.Buffer = []rune(content)
	var ibeam bool
	return func(w *nucular.Window) {
		w.Row(20).Dynamic(4)
		w.MenubarBegin()
		w.Label("Editing "+filename, "LT")
		if w.ButtonText("Rename") {
			var renamer nucular.TextEditor
			renamer.Buffer = []rune(filename)
			w.Master().PopupOpen("Rename File", nucular.WindowTitle|nucular.WindowClosable|nucular.WindowBorder|nucular.WindowMovable|nucular.WindowScalable, rect.Rect{W: 400, H: 400}, true, func(w *nucular.Window) {
				w.Row(0).Dynamic(1)
				renamer.Flags = nucular.EditMultiline | nucular.EditSelectable | nucular.EditClipboard
				if ibeam {
					renamer.Flags |= nucular.EditIbeamCursor
				}
				renamer.Edit(w)
				filename = string(renamer.Buffer)
			})
		}
		w.CheckboxText("I-Beam cursor", &ibeam)
		if w.ButtonText("Save File") {
			_ = ioutil.WriteFile(filename, []byte(string(multilineTextEditor.Buffer)), 0644)
		}
		if w.ButtonText("Run Code") {
			var reader, writer, _ = os.Pipe()
			var term nucular.TextEditor
			var ptr atomic.Value
			ptr.Store(uint16(0))
			//The said memory the interep will use.
			memory := make(map[uint16]uint8)
			b := &brainfuck_interpreter.BrainFucker{Memory: memory, Ptr: ptr, Stdout: writer, Stdin: reader}
			b.Run(string(multilineTextEditor.Buffer))
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			term.Buffer = []rune(string(buf[:n]))
			w.Master().PopupOpen("Terminal", nucular.WindowTitle|nucular.WindowClosable|nucular.WindowBorder|nucular.WindowMovable|nucular.WindowScalable, rect.Rect{W: 400, H: 400}, true, func(w *nucular.Window) {
				w.Row(0).Dynamic(1)
				term.Flags = nucular.EditMultiline | nucular.EditSelectable | nucular.EditClipboard
				if ibeam {
					term.Flags |= nucular.EditIbeamCursor
				}
				term.Edit(w)
			})
		}
		w.MenubarEnd()
		w.Row(0).Dynamic(1)
		multilineTextEditor.Flags = nucular.EditMultiline | nucular.EditSelectable | nucular.EditClipboard
		if ibeam {
			multilineTextEditor.Flags |= nucular.EditIbeamCursor
		}
		multilineTextEditor.Edit(w)
	}
}
