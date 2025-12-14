package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	// Initialize working directory
	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting CWD: %v\n", err)
		// Non-fatal, just continue with empty workingDir
	}

	app := tview.NewApplication()
	pages := tview.NewPages()

	// Code Editor (Left Pane)
	editor := tview.NewTextArea()
	editor.SetPlaceholder("Enter Go code here...")
	editor.SetTitle("Code (Not Saved)")
	editor.SetBorder(true)

	// Clipboard integration
	editor.SetClipboard(func(text string) {
		clipboard.WriteAll(text)
	}, func() string {
		text, _ := clipboard.ReadAll()
		return text
	})

	// ... (Existing code)

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
	footerRight := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight).SetText("Press F1 / Alt+H to see all keys")

	footer := tview.NewFlex().
		AddItem(footerLeft, 0, 1, false).
		AddItem(footerRight, 40, 1, false)

	// Window Split Ratio (50 means 50/50)
	splitRatio := 50

	// Layout
	mainFlex := tview.NewFlex()
	// Function to refresh layout based on splitRatio
	refreshLayout := func() {
		mainFlex.Clear()
		mainFlex.AddItem(editor, 0, splitRatio, true)
		mainFlex.AddItem(outputView, 0, 100-splitRatio, false)
	}
	refreshLayout() // Initial layout

	layout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, true).
		AddItem(footer, 1, 1, false)

	pages.AddPage("main", layout, true, true)

	// Update cursor position in footer
	updateFooter := func() {
		row, _, _, _ := editor.GetCursor()
		cwdText := ""
		if workingDir != "" {
			cwdText = fmt.Sprintf("  [CWD: %s]", workingDir)
		}
		footerLeft.SetText(fmt.Sprintf("Line: %d%s", row+1, cwdText))
	}
	editor.SetMovedFunc(updateFooter)
	// Initial update
	updateFooter()

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
		{"Ctrl+g", "Go to Line"},
		{"Ctrl+p", "Settings"},
		{"Ctrl+b", "Build App"},
		{"F2 / Alt+[", "Shrink Left"},
		{"F3 / Alt+]", "Grow Left"},
		{"Ctrl+q", "Quit"},
		{"Esc", "Close"},
		{"F1 / Alt+h", "Help"},
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
			AddItem(helpTable, 18, 1, true).
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
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Debug logging
		logDebug(fmt.Sprintf("Key: %v, Rune: %c, Mod: %v", event.Key(), event.Rune(), event.Modifiers()))

		// Quit: Ctrl+Q
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
			return nil
		}

		// Help: F1 or Alt+H
		if event.Key() == tcell.KeyF1 || ((event.Rune() == 'h' || event.Rune() == 'H') && event.Modifiers()&tcell.ModAlt != 0) {
			pages.AddPage("help", helpFlex, true, true)
			app.SetFocus(helpTable)
			return nil
		}

		// Resize Shortcuts
		// Shrink: F2 or Alt+[
		if event.Key() == tcell.KeyF2 ||
			(event.Key() == tcell.KeyCtrlLeftSq) || // Ctrl+[
			(event.Rune() == '[' && event.Modifiers()&tcell.ModAlt != 0) {
			if splitRatio > 10 {
				splitRatio -= 5
				refreshLayout()
			}
			return nil
		}
		// Grow: F3 or Alt+]
		if event.Key() == tcell.KeyF3 ||
			(event.Key() == tcell.KeyCtrlRightSq) || // Ctrl+]
			(event.Rune() == ']' && event.Modifiers()&tcell.ModAlt != 0) {
			if splitRatio < 90 {
				splitRatio += 5
				refreshLayout()
			}
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
		case tcell.KeyCtrlG:
			jumpToLine(app, editor, outputView, pages)
			return nil
		case tcell.KeyCtrlP:
			showSettings(app, editor, outputView, pages, updateFooter)
			return nil
		case tcell.KeyCtrlB:
			buildExecutable(app, editor, outputView, pages)
			return nil
		case tcell.KeyCtrlC:
			// Manually handle Copy to prevent app exit
			if editor.HasFocus() {
				selectedText, _, _ := editor.GetSelection()
				if selectedText != "" {
					clipboard.WriteAll(selectedText)
					outputView.SetText("Copied to clipboard.")
				}
			}
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
