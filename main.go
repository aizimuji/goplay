package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// Code Editor (Left Pane)
	editor := tview.NewTextArea()
	editor.SetPlaceholder("Enter Go code here...")
	editor.SetTitle("Code (Not Saved)")
	editor.SetBorder(true)

	// Track changes for Undo/Redo and Title
	editor.SetChangedFunc(func() {
		saveSnapshot(editor)
	})

	// Output View (Right Pane)
	outputView := tview.NewTextView()
	outputView.SetDynamicColors(true)
	outputView.SetRegions(true)
	outputView.SetWordWrap(true)
	outputView.SetTitle("Output")
	outputView.SetBorder(true)
	outputView.SetChangedFunc(func() {
		app.Draw()
	})

	// Footer
	footerLeft := tview.NewTextView().SetDynamicColors(true).SetText("Line: 1")
	footerCenter := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText("Press Ctrl+H to see all keys")

	footer := tview.NewFlex().
		AddItem(footerLeft, 15, 1, false).
		AddItem(footerCenter, 0, 1, false)

	// Layout
	mainFlex := tview.NewFlex().
		AddItem(editor, 0, 1, true).
		AddItem(outputView, 0, 1, false)

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, true).
		AddItem(footer, 1, 1, false)

	pages.AddPage("main", layout, true, true)

	// Update cursor position in footer
	editor.SetMovedFunc(func() {
		row, _, _, _ := editor.GetCursor()
		footerLeft.SetText(fmt.Sprintf("Line: %d", row+1))
	})

	// Help Popup
	helpTable := tview.NewTable().
		SetBorders(false).
		SetSelectable(false, false)

	shortcuts := []struct{ Key, Desc string }{
		{"Ctrl+r", "Run Code"},
		{"Ctrl+k", "Compile Code"},
		{"Ctrl+f", "Format Code"},
		{"Ctrl+s", "Save File"},
		{"Ctrl+o", "Open File"},
		{"Ctrl+n", "New File"},
		{"Ctrl+t", "Load Template"},
		{"Ctrl+l", "Clear Code"},
		{"Ctrl+z", "Undo"},
		{"Ctrl+y", "Redo"},
		{"Ctrl+Shift+q", "Quit"},
		{"Esc", "Close Help"},
	}

	for i, s := range shortcuts {
		helpTable.SetCell(i, 0, tview.NewTableCell(s.Key).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignRight).SetExpansion(1))
		helpTable.SetCell(i, 1, tview.NewTableCell("   "+s.Desc).SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft).SetExpansion(1))
	}

	helpTable.SetBorder(true).SetTitle("Shortcuts")

	helpFlex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(helpTable, 16, 1, true).
			AddItem(nil, 0, 1, false), 50, 1, true).
		AddItem(nil, 0, 1, false)

	helpTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEnter {
			pages.RemovePage("help")
			app.SetFocus(editor)
			return nil
		}
		return event
	})

	// Global Shortcuts
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Debug logging
		logDebug(fmt.Sprintf("Key: %v, Rune: %c, Mod: %v", event.Key(), event.Rune(), event.Modifiers()))

		// Quit: Ctrl+Shift+Q
		if event.Key() == tcell.KeyCtrlQ && event.Modifiers()&tcell.ModShift != 0 {
			app.Stop()
			return nil
		}

		// Help: Ctrl+H
		if event.Key() == tcell.KeyCtrlH {
			pages.AddPage("help", helpFlex, true, true)
			app.SetFocus(helpTable)
			return nil
		}

		// Map other shortcuts
		switch event.Key() {
		case tcell.KeyCtrlR:
			runCode(app, editor, outputView)
			return nil
		case tcell.KeyCtrlK: // Build
			compileCode(app, editor, outputView)
			return nil
		case tcell.KeyCtrlF:
			formatCode(app, editor, outputView)
			return nil
		case tcell.KeyCtrlS:
			saveFile(app, editor, outputView, pages)
			return nil
		case tcell.KeyCtrlO:
			openFile(app, editor, outputView, pages)
			return nil
		case tcell.KeyCtrlN:
			newFile(app, editor, outputView)
			return nil
		case tcell.KeyCtrlL:
			saveSnapshot(editor)
			editor.SetText("", false)
			return nil
		case tcell.KeyCtrlZ:
			undo(app, editor, outputView)
			return nil
		case tcell.KeyCtrlY:
			redo(app, editor, outputView)
			return nil
		case tcell.KeyCtrlT:
			loadTemplate(app, editor, outputView)
			return nil
		}

		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func logDebug(msg string) {
	tmpDir := os.TempDir()
	logFile := filepath.Join(tmpDir, "goplay_debug.log")
	f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(msg + "\n")
}
