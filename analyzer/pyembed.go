package analyzer

import (
	_ "embed"
	"errors"
	"os"
	"path"
)

const analyzeScriptName = "analyze.py"

//go:embed model_analyzer.py
var analyzePy []byte

var analyzePyPath string

func saveAnalyzeScript() error {
	// TODO: use better path, cache the script between runs
	path := path.Join("/tmp", analyzeScriptName)

	err := os.WriteFile(path, analyzePy, 0666)
	if err != nil {
		return err
	}

	analyzePyPath = path
	return nil
}

func getAnalyzeScript() (string, error) {
	if analyzePyPath != "" {
		return analyzePyPath, nil
	}

	err := saveAnalyzeScript()
	if err != nil {
		return "", err
	}

	return analyzePyPath, nil
}

func getPythonBin() (string, error) {
	// TODO: find a better solution

	venv := os.Getenv("VIRTUAL_ENV")
	if venv == "" {
		return "", errors.New("no python venv active")
	}

	return path.Join(venv, "bin", "python3"), nil
}
