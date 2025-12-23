## 命令行go playground
这是一个用go语言实现命令行实现得类似 `go playground` 得程序, 有以下功能
- 主界面分成左右两部分, 左边可以输入代码, 右边显示运行结果
- 为了简化实现, 代码编辑器只需要最简单得文本输入功能, 无需高亮等语法提示
- 支持快捷键, 快速实现一些功能
    - Ctrl + r, 运行代码, 并在右边显示运行结果
    - Ctrl + k, 编译当前代码, 并在右边显示结果, 如果有错误得话
    - Ctrl + s, 保存代码, 如果当前没有文件名, 提示输入文件名
    - Ctrl + l, 清空代码
    - Ctrl + shift + q, 退出
    - Ctrl + f, 使用go fmt 格式化当前代码
    - Ctrl + n, 创建新文件
    - Ctrl + o, 打开文件, 提示输入文件名
    - Ctrl + z, 撤销
    - Ctrl + y, 重做
    - Ctrl + t, 加载 .template 模板到当前编辑器


new feature:
    - Ctrl + g, 提示输入行数, 光标跳转到指定行
    - Ctrl + p, 显示设置窗口, 目前可以显示并设置 `当前工作目录(CWD)`, 这样输入的文件名如果是相对路径,都是相对于该目录的, 软件启动时, 将启动的CWD 自动设置为当前工作目录, 如果当前文件名是绝对路径, 则不使用CWD,直接使用路径打开
    - Ctrl + b, build 当前文件为可执行文件, 并提示输入输出文件名, 默认为当前文件名, 考虑平台差异, windows 输出 app.exe, linux, macos 输出 app, 如果当前文件没保存, 提示先保存文件


new feature:
    - Alt + [, 将窗口分割线左移, 让左边窗口变小, 右边变大
    - Alt + ], 将窗口分割线右移, 让左边窗口变大, 右边变小
    - Alt + h, 显示帮助窗口, 显示所有快捷键

new feature:
    - 添加一个程序默认模板, 当按`Ctrl+t`, 检查当前cwd 是否有`.template`文件, 如果有, 则加载该文件到编辑器, 如果没有, 则使用默认模板
    ```
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello Goplay")
}
    ```


goplay改进 - v0.3.0

1. 增加参数处理
- `goplay file.go` 如果是文件, 则直接打开该文件
- `goplay d:\dev\go` 如果是目录, 则打开该目录, 并将当前工作目录设置为该目录
- `goplay .` 则打开当前目录

2. 改进快捷键, 增加更多F键, 取消所有alt键, 具体如下
- `Alt + [` Or F2 -> F2
- `Alt + ]` Or F3 -> F3
- `Alt + h` Or F1 -> F1
- `Ctrl + t` -> `Ctrl + t` Or F4
- `Ctrl + r` -> `Ctrl + r` Or F5
- `Ctrl + k` -> `Ctrl + k` Or F6
- `Ctrl + b` -> `Ctrl + b` Or F7
- `Ctrl + q` -> `Ctrl + q` Or F12

3. 增加一些简单自动格式化功能
- 按回车后保持和上一行一致缩进
- 请建议其他常用命令行格式化功能

goplay v0.3.0 fix
bugs
- main screen bottom right - should be `Press F1 to see all keys`
- auto format not work, when press enter, cursor go to top of window

enhance
- add F9 -> `Ctrl + p`
- add a pop up, use F10 with `about goplay`, with following
```
goplay v0.3.0
command line go playground

author: your name
license: MIT

```
- 