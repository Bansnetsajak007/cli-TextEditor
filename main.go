/*
	can't edit every line can only appded content in file

*/

/*
		V1
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("USage: go run main.go <filename>")
		return
	}

	filename := os.Args[1]
	content , err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Creating new file:", filename)
		content = []byte{}
	}

	//splits the file content into individual lines.
	lines := strings.Split(string(content), "\n")

	fmt.Println("CLI Text Editor")
	fmt.Println("Commands: :w (save), :q (quit), :wq (save and quit)")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		for i, line := range lines {
			fmt.Printf("%d: %s\n", i+1, line)
		}
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case ":q":
			return
		case ":w":
			saveFile(filename, lines)
		case ":wq":
			saveFile(filename, lines)
			return
		default:
			lines = append(lines, input)
		}
	}
}

func saveFile(filename string, lines []string) {
	content := strings.Join(lines, "\n")
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
	} else {
		fmt.Println("File saved successfully")
	}
}

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Editor struct {
	filename string
	lines []string
	currentLine int
	changed bool
}

func NewEditor(filename string) *Editor {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Creating new file:", filename)
		return &Editor{filename: filename, lines: []string{""}, currentLine: 0}
	}

	lines := strings.Split(string(content), "\n")
	return &Editor{filename: filename, lines: lines, currentLine: 0}
}

//method on struct
func (e * Editor) display() {
	fmt.Print("\033[2J")  // Clear screen
	fmt.Print("\033[H")   // Move cursor to top-left corner
	for i, line := range e.lines {
		if i == e.currentLine {
			fmt.Printf("> %d: %s\n", i+1, line)
		} else {
			fmt.Printf("  %d: %s\n", i+1, line)
		}
	}
	fmt.Println("\nCommands: :w (save), :q (quit), :wq (save and quit)")
	fmt.Println("          :u (move up), :d (move down)")
}

func (e * Editor) handelCommand(cmd string) bool {
	switch cmd {
		case ":q":
			if e.changed {
			fmt.Print("Unsaved changes. Are you sure you want to quit? (y/n): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" {
				return false
			}
		}
			return true
	case ":w":
		e.saveFile()
	case ":wq":
		e.saveFile()
		return true
	case ":u":
		if e.currentLine > 0 {
			e.currentLine--
		}
	case ":d":
		if e.currentLine < len(e.lines) {
			e.currentLine++
		}
	default:
		e.lines[e.currentLine] = cmd
		e.changed = true
		if e.currentLine == len(e.lines)-1 {
			e.lines = append(e.lines, "")
		}
		e.currentLine++
}
	return false
}

func (e *Editor) saveFile(){
	content := strings.Join(e.lines, "\n")
	err := os.WriteFile(e.filename, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error saving file:", err)
	} else {
		fmt.Printf("File saved successfully")
		e.changed = false
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run editor.go <filename>")
		return	
	}

	editor := NewEditor(os.Args[1])
	scanner := bufio.NewScanner(os.Stdin)

	for {
		editor.display()
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if editor.handelCommand(input) {
			break
		}
	}
}