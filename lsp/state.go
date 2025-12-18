package lsp

import (
	"log/slog"

	"github.com/jimmyl0l3c/django-ls/analyzer"
	"github.com/jimmyl0l3c/django-ls/safemap"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var documentStates = safemap.New[*DocumentState]()
var djangoWorkspace *analyzer.DjangoWorkspace

type DocumentState struct {
	Content     string
	DocumentURI protocol.DocumentUri
}

func getDocument(key string) (*DocumentState, bool) {
	doc, ok := documentStates.Load(key)
	if !ok {
		return nil, false
	}

	return doc, ok
}

func Init(params *protocol.InitializeParams) {
	if len(params.WorkspaceFolders) == 0 {
		slog.Error("No workspace folders specified.")
		return
	}

	for _, wf := range params.WorkspaceFolders {
		if dw, err := analyzer.AnalyzeWorkspace(wf); err == nil && dw != nil {
			djangoWorkspace = dw
			return
		}
	}

	slog.Error("No valid workspace specified.")
}
