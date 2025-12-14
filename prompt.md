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