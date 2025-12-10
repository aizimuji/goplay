# GoPlay CLI

**GoPlay CLI** is a terminal-based Go Playground that allows you to write, format, compile, and run Go code directly from your command line interface. It features a clean, split-pane TUI (Text User Interface) built with [tview](https://github.com/rivo/tview), offering a seamless coding experience without leaving your terminal.

## Features

*   **Split-Pane Interface**: Edit code on the left, view output on the right.
*   **Instant Execution**: Run (`go run`) or Compile (`go build`) your code with a single keystroke.
*   **File Management**: Open, Save, and Create new files easily.
*   **Template Support**: Quickly load a predefined `.template` file to bootstrap your coding.
*   **Editor Tools**:
    *   Auto-formatting using `go fmt`.
    *   Snapshot-based Undo/Redo.
    *   Real-time cursor line tracking.
    *   Modification status indicators.
*   **Built-in Help**: Access a quick reference of all shortcuts anytime.

## Installation

### Prerequisites
*   [Go](https://go.dev/dl/) installed on your system.

### Build from Source

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/goplay.git
    cd goplay
    ```

2.  Build the application:
    ```bash
    go build -o goplay
    ```

3.  Run it:
    ```bash
    ./goplay
    ```

## Usage

Start the application by running the executable. You can start typing Go code immediately.

### Key Bindings

| Shortcut | Action | Description |
| :--- | :--- | :--- |
| **Ctrl+r** | Run | Runs the current code and shows output. |
| **Ctrl+k** | Compile | Compiles the code and reports errors/success. |
| **Ctrl+f** | Format | Formats the code using `go fmt`. |
| **Ctrl+s** | Save | Saves the current code to a file. |
| **Ctrl+o** | Open | Opens an existing file. |
| **Ctrl+n** | New | Clears the editor for a new file. |
| **Ctrl+t** | Template | Loads content from `.template` file. |
| **Ctrl+l** | Clear | Clears the editor content. |
| **Ctrl+z** | Undo | Undoes the last action. |
| **Ctrl+y** | Redo | Redoes the last undone action. |
| **Ctrl+h** | Help | Shows the shortcuts popup. |
| **Ctrl+Shift+q** | Quit | Exits the application. |

## Customization

You can create a `.template` file in the same directory as the executable. Pressing `Ctrl+t` will load the contents of this file into the editor, which is useful for keeping a boilerplate `main` function ready.

**Example `.template`:**
```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, GoPlay!")
}
```

## License

[MIT](LICENSE)
