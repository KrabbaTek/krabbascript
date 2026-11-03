package formatter

// this file is for the automatic semicolon system in KrabbaScript
// it's job is to automatically insert semicolons at the end of lines so you can omit semicolons in your code.

// the formatter — namely FormatDir and FormatFile

// - Reece Turner (reeceturner1358@gmail.com)

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Nothing really special, but it returns an array of strings and not a stupid array of bytes.
func read(path string) ([]string, error) {
	// read file
	lines, readError := os.ReadFile(path)

	if readError != nil {
		return nil, readError
	}

	// convert []byte to []string
	stringLines := strings.Split(strings.TrimSpace(string(lines)), "\n")

	return stringLines, nil
}

// receives lines without semicolons;
// returns lines with semicolons
// for example, this turns
// var foo: I32 = 0
// into
// var foo: I32 = 0;
func insertSemicolons(lines []string) []string {
	// the lines with semicolons we'll be returning.
	var newLines []string

	for _, line := range lines {
		// keep the leading tabspace, only strip trailing whitespace
		line = strings.TrimRight(line, " \t\r")

		// keep empty lines as they are
		if strings.TrimSpace(line) == "" {
			newLines = append(newLines, line)
			continue
		}

		// append lines without changes if they already end with a semicolon or don't need one
		if strings.HasSuffix(line, ";") ||
			strings.HasSuffix(line, "{") ||
			strings.HasSuffix(line, "}") ||
			strings.HasSuffix(line, "(") ||
			strings.HasSuffix(line, "[") ||
			strings.HasSuffix(line, ",") ||
			strings.HasSuffix(line, ":") ||
			strings.HasSuffix(line, "&&") ||
			strings.HasSuffix(line, "||") ||
			strings.HasSuffix(line, "+") ||
			strings.HasSuffix(line, "-") ||
			strings.HasSuffix(line, "*") ||
			strings.HasSuffix(line, "/") ||
			strings.HasSuffix(line, "%") ||
			strings.HasPrefix(strings.TrimSpace(line), "#") {
			newLines = append(newLines, line)
		} else {
			// append semicolon
			newLine := line + ";"
			newLines = append(newLines, newLine)
		}
	}

	return newLines
}

// format a file — read the file into lines, add semicolons to the lines, and overwrite the file with the new semicolons.
func FormatFile(file string) error {
	lines, err := read(file)

	if err != nil {
		return err
	}

	formattedLines := insertSemicolons(lines)

	// replace old lines with formatted lines
	// the line noise is because I have to convert []strings into []bytes
	writeError := os.WriteFile(file, []byte(strings.Join(formattedLines, "\n")), 0644)

	if writeError != nil {
		return writeError
	}

	return nil
}

// batch format an entire directory
// once again, I use the fancy word format, but I really just mean insert semicolons.
// recursively walks a directory looking for krabbascript (.ks) files and formats each file
func FormatDirectory(dir string) error {
	// extract only the KrabbaScript (.ks) files
	var filePaths []string

	// I just copied this from https://zetcode.com/golang/find-file/
	// Don't hate the player, hate the game.
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		// only append if path is a file, and is a KrabbaScript file
		if !d.IsDir() && filepath.Ext(path) == ".ks" {
			filePaths = append(filePaths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// now format each KrabbaScript file
	for _, file := range filePaths {
		err := FormatFile(file)

		if err != nil {
			return err
		}
	}

	// no error
	return nil
}
