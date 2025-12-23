# GoPlay CLI

[English](README.md)

**GoPlay CLI** is a terminal-based Go Playground that allows you to write, format, compile, and run Go code directly from your command line interface. It features a clean, split-pane TUI (Text User Interface) built with [tview](https://github.com/rivo/tview), offering a seamless coding experience without leaving your terminal.
![](./files/goplay_screen1.png)
## 功能特性

*   **分屏界面**: 左侧编辑代码，右侧查看输出。
*   **即时执行**: 一键运行 (`go run`) 或编译检查 (`go build`) 代码。
*   **文件管理**: 轻松打开、保存和创建新文件。
*   **模板支持**: 快速加载预定义的 `.template` 文件或使用内置默认模板。
*   **构建功能**: 将代码构建为独立的可执行文件 (`Ctrl+b`)。
*   **系统剪贴板支持**: 支持与系统剪贴板的复制/粘贴交互。
*   **编辑器工具**:
    *   使用 `go fmt` 自动格式化。
    *   基于快照的撤销/重做。
    *   实时光标行号追踪。
    *   修改状态指示。
*   **内置帮助**: 随时通过 `F1` 查看快捷键参考。

## 安装

### 前置要求
*   系统已安装 [Go](https://go.dev/dl/)。

### 源码构建

1.  克隆仓库:
    ```bash
    git clone https://github.com/yourusername/goplay.git
    cd goplay
    ```

2.  构建应用:
    ```bash
    go build -o goplay
    ```

3.  运行:
    ```bash
    ./goplay
    ```

## 使用指南

### 命令行参数

你可以通过以下方式启动 GoPlay：

```bash
# 启动空编辑器
./goplay

# 直接打开指定文件
./goplay main.go

# 打开指定目录（并将其设置为工作目录）
./goplay /path/to/project

# 打开当前目录
./goplay .
```

### 快捷键列表

| 快捷键 | 动作 | 说明 |
| :--- | :--- | :--- |
| **F1** | 帮助 | 显示快捷键帮助窗口。 |
| **F2** | 缩小窗口 | 减小左侧编辑器窗口宽度。 |
| **F3** | 扩大窗口 | 增加左侧编辑器窗口宽度。 |
| **F4** / **Ctrl+T** | 模板 | 加载模板 (自定义或默认)。 |
| **F5** / **Ctrl+R** | 运行 | 运行当前代码并显示输出。 |
| **F6** / **Ctrl+K** | 编译检查 | 编译代码以检查错误/成功状态。 |
| **F7** / **Ctrl+B** | 构建应用 | 将代码构建为独立的可执行文件。 |
| **F9** / **Ctrl+P** | 设置 | 打开设置 (设置工作目录)。 |
| **F10** | 关于 | 显示关于信息弹窗。 |
| **F12** / **Ctrl+Q** | 退出 | 退出应用程序。 |
| **Ctrl+F** | 格式化 | 使用 `go fmt` 格式化代码。 |
| **Ctrl+S** | 保存 | 保存当前代码到文件。 |
| **Ctrl+O** | 打开 | 打开现有文件。 |
| **Ctrl+N** | 新建 | 清空编辑器以创建新文件。 |
| **Ctrl+L** | 清空 | 清空编辑器内容。 |
| **Ctrl+Z** | 撤销 | 撤销上一步操作。 |
| **Ctrl+Y** | 重做 | 重做上一步操作。 |
| **Ctrl+G** | 跳转行 | 跳转到指定的行号。 |
| **Ctrl+C** | 复制 | 将选中的文本复制到系统剪贴板。 |
| **Esc** | 关闭 | 关闭弹窗或帮助菜单。 |

### 详细功能说明

#### 1. 模板功能 (F4 / Ctrl+T)
按 `F4` 或 `Ctrl+T` 加载模板。
- **默认模板**: 如果当前工作目录下没有 `.template` 文件，程序将加载内置的 "Hello Goplay" 示例代码。
- **自定义模板**: 在当前目录下创建一个名为 `.template` 的文件，即可使用你自己的样板代码。

#### 2. 编译与构建 (Compile & Build)
- **编译检查 (`F6` / `Ctrl+K`)**: 仅执行编译过程以检查语法错误，不生成文件。
- **构建应用 (`F7` / `Ctrl+B`)**: 会提示你输入输出文件名，然后在当前目录（或指定的工作目录）生成 `.exe` (Windows) 或二进制文件。

#### 3. 窗口调整
你可以调整左右分屏的比例：
- 按 **F2** 向左收缩编辑器，增大输出窗口。
- 按 **F3** 向右扩展编辑器，减小输出窗口。

#### 4. 复制与粘贴 (Copy & Paste)
- **复制**: 在编辑器中选中文本，按 **`Ctrl+C`** 将其复制到系统剪贴板。
  - *注意*: 只有当编辑器拥有焦点时 `Ctrl+C` 才会执行复制，否则它不会退出程序。
- **粘贴**: 使用系统标准的粘贴快捷键（如 **`Ctrl+V`** 或终端右键粘贴）将代码粘贴到编辑器中。

#### 5. 工作目录 (Working Directory)
使用 **F9** 或 **Ctrl+P** 设置工作目录。相对路径的文件操作都将相对于此目录进行。

#### 6. 退出程序
使用 **`F12`** 或 **`Ctrl+Q`** 安全退出应用程序。

## 许可证

[MIT](LICENSE)
