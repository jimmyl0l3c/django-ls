package lsp

import (
	"log/slog"

	"github.com/jimmyl0l3c/django-ls/parser"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	_ "github.com/tliron/commonlog/simple"
)

func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	slog.Debug("Completion", "params", params)

	doc, ok := getDocument(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}

	call := parser.ParseMethodCall(doc.Content, uint(params.Position.Line), uint(params.Position.Character))
	if call == nil {
		return nil, nil
	}

	slog.Debug("MethodCall found", "call", call)

	lookups := djangoWorkspace.GetLookups(call)
	if lookups == nil {
		return nil, nil
	}

	var completionItems []protocol.CompletionItem
	kind := protocol.CompletionItemKindField

	for _, field := range lookups {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      field.Name,
			Detail:     &field.TypeName,
			InsertText: &field.Name,
			Kind:       &kind,
		})
	}

	return completionItems, nil
}
