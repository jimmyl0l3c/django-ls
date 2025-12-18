package analyzer

import (
	_ "embed"
	"errors"
	"os"
	"path"
)

//go:embed model_analyzer.py
var modelAnalyzerPy []byte

func getPythonBin() (string, error) {
	// TODO: find a better solution

	venv := os.Getenv("VIRTUAL_ENV")
	if venv == "" {
		return "", errors.New("no python venv active")
	}

	return path.Join(venv, "bin", "python3"), nil
}
