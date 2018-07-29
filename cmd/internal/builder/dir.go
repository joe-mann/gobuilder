package builder

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/joe-mann/gobuilder/cmd/internal/executor"
	"github.com/sirupsen/logrus"
)

const (
	BuildDir = ".build"
)

type Dir struct {
	path    string
	relpath string
	files   []string
	stdout  bytes.Buffer
}

func NewDir(path, relpath string) (*Dir, error) {
	d := &Dir{path: path, relpath: relpath}
	var err error
	d.files, err = filepath.Glob(filepath.Join(path, "*_build.go"))
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Dir) Build() bool {
	if len(d.files) == 0 {
		return true
	}

	e := executor.New(d.relpath)
	for _, path := range d.files {
		f, err := os.Open(path)
		if err != nil {
			d.stdout.WriteString(err.Error())
			return false
		}
		if err := e.AddBuildFile(f, path); err != nil {
			d.stdout.WriteString(err.Error())
			return false
		}
	}

	testDir := filepath.Join(d.path, ".build")
	if err := os.RemoveAll(testDir); err != nil {
		d.stdout.WriteString(err.Error())
		return false
	}

	if err := os.MkdirAll(testDir, os.ModePerm); err != nil {
		d.stdout.WriteString(err.Error())
		return false
	}

	f, err := os.Create(filepath.Join(testDir, "build_main.go"))
	if err != nil {
		d.stdout.WriteString(err.Error())
		return false
	}

	e.WriteBuildMain(f)
	return true
}

func (d *Dir) NumBuildFiles() uint {
	return uint(len(d.files))
}

func (d *Dir) HasBuildFiles() bool {
	return len(d.files) != 0
}

func (d *Dir) String() string {
	return d.relpath
}

func (d *Dir) Print(logger *logrus.Logger) {
	if d.stdout.Len() != 0 {
		logger.Print(d.stdout.String())
	}
}
