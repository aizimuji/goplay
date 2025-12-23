package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	}

	// 1. Parameter Handling
	initialFileContent := ""
	if len(os.Args) > 1 {
		arg := os.Args[1]
		stat, err := os.Stat(arg)
		if err == nil {
			if stat.IsDir() {
				// If directory, change CWD
				os.Chdir(arg)
				workingDir, _ = os.Getwd()
			} else {
				// If file, prepare to open
				content, err := os.ReadFile(arg)
				if err == nil {
					initialFileContent = string(content)
					currentFilename, _ = filepath.Abs(arg)
					// If the file is in a different dir, switch CWD to it?
					// The requirement says:
					// "goplay d:\dev\go 如果是目录, 则打开该目录, 并将当前工作目录设置为该目录"
					// "goplay file.go 如果是文件, 则直接打开该文件" (Does not specify changing CWD, but usually editors do or don't. Let's stick to CWD is where you launched it, unless you opened a dir)
				}
			}
		}
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

	if initialFileContent != "" {
		editor.SetText(initialFileContent, false)
		// Update title will be handled by logic inside actions or initial update
	}

	// 不拦截编辑器的键盘输入，让 tview 自然处理

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
	footerRight := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight).SetText("Press F1 to see all keys")

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
		{"Ctrl+r / F5", "Run Code"},
		{"Ctrl+k / F6", "Compile Code"},
		{"Ctrl+f", "Format Code"},
		{"Ctrl+s", "Save File"},
		{"Ctrl+o", "Open File"},
		{"Ctrl+n", "New File"},
		{"Ctrl+t / F4", "Load Template"},
		{"Ctrl+l", "Clear Code"},
		{"Ctrl+z", "Undo"},
		{"Ctrl+y", "Redo"},
		{"Ctrl+g", "Go to Line"},
		{"Ctrl+p", "Settings"},
		{"Ctrl+b / F7", "Build App"},
		{"F2", "Shrink Left"},
		{"F3", "Grow Left"},
		{"Ctrl+q / F12", "Quit"},
		{"Esc", "Close"},
		{"F1", "Help"},
		{"Ctrl+p / F9", "Settings"},
		{"F10", "About"},
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

		// Quit: Ctrl+Q, F12
		if event.Key() == tcell.KeyCtrlQ || event.Key() == tcell.KeyF12 {
			app.Stop()
			return nil
		}

		// Help: F1
		if event.Key() == tcell.KeyF1 {
			pages.AddPage("help", helpFlex, true, true)
			app.SetFocus(helpTable)
			return nil
		}

		// About: F10
		if event.Key() == tcell.KeyF10 {
			version := "v0.3.0"
			authorName := "aizimuji"
			modal := tview.NewModal().
				SetText(fmt.Sprintf("goplay %s\ncommand line go playground\n\nauthor: %s\nlicense: MIT", version, authorName)).
				AddButtons([]string{"OK"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					pages.RemovePage("about")
					app.SetFocus(editor)
				})
			pages.AddPage("about", modal, true, true)
			app.SetFocus(modal)
			return nil
		}

		// Resize Shortcuts
		// Shrink: F2
		if event.Key() == tcell.KeyF2 {
			if splitRatio > 10 {
				splitRatio -= 5
				refreshLayout()
			}
			return nil
		}
		// Grow: F3
		if event.Key() == tcell.KeyF3 {
			if splitRatio < 90 {
				splitRatio += 5
				refreshLayout()
			}
			return nil
		}

		// Map other shortcuts
		switch event.Key() {
		case tcell.KeyCtrlR, tcell.KeyF5:
			runCode(app, editor, outputView)
			return nil
		case tcell.KeyCtrlK, tcell.KeyF6: // Check/Compile
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
		case tcell.KeyCtrlT, tcell.KeyF4:
			loadTemplate(app, editor, outputView)
			return nil
		case tcell.KeyCtrlG:
			jumpToLine(app, editor, outputView, pages)
			return nil
		case tcell.KeyCtrlP, tcell.KeyF9:
			showSettings(app, editor, outputView, pages, updateFooter)
			return nil
		case tcell.KeyCtrlB, tcell.KeyF7:
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

func insertAtCursor(editor *tview.TextArea, text string) {
	row, screenCol, _, _ := editor.GetCursor()
	content := editor.GetText()

	// 标准化换行符为 \n
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	lines := strings.Split(content, "\n")

	// Calculate byte offset to insertion point
	offset := 0
	if row < len(lines) {
		for i := 0; i < row; i++ {
			offset += len(lines[i]) + 1 // +1 for newline
		}

		line := lines[row]

		// 将屏幕列位置转换为字节索引
		// Tab 在屏幕上显示为多个空格（通常到下一个 4 的倍数位置）
		byteIndex := 0
		currentScreenCol := 0
		for i, r := range line {
			if currentScreenCol >= screenCol {
				byteIndex = i
				break
			}
			if r == '\t' {
				// Tab 跳到下一个 4 的倍数
				currentScreenCol = ((currentScreenCol / 4) + 1) * 4
			} else {
				currentScreenCol++
			}
			byteIndex = i + len(string(r)) // 移动到下一个字符
		}

		// 如果 screenCol 超过了行尾，byteIndex 就是整行长度
		if currentScreenCol < screenCol {
			byteIndex = len(line)
		}

		offset += byteIndex
	} else {
		offset = len(content)
	}

	if offset > len(content) {
		offset = len(content)
	}

	prefix := content[:offset]
	suffix := content[offset:]
	newContent := prefix + text + suffix

	editor.SetText(newContent, false)

	// 计算光标应该在的新位置
	newOffset := offset + len(text)

	// 只用 Select 定位光标，不滚动视图
	editor.Select(newOffset, newOffset)
}

func logDebug(msg string) {
	tmpDir := os.TempDir()
	logFile := filepath.Join(tmpDir, "goplay_debug.log")
	f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString(msg + "\n")
}
