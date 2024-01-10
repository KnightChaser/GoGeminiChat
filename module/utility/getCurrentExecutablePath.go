package utility

import (
	"os"
	"path/filepath"
)

func GetCurrentExecutablePath() string {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	executablePath := filepath.Dir(executable)

	return executablePath
}
