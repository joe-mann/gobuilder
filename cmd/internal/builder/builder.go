package builder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joe-mann/gobuilder/cmd/internal/view"
	"github.com/sirupsen/logrus"
)

type Builder struct {
	config Config
	logger *logrus.Logger
}

func New(config Config, logger *logrus.Logger) *Builder {
	return &Builder{config: config, logger: logger}
}

func (b *Builder) Build() bool {
	dirs := b.dirs()
	for i := len(dirs) - 1; i >= 0; i-- {
		if !b.build(dirs[i]) {
			return false
		}
	}
	return true
}

func (b *Builder) build(d Dir) bool {
	if !d.HasBuildFiles() {
		view.Verbose(b.logger, "?      %s      [%s]", d.String(), buildFilesString(0))
		return true
	}
	if !d.Build() {
		view.Always(b.logger, "✗      %s      [%s]", d.relpath, buildFilesString(len(d.files)))
		d.Print(b.logger)
		return false
	}
	view.Always(b.logger, "✓      %s      [%s]", d.String(), buildFilesString(len(d.files)))
	d.Print(b.logger)
	return true
}

func (b *Builder) dirs() []Dir {
	dirs := []Dir{}
	filepath.Walk(b.config.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			return nil
		}
		if b.shouldSkipDir(info.Name()) {
			return filepath.SkipDir
		}

		dir, err := NewDir(path, b.relpath(path))
		if err != nil {
			log.Fatal(err)
		}
		dirs = append(dirs, *dir)
		return nil
	})
	return dirs
}

func (b *Builder) shouldSkipDir(name string) bool {
	for _, dir := range b.config.excludedDirs {
		if name == dir {
			return true
		}
	}
	return false
}

func (b *Builder) relpath(abs string) string {
	srcDir := filepath.Join(b.config.gopath, "src/")
	if !strings.HasPrefix(abs, srcDir) {
		b.logger.Fatalf("'%s' is not in GOPATH")
	}
	relpath := strings.TrimPrefix(abs, srcDir)
	return strings.TrimPrefix(relpath, "/")
}

func buildFilesString(n int) string {
	switch n {
	case 0:
		return "no build files"
	case 1:
		return "1 build file"
	default:
		return fmt.Sprintf("%d build files", n)
	}
}
