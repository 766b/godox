package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	flgKeyword = flag.String("keys", "todo,bug,fix", "Change keywords")

	fset *token.FileSet

	keywords [][]byte
	dir      string
)

type comment struct {
	c        *ast.Comment
	b        *bufio.Reader
	g        *godox
	tokenPos token.Position
}

func (c comment) pos() *token.Position {
	if !c.tokenPos.IsValid() {
		c.tokenPos = fset.Position(c.c.Pos())
	}
	return &c.tokenPos
}

func (c comment) path() string {
	path := filepath.Join(c.pos().Filename)
	return path
}

func (c comment) printTodoLines(w io.Writer) {
	var lineNum int
	for {
		line, _, err := c.b.ReadLine()
		if err != nil {
			break
		}
		sComment := bytes.TrimSpace(line)
		if len(sComment) < 4 {
			lineNum++
			continue
		}
		for _, kw := range keywords {
			if bytes.EqualFold(kw, sComment[0:len(kw)]) {
				fmt.Fprintf(w, "%s:%d:%d:%s\n", c.path(), c.pos().Line+lineNum, c.pos().Column, sComment)
				break
			}
		}
		lineNum++
	}
	return
}

func (g *godox) newComment(c *ast.Comment) comment {
	cT := c.Text
	switch cT[1] {
	case '/':
		cT = cT[2:]
		if len(cT) > 0 && cT[0] == ' ' {
			cT = cT[1:]
		}
	case '*':
		cT = cT[2 : len(cT)-2]
	}

	return comment{
		c: c,
		b: bufio.NewReader(bytes.NewBufferString(cT)),
		g: g,
	}
}

func (g *godox) parse() {
	for _, k := range strings.Split(*flgKeyword, ",") {
		keywords = append(keywords, []byte(k))
	}

	fset = token.NewFileSet()
	f, err := parser.ParseDir(fset, g.dir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for _, pkg := range f {
		for _, file := range pkg.Files {
			for _, c := range file.Comments {
				for _, ci := range c.List {
					g.newComment(ci).printTodoLines(g.out)
				}
			}
		}
	}
}

type godox struct {
	dir string
	out io.Writer
}

func new(dir string, out io.Writer) *godox {
	return &godox{
		dir,
		out,
	}
}

func main() {
	flag.Parse()
	if len(os.Args) == 1 {
		new(".", os.Stdout).parse()
		return
	}

	for _, arg := range os.Args {
		fs, err := os.Stat(arg)
		if err != nil {
			panic(err)
		}
		if fs.IsDir() {
			new(arg, os.Stdout).parse()
		}
	}
}
