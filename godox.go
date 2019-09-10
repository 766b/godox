package godox

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
)

var (
	defaultKeywords = []string{"TODO", "BUG", "FIXME"}
)

// Message contains a message and position
type Message struct {
	Pos     token.Position
	Message string
}

func getMessages(c *ast.Comment, fset *token.FileSet, keywords []string) []Message {
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
	var comments []Message

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
			if bytes.EqualFold([]byte(kw), sComment[0:len(kw)]) {
				pos := fset.Position(c.Pos())
				comments = append(comments, Message{
					Pos:     pos,
					Message: fmt.Sprintf("%s:%d:%d:%s", filepath.Join(pos.Filename), pos.Line+lineNum, pos.Column, sComment),
				})
				break
			}
		}
	}
	return comments
}

func Run(file *ast.File, fset *token.FileSet, keywords ...string) []Message {
	if keywords == nil || len(keywords) == 0 {
		keywords = defaultKeywords
	}
	var messages []Message
	for _, c := range file.Comments {
		for _, ci := range c.List {
			messages = append(messages, getMessages(ci, fset, keywords)...)
		}
	}
	return messages
}
