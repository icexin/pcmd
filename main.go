package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	nworker = flag.Int("c", 50, "concurrent")

	seq int32
)

type lineWriter struct {
	buf bytes.Buffer
}

func (w *lineWriter) Write(b []byte) (int, error) {
	for _, c := range b {
		if c == '\n' {
			w.Flush()
		} else {
			w.buf.WriteByte(c)
		}
	}
	return len(b), nil
}

func (w *lineWriter) Flush() {
	if w.buf.Len() != 0 {
		fmt.Println(w.buf.String())
		w.buf.Reset()
	}
}

type item struct {
	data string
	idx  int
}

func render(cmd string, x item) string {
	fs := strings.Fields(x.data)
	for i := 1; i < 10; i++ {
		old := fmt.Sprintf("{{%d}}", i)
		new := ""
		if i <= len(fs) {
			new = fs[i-1]
		}
		cmd = strings.Replace(cmd, old, new, -1)
	}
	cmd = strings.Replace(cmd, "{{i}}", fmt.Sprint(x.idx), -1)
	return cmd
}

func worker(cmd string, ch chan item, wg *sync.WaitGroup) {
	defer wg.Done()
	stdout := new(lineWriter)
	stderr := new(lineWriter)
	for item := range ch {
		cmdstr := render(cmd, item)
		c := exec.Command("bash", "-c", cmdstr)
		c.Stdout = stdout
		c.Stderr = stderr
		c.Run()
		stdout.Flush()
		stderr.Flush()
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s $cmd", os.Args[0])
		return
	}

	wg := new(sync.WaitGroup)
	ch := make(chan item, *nworker)
	for i := 0; i < *nworker; i++ {
		wg.Add(1)
		go worker(flag.Arg(0), ch, wg)
	}

	idx := 0
	r := bufio.NewScanner(os.Stdin)
	for r.Scan() {
		idx++
		line := r.Text()
		ch <- item{line, idx}
	}
	close(ch)
	wg.Wait()
}
