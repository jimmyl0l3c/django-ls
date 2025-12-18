package lsp

import (
	"log/slog"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// protocol.TextDocumentDidOpenFunc signature
func TextDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	slog.Debug("DocumentOpen", "params", params)

	if params.TextDocument.LanguageID != "python" {
		slog.Debug("Skipping non-python file.", "languageID", params.TextDocument.LanguageID)
		return nil
	}

	documentStates.Store(params.TextDocument.URI, &DocumentState{
		DocumentURI: params.TextDocument.URI,
		Content:     params.TextDocument.Text,
	})

	return nil
}

// protocol.TextDocumentDidChangeFunc signature
func TextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	slog.Debug("DocumentChange", "params", params)

	doc, ok := getDocument(params.TextDocument.URI)
	if !ok {
		return nil
	}

	for _, change := range params.ContentChanges {
		if e, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
			start, end := e.Range.IndexesIn(doc.Content)
			content := doc.Content[:start] + e.Text + doc.Content[end:]
			doc.Content = content
		} else if e, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			doc.Content = e.Text
		}
	}

	return nil
}

// protocol.TextDocumentDidCloseFunc signature
func TextDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	slog.Debug("DocumentClose", "params", params)

	documentStates.Delete(params.TextDocument.URI)

	return nil
}
