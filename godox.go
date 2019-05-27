package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var (
	flgKeyword = flag.String("keys", "todo,bug,fix", "Change keywords")
	fset       *token.FileSet
	keywords   [][]byte
)

func newComment(c *ast.Comment) []string {
	commentText := c.Text
	switch commentText[1] {
	case '/':
		commentText = commentText[2:]
		if len(commentText) > 0 && commentText[0] == ' ' {
			commentText = commentText[1:]
		}
	case '*':
		commentText = commentText[2 : len(commentText)-2]
	}

	b := bufio.NewReader(bytes.NewBufferString(commentText))

	var comments []string
	for lineNum := 0; ; lineNum++ {
		line, _, err := b.ReadLine()
		if err != nil {
			break
		}
		sComment := bytes.TrimSpace(line)
		if len(sComment) < 4 {
			continue
		}
		for _, kw := range keywords {
			if bytes.EqualFold(kw, sComment[0:len(kw)]) {
				pos := fset.Position(c.Pos())
				comments = append(comments, fmt.Sprintf("%s:%d:%d:%s", filepath.Join(pos.Filename), pos.Line+lineNum, pos.Column, sComment))
				break
			}
		}
	}

	return comments
}

func godox(rootPath string, includeTests bool) ([]string, error) {
	for _, k := range strings.Split(*flgKeyword, ",") {
		keywords = append(keywords, []byte(k))
	}

	const recursiveSuffix = string(filepath.Separator) + "..."
	recursive := false
	if strings.HasSuffix(rootPath, recursiveSuffix) {
		recursive = true
		rootPath = rootPath[:len(rootPath)-len(recursiveSuffix)]
	}

	var messages []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if !recursive && path != rootPath {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if !includeTests && strings.HasSuffix(path, "_test.go") {
			return nil
		}

		fset = token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, c := range f.Comments {
			for _, ci := range c.List {
				messages = append(messages, newComment(ci)...)
			}
		}
		return nil
	})

	return messages, err
}
