package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

// Global state
var currentFilename string = ""
var workingDir string = ""
var isModified bool = false

// Undo/Redo
type EditorState struct {
	Text string
}

var undoStack []EditorState
var redoStack []EditorState
var isUndoingOrRedoing bool

func updateTitle(editor *tview.TextArea) {
	title := "Code"
	if currentFilename == "" {
		title += " (Not Saved)"
	} else {
		title += fmt.Sprintf(" (%s)", currentFilename)
	}

	if isModified {
		title += " (modified)"
	}

	editor.SetTitle(title)
}

func saveSnapshot(editor *tview.TextArea) {
	if isUndoingOrRedoing {
		return
	}
	if len(undoStack) > 1000 {
		undoStack = undoStack[1:]
	}
	undoStack = append(undoStack, EditorState{Text: editor.GetText()})
	redoStack = []EditorState{}

	isModified = true
	updateTitle(editor)
}

func undo(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView) {
	if len(undoStack) == 0 {
		outputView.SetText("Nothing to undo")
		return
	}

	isUndoingOrRedoing = true
	defer func() { isUndoingOrRedoing = false }()

	redoStack = append(redoStack, EditorState{Text: editor.GetText()})

	lastState := undoStack[len(undoStack)-1]
	undoStack = undoStack[:len(undoStack)-1]

	editor.SetText(lastState.Text, false)
	outputView.SetText("Undid last action")

	isModified = true
	updateTitle(editor)
}

func redo(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView) {
	if len(redoStack) == 0 {
		outputView.SetText("Nothing to redo")
		return
	}

	isUndoingOrRedoing = true
	defer func() { isUndoingOrRedoing = false }()

	undoStack = append(undoStack, EditorState{Text: editor.GetText()})

	nextState := redoStack[len(redoStack)-1]
	redoStack = redoStack[:len(redoStack)-1]

	editor.SetText(nextState.Text, false)
	outputView.SetText("Redid action")

	isModified = true
	updateTitle(editor)
}

func runCode(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView) {
	outputView.SetText("Running...")

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "goplay_run.go")

	err := os.WriteFile(tmpFile, []byte(editor.GetText()), 0644)
	if err != nil {
		outputView.SetText(fmt.Sprintf("Error writing temp file: %v", err))
		return
	}

	cmd := exec.Command("go", "run", tmpFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputView.SetText(fmt.Sprintf("Error running code:\n%s\n%v", string(output), err))
	} else {
		outputView.SetText(string(output))
	}
}

func compileCode(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView) {
	outputView.SetText("Compiling...")

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "goplay_build.go")

	err := os.WriteFile(tmpFile, []byte(editor.GetText()), 0644)
	if err != nil {
		outputView.SetText(fmt.Sprintf("Error writing temp file: %v", err))
		return
	}

	exeFile := filepath.Join(tmpDir, "goplay_run.exe")
	cmd := exec.Command("go", "build", "-o", exeFile, tmpFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputView.SetText(fmt.Sprintf("Compilation Error:\n%s\n%v", string(output), err))
	} else {
		outputView.SetText("Compilation Successful!")
	}
}

func formatCode(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView) {
	saveSnapshot(editor)

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "goplay_fmt.go")

	err := os.WriteFile(tmpFile, []byte(editor.GetText()), 0644)
	if err != nil {
		outputView.SetText(fmt.Sprintf("Error writing temp file: %v", err))
		return
	}

	cmd := exec.Command("go", "fmt", tmpFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		outputView.SetText(fmt.Sprintf("Format Error:\n%s\n%v", string(output), err))
		return
	}

	formatted, err := os.ReadFile(tmpFile)
	if err != nil {
		outputView.SetText(fmt.Sprintf("Error reading formatted file: %v", err))
		return
	}

	editor.SetText(string(formatted), false)
	outputView.SetText("Formatted code.")

	isModified = true
	updateTitle(editor)
}

func resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if workingDir != "" {
		return filepath.Join(workingDir, path)
	}
	return path
}

func saveFile(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView, pages *tview.Pages) {
	if currentFilename == "" {
		promptForInput(app, pages, "Save as:", "", func(filename string) {
			fullPath := resolvePath(filename)
			currentFilename = fullPath
			saveFileContent(app, editor, outputView, fullPath)
		})
	} else {
		saveFileContent(app, editor, outputView, currentFilename)
	}
}

