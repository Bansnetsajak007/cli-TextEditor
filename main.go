/*
	can't edit every line can only appded content in file

*/

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



