package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rivo/tview"
)

// Global state
var currentFilename string = ""
var isModified bool = false

// Undo/Redo
type EditorState struct {
	Text string
	// Cursor position could be added if TextArea exposes it easily
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
	// Limit stack size if needed, e.g. 1000
	if len(undoStack) > 1000 {
		undoStack = undoStack[1:]
	}
	undoStack = append(undoStack, EditorState{Text: editor.GetText()})
	redoStack = []EditorState{} // Clear redo stack on new change

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

	// Current state goes to redo stack
	redoStack = append(redoStack, EditorState{Text: editor.GetText()})

	// Pop from undo stack
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

	// Current state goes to undo stack
	undoStack = append(undoStack, EditorState{Text: editor.GetText()})

	// Pop from redo stack
	nextState := redoStack[len(redoStack)-1]
	redoStack = redoStack[:len(redoStack)-1]

	editor.SetText(nextState.Text, false)
	outputView.SetText("Redid action")

	isModified = true
	updateTitle(editor)
}

func runCode(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView) {
	outputView.SetText("Running...")

	// Create a temporary file
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

	// Create a temporary file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "goplay_build.go")

	err := os.WriteFile(tmpFile, []byte(editor.GetText()), 0644)
	if err != nil {
		outputView.SetText(fmt.Sprintf("Error writing temp file: %v", err))
		return
	}

	// Build to a temporary executable
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
	// Snapshot before format
	saveSnapshot(editor)

	// Create a temporary file
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

	// Read back the formatted file
	formatted, err := os.ReadFile(tmpFile)
	if err != nil {
		outputView.SetText(fmt.Sprintf("Error reading formatted file: %v", err))
		return
	}

	editor.SetText(string(formatted), false)
	outputView.SetText("Formatted code.")

	// Formatting might change content, but semantically it's just formatting.
	// However, we treat it as modification.
	isModified = true
	updateTitle(editor)
}

func saveFile(app *tview.Application, editor *tview.TextArea, outputView *tview.TextView, pages *tview.Pages) {
	if currentFilename == "" {
		promptForFilename(app, pages, "Save as:", func(filename string) {
			currentFilename = filename
			saveFileContent(app, editor, outputView, filename)
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
	promptForFilename(app, pages, "Open file:", func(filename string) {
		content, err := os.ReadFile(filename)
		if err != nil {
			outputView.SetText(fmt.Sprintf("Error opening file: %v", err))
			return
		}
		// Snapshot before opening new file? Or maybe clear history?
		// Let's snapshot so we can undo opening a file if we want (restore previous content)
		saveSnapshot(editor)

		currentFilename = filename
		editor.SetText(string(content), false)
		outputView.SetText(fmt.Sprintf("Opened %s", filename))

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

func promptForFilename(app *tview.Application, pages *tview.Pages, title string, callback func(string)) {
	// Re-implementing prompt with a custom Flex to capture input
	form := tview.NewForm().
		AddInputField("Filename", "", 30, nil, nil).
		AddButton("OK", nil).
		AddButton("Cancel", func() {
			pages.RemovePage("prompt")
		})

	form.GetButton(0).SetSelectedFunc(func() {
		filename := form.GetFormItem(0).(*tview.InputField).GetText()
		if filename != "" {
			callback(filename)
		}
		pages.RemovePage("prompt")
	})

	form.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)

	// Center the form
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
	// Check for .template file in current directory
	templateFile := ".template"
	content, err := os.ReadFile(templateFile)
	if err != nil {
		outputView.SetText(fmt.Sprintf("Error loading template: %v", err))
		return
	}

	saveSnapshot(editor)
	editor.SetText(string(content), false)
	outputView.SetText("Loaded template.")

	isModified = true
	updateTitle(editor)
}
