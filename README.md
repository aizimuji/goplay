# GoPlay CLI

[中文文档](README_CN.md)

**GoPlay CLI** is a terminal-based Go Playground that allows you to write, format, compile, and run Go code directly from your command line interface. It features a clean, split-pane TUI (Text User Interface) built with [tview](https://github.com/rivo/tview), offering a seamless coding experience without leaving your terminal.

## Features

*   **Split-Pane Interface**: Edit code on the left, view output on the right.
*   **Instant Execution**: Run (`go run`) or Compile (`go build`) your code with a single keystroke.
*   **File Management**: Open, Save, and Create new files easily.
*   **Template Support**: Quickly load a predefined `.template` file or use the built-in default to bootstrap your coding.
*   **Build Support**: Compile your code into a standalone executable (`Ctrl+b`).
*   **System Clipboard**: Seamless Copy/Paste support with the system clipboard.
*   **Editor Tools**:
    *   Auto-formatting using `go fmt`.
    *   Snapshot-based Undo/Redo.
    *   Real-time cursor line tracking.
    *   Modification status indicators.
*   **Built-in Help**: Access a quick reference of all shortcuts anytime (`F1` or `Alt+h`).

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
| **Ctrl+k** | Compile Check | Compiles the code to check for errors (no binary output). |
| **Ctrl+b** | Build Binary | Builds the code into a standalone executable. |
| **Ctrl+f** | Format | Formats the code using `go fmt`. |
| **Ctrl+s** | Save | Saves the current code to a file. |
| **Ctrl+o** | Open | Opens an existing file. |
| **Ctrl+n** | New | Clears the editor for a new file. |
| **Ctrl+t** | Template | Loads template (Custom `.template` or Default). |
| **Ctrl+l** | Clear | Clears the editor content. |
| **Ctrl+z** | Undo | Undoes the last action. |
| **Ctrl+y** | Redo | Redoes the last undone action. |
| **Ctrl+q** | Quit | Exits the application. |
| **Ctrl+c** | Copy | Copies selected text to system clipboard. |
| **F1** / **Alt+h** | Help | Shows the shortcuts popup. |
| **F2** / **Alt+[** | Shrink Window | Decreases the left editor window width. |
| **F3** / **Alt+]** | Grow Window | Increases the left editor window width. |
| **Esc** | Close | Closes popups or help menu. |

### Detailed Help

#### 1. Templates (Ctrl+t)
Press `Ctrl+t` to load a template.
- **Default Template**: If no `.template` file exists in the current working directory, the app loads a built-in "Hello Goplay" example.
- **Custom Template**: Create a `.template` file in your working directory to use your own boilerplate.

#### 2. Compiling & Building
- **Compile Check (`Ctrl+k`)**: Fast check for syntax errors. Does not leave a binary file.
- **Build Binary (`Ctrl+b`)**: Prompts for an output filename and builds a standalone `.exe` (on Windows) or binary.

#### 3. Window Adjustment
You can adjust the split ratio between the Editor and Output:
- **F2** (or `Alt+[`): Shrink the Editor pane.
- **F3** (or `Alt+]`): Grow the Editor pane.

#### 4. Copy & Paste
- **Copy**: Select text in the editor and press **`Ctrl+c`**. The text is copied to your system clipboard.
  - *Note*: Ensure the editor is focused. `Ctrl+c` will NOT exit the application.
- **Paste**: Use your system's standard paste shortcut (e.g., **`Ctrl+v`** or right-click) to paste text into the editor.

#### 5. Exiting
Use **`Ctrl+q`** to safely exit the application.

## License

[MIT](LICENSE)
