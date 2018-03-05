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
	flgKeyword = flag.String("keys", "todo,fix", "Change keyword")

	fset *token.FileSet

	keywords []string
	dir      string
)

type comment struct {
	c        *ast.Comment
	b        *bufio.Reader
	line     int
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
	for {
		line, _, err := c.b.ReadLine()
		if err != nil {
			break
		}
		sComment := string(bytes.TrimSpace(line))
		if len(sComment) < 4 {
			c.line++
			continue
		}
		for _, kw := range keywords {
			if strings.EqualFold(kw, string(sComment[0:len(kw)])) {
				fmt.Fprintf(w, "%s:%d:%d:%s\n", c.path(), c.pos().Line+c.line, c.pos().Column, sComment)
			}
		}

		c.line++
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
	keywords = strings.Split(*flgKeyword, ",")
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
