package analyzer

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"regexp"
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func AnalyzeWorkspace(wf protocol.WorkspaceFolder) (*DjangoWorkspace, error) {
	wpath, ok := strings.CutPrefix(wf.URI, "file://")
	if !ok {
		slog.Error("Invalid or unsupported workspace URI.", "uri", wf.URI)
		return nil, fmt.Errorf("invalid or unsupported workspace URI: %s", wf.URI)
	}

	workspace := findDjangoWorkspace(wpath)
	if workspace == nil {
		return nil, errors.New("no valid workspace found")
	}

	go workspace.runAnalyzer()

	return workspace, nil
}

func findDjangoWorkspace(root string) *DjangoWorkspace {
	dfs := os.DirFS(root)

	matches, err := fs.Glob(dfs, "manage.py")
	if err != nil {
		slog.Error("Settings search failed.", "error", err)
		return nil
	}

	re := regexp.MustCompile(`os\.environ\.setdefault\("DJANGO_SETTINGS_MODULE", "([^"]+)"\)`)

	for _, m := range matches {
		f, err := dfs.Open(m)
		if err != nil {
			slog.Error("Could not open 'manage.py'.", "error", err)
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			sm := re.FindStringSubmatch(scanner.Text())
			if sm == nil {
				continue
			}

			if len(sm) == 2 {
				return newWorkspace(path.Join(root, path.Dir(m)), sm[1])
			}
		}

		if err := scanner.Err(); err != nil {
			slog.Error("Reading 'manage.py' failed.", "error", err)
			return nil
		}
	}

	return nil
}
