package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
)

type Editor struct {
	filename    string
	lines       []string
	currentLine int
	cursorX     int
	changed     bool
}

func NewEditor(filename string) *Editor {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Creating new file:", filename)
		return &Editor{filename: filename, lines: []string{""}, currentLine: 0, cursorX: 0}
	}

	lines := strings.Split(string(content), "\n")
	return &Editor{filename: filename, lines: lines, currentLine: 0, cursorX: 0}
}

func (e *Editor) display() {
	fmt.Print("\033[2J")  // Clear screen
	fmt.Print("\033[H")   // Move cursor to top-left corner
	for i, line := range e.lines {
		if i == e.currentLine {
			fmt.Printf("> %d: %s\n", i+1, line)
		} else {
			fmt.Printf("  %d: %s\n", i+1, line)
		}
	}
	fmt.Println("\nCommands: Ctrl+S (save), Ctrl+Q (quit)")
	fmt.Println("          Arrow keys to move, Enter to insert new line")
	fmt.Printf("\033[%d;%dH", e.currentLine+1, e.cursorX+5) // Move cursor to current position
}

func (e *Editor) handleKey(char rune, key keyboard.Key) bool {
	switch key {
	case keyboard.KeyArrowUp:
		if e.currentLine > 0 {
			e.currentLine--
			if e.cursorX > len(e.lines[e.currentLine]) {
				e.cursorX = len(e.lines[e.currentLine])
			}
		}
	case keyboard.KeyArrowDown:
		if e.currentLine < len(e.lines)-1 {
			e.currentLine++
			if e.cursorX > len(e.lines[e.currentLine]) {
				e.cursorX = len(e.lines[e.currentLine])
			}
		}
	case keyboard.KeyArrowLeft:
		if e.cursorX > 0 {
			e.cursorX--
		}
	case keyboard.KeyArrowRight:
		if e.cursorX < len(e.lines[e.currentLine]) {
			e.cursorX++
		}
	case keyboard.KeyEnter:
		rightPart := e.lines[e.currentLine][e.cursorX:]
		e.lines[e.currentLine] = e.lines[e.currentLine][:e.cursorX]
		e.lines = append(e.lines[:e.currentLine+1], append([]string{rightPart}, e.lines[e.currentLine+1:]...)...)
		e.currentLine++
		e.cursorX = 0
		e.changed = true
	case keyboard.KeyBackspace, keyboard.KeyBackspace2:
		if e.cursorX > 0 {
			e.lines[e.currentLine] = e.lines[e.currentLine][:e.cursorX-1] + e.lines[e.currentLine][e.cursorX:]
			e.cursorX--
			e.changed = true
		} else if e.currentLine > 0 {
			e.cursorX = len(e.lines[e.currentLine-1])
			e.lines[e.currentLine-1] += e.lines[e.currentLine]
			e.lines = append(e.lines[:e.currentLine], e.lines[e.currentLine+1:]...)
			e.currentLine--
			e.changed = true
		}
	case keyboard.KeyCtrlS:
		e.saveFile()
	case keyboard.KeyCtrlQ:
		if e.changed {
			fmt.Print("Unsaved changes. Are you sure you want to quit? (y/n): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" {
				return false
			}
		}
		return true
	default:
		if char != 0 {
			e.lines[e.currentLine] = e.lines[e.currentLine][:e.cursorX] + string(char) + e.lines[e.currentLine][e.cursorX:]
			e.cursorX++
			e.changed = true
		}
	}
	return false
}

func (e *Editor) saveFile() {
	content := strings.Join(e.lines, "\n")
	err := os.WriteFile(e.filename, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
	} else {
		fmt.Println("File saved successfully")
		e.changed = false
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run editor.go <filename>")
		return
	}

	editor := NewEditor(os.Args[1])

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		editor.display()

		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if editor.handleKey(char, key) {
			break
		}
	}
}