func saveFileContent(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView, filename string) {
	err := os.WriteFile(filename, []byte(editor.GetText()), 0644)
	if err != nil {
		outputView.SetText(fmt.Sprintf("Error saving file: %v", err))
	} else {
		outputView.SetText(fmt.Sprintf("Saved to %s", filename))
		isModified = false
		updateTitle(editor)
	}
}

func openFile(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView, pages *tview.Pages) {
	promptForInput(app, pages, "Open file:", "", func(filename string) {
		fullPath := resolvePath(filename)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			outputView.SetText(fmt.Sprintf("Error opening file: %v", err))
			return
		}
		saveSnapshot(editor)

		currentFilename = fullPath
		editor.SetText(string(content), false)
		outputView.SetText(fmt.Sprintf("Opened %s", fullPath))

		isModified = false
		updateTitle(editor)
	})
}

func newFile(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView) {
	saveSnapshot(editor)
	currentFilename = ""
	editor.SetText("", false)
	outputView.SetText("New file created")

	isModified = false
	updateTitle(editor)
}

func promptForInput(app *tview.Application, pages *tview.Pages, title string, initialText string, callback func(string)) {
	form := tview.NewForm().
		AddInputField("Input", initialText, 30, nil, nil).
		AddButton("OK", nil).
		AddButton("Cancel", func() {
			pages.RemovePage("prompt")
			app.SetFocus(app.GetFocus())
		})

	form.GetButton(0).SetSelectedFunc(func() {
		text := form.GetFormItem(0).(*tview.InputField).GetText()
		if text != "" {
			callback(text)
		}
		pages.RemovePage("prompt")
	})

	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 7, 1, true).
			AddItem(nil, 0, 1, false), 40, 1, true).
		AddItem(nil, 0, 1, false)

	pages.AddPage("prompt", flex, true, true)
	app.SetFocus(form)
}

func loadTemplate(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView) {
	templateFile := ".template"
	if workingDir != "" {
		templateFile = filepath.Join(workingDir, ".template")
	}

	content, err := os.ReadFile(templateFile)
	if err != nil {
		if workingDir != "" {
			content, err = os.ReadFile(".template")
		}
		if err != nil {
			outputView.SetText(fmt.Sprintf("Error loading template: %v", err))
			return
		}
	}

	saveSnapshot(editor)
	editor.SetText(string(content), false)
	outputView.SetText("Loaded template.")

	isModified = true
	updateTitle(editor)
}

func jumpToLine(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView, pages *tview.Pages) {
	promptForInput(app, pages, "Go to Line:", "", func(input string) {
		lineNum, err := strconv.Atoi(input)
		if err != nil {
			outputView.SetText(fmt.Sprintf("Invalid line number: %s", input))
			return
		}

		targetRow := lineNum - 1
		if targetRow < 0 {
			targetRow = 0
		}

		text := editor.GetText()
		lines := strings.Split(text, "\n")

		if targetRow >= len(lines) {
			targetRow = len(lines) - 1
		}

		editor.SetOffset(targetRow, 0)

		offset := 0
		for i := 0; i < targetRow; i++ {
			offset += len(lines[i]) + 1
		}

		editor.Select(offset, offset)
		outputView.SetText(fmt.Sprintf("Jumped to line %d", targetRow+1))
	})
}

func showSettings(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView, pages *tview.Pages, updateFooter func()) {
	promptForInput(app, pages, "Set Working Directory:", workingDir, func(input string) {
		workingDir = input
		outputView.SetText(fmt.Sprintf("Working Directory set to: %s", workingDir))
		updateFooter()
	})
}

func buildExecutable(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView, pages *tview.Pages) {
	if isModified || currentFilename == "" {
		outputView.SetText("Please save the file before building.")
		return
	}

	baseName := filepath.Base(currentFilename)
	ext := filepath.Ext(baseName)
	defaultOutput := strings.TrimSuffix(baseName, ext)
	if runtime.GOOS == "windows" {
		defaultOutput += ".exe"
	}

	promptForInput(app, pages, "Build Output Name:", defaultOutput, func(outputName string) {
		outputView.SetText("Building executable...")

		outputPath := resolvePath(outputName)

		cmd := exec.Command("go", "build", "-o", outputPath, currentFilename)
		if workingDir != "" {
			cmd.Dir = workingDir
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			outputView.SetText(fmt.Sprintf("Build Error:\n%s\n%v", string(output), err))
		} else {
			outputView.SetText(fmt.Sprintf("Successfully built: %s", outputPath))
		}
	})
}
