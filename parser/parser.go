package parser

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

type MethodCall struct {
	Class  string
	Method string
}

func ParseMethodCall(content string, row uint, column uint) *MethodCall {
	code := []byte(content)

	parser := tree_sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_python.Language()))

	tree := parser.Parse(code, nil)
	defer tree.Close()

	root := tree.RootNode()

	point := tree_sitter.Point{Row: row, Column: column}
	node := root.NamedDescendantForPointRange(point, point)

	if node == nil {
		return nil
	}

	gname := node.GrammarName()
	parent := node.Parent()

	if gname == "keyword_argument" {
		node = parent
	}

	if gname != "argument_list" && gname != "call" && parent != nil && parent.GrammarName() == "argument_list" {
		node = parent.Parent()
		parent = node.Parent()
	}

	if node.GrammarName() == "argument_list" {
		node = parent
	}

	if node.GrammarName() != "call" {
		return nil
	}

	fn := node.ChildByFieldName("function")
	if fn == nil {
		return nil
	}

	attr := fn.ChildByFieldName("attribute")
	if attr == nil || attr.GrammarName() != "identifier" {
		return nil
	}

	// TODO: also find what package the class comes from
	// TODO: make this more robust (check that it is a QuerySet/Manager, handle substitution)

	call := MethodCall{Class: getClassName(code, node), Method: nodeContent(code, attr)}
	if call.Class == "" || call.Method == "" {
		return nil
	}

	return &call
}

func nodeContent(src []byte, node *tree_sitter.Node) string {
	return string(src[node.StartByte():node.EndByte()])
}

func getClassName(src []byte, node *tree_sitter.Node) string {
	if node == nil {
		return ""
	}

	switch node.GrammarName() {
	case "identifier":
		return nodeContent(src, node)
	case "attribute":
		return getClassName(src, node.ChildByFieldName("object"))
	case "call":
		return getClassName(src, node.ChildByFieldName("function"))
	case "subscript":
		return getClassName(src, node.ChildByFieldName("value"))
	}

	return nodeContent(src, node)
}
