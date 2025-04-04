package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	_ "github.com/tliron/commonlog/simple"
)

func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	var completionItems []protocol.CompletionItem

	value := "Bar"

	completionItems = append(completionItems, protocol.CompletionItem{
		Label:      "foo",
		Detail:     &value,
		InsertText: &value,
	})

	return completionItems, nil
}
