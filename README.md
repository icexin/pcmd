# Paralleled command execution tool

## Overview

假设你有一个url列表文件，每行一个url，你想并发使用curl命令把列表里的地址抓下来保存到文件里面，pcmd或许能帮助你。

```
cat url.txt | pcmd -c 100 'echo "downloading {{1}}" && curl {{1}} > {{1}}.html'
```

上面命令将以100的并发来下载文件。


## Install

`go get github.com/icexin/pcmd`

## Usage

pcmd从标准输入里面把每一行按空格，tab分割，第一个字段是`{{1}}`，第二个是`{{2}}`，以此类推，总共可以支持9个占位符。

`{{i}}`表示数据的行号

-c 指定并发数
