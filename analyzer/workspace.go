package analyzer

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os/exec"

	"github.com/jimmyl0l3c/django-ls/parser"
	"github.com/jimmyl0l3c/django-ls/safemap"
)

type DjangoWorkspace struct {
	RootPath       string
	SettingsModule string

	models *safemap.SafeMap[*DjangoModel]
}

func newWorkspace(root string, settings string) *DjangoWorkspace {
	return &DjangoWorkspace{
		RootPath:       root,
		SettingsModule: settings,
		models:         safemap.New[*DjangoModel](),
	}
}

func (dw *DjangoWorkspace) GetModel(name string) *DjangoModel {
	if dw == nil {
		return nil
	}

	model, ok := dw.models.Load(name)
	if !ok {
		return nil
	}

	return model
}

func (dw *DjangoWorkspace) GetLookups(call *parser.MethodCall) []FieldLookup {
	m := dw.GetModel(call.Class)
	if m == nil {
		return nil
	}

	// TODO: check method

	return m.GetLookups()
}

func (gw *DjangoWorkspace) runAnalyzer() error {
	py, err := getPythonBin()
	if err != nil {
		slog.Error("Could not get python bin.", "error", err)
		return err
	}

	cmd := exec.Command(py, "-", gw.RootPath, "-s", gw.SettingsModule)

	cmd.Stdin = bytes.NewBuffer(modelAnalyzerPy)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("Could not pipe stdout.", "error", err)
		return err
	}
	defer stdout.Close()

	slog.Debug("Starting analyzer.", "interpreter", py, "root", gw.RootPath, "settings", gw.SettingsModule)

	if err := cmd.Start(); err != nil {
		slog.Error("Could not start cmd.", "error", err)
		return err
	}

	var jsonData []DjangoModel
	if err := json.NewDecoder(stdout).Decode(&jsonData); err != nil {
		slog.Error("Could not decode data.", "error", err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		slog.Error("Cmd wait failed.", "error", err)
		return err
	}

	for _, v := range jsonData {
		gw.models.Store(v.Name, &v)
		go v.syncLookups()
	}

	slog.Debug("Models loaded.")

	return nil
}
