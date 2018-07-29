package cmd

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/joe-mann/gobuilder/cmd/internal/view"

	"github.com/joe-mann/gobuilder/cmd/internal/builder"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		build()
	},
}

type formatter struct{}

func (f *formatter) Format(e *logrus.Entry) ([]byte, error) {
	return []byte(e.Message + "\n"), nil
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func build() bool {
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	cwd := cwd()
	logger.Debugf("cwd = '%s'", cwd)
	gopath := gopath()
	logger.Debugf("GOPATH = '%s'\n", gopath)

	b := builder.New(config(cwd, gopath), &logger)
	if !b.Build() {
		view.Always(&logger, "BUILD FAILED!")
		return false
	}

	view.Always(&logger, "BUILD SUCCESSFUL")
	return true
}

func config(cwd, gopath string) builder.Config {
	c := builder.NewConfig(cwd, gopath)
	if len(excludeDirs) != 0 {
		c.SetExcludedDirs(excludeDirs)
	}
	return *c
}

func cwd() string {
	path, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		logger.Fatal(err)
	}
	return path
}

func gopath() string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		gopath, err := filepath.Abs(gopath)
		if err != nil {
			logger.Fatal(err)
		}
		return gopath
	}

	logger.Debugln("GOPATH not set, using default")
	usr, err := user.Current()
	if err != nil {
		logger.Fatal(err)
	}
	gopath, err = filepath.Abs(usr.HomeDir)
	if err != nil {
		logger.Fatal(err)
	}
	return gopath
}
