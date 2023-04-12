package gowrap

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"

	gowrap "github.com/hexdigest/gowrap/generator"
	"github.com/ythosa/gowrapper/internal/logging"
	"github.com/ythosa/gowrapper/internal/model"
)

type Finder interface {
	Find(initial string, excludedDirs []string, excludedFiles []string) model.InterfacesByFileMap
}

type OptionsBuilder interface {
	Build(model.InterfacesByFileMap) []gowrap.Options
}

type Generator struct {
	logger logging.Logger

	finder  Finder
	builder OptionsBuilder

	initialDirectory string
	excludedDirs     []string
	excludedFiles    []string
}

func NewGenerator(finder Finder, builder OptionsBuilder) *Generator {
	return &Generator{
		logger:  logging.New("generator"),
		finder:  finder,
		builder: builder,
	}
}

func (g *Generator) SetInitialDirectory(initialDirectory string) {
	if initialDirectory == "" {
		g.logger.Fatal("initial directory is empty")
	}

	g.initialDirectory = initialDirectory
}

func (g *Generator) SetExcludedDirs(excludedDirs []string) {
	g.excludedDirs = excludedDirs
}

func (g *Generator) SetExcludedFiles(excludedFiles []string) {
	g.excludedFiles = excludedFiles
}

func (g *Generator) Generate() {
	interfacesByFile := g.finder.Find(g.initialDirectory, g.excludedDirs, g.excludedFiles)
	options := g.builder.Build(interfacesByFile)

	var wg sync.WaitGroup

	for _, option := range options {
		wg.Add(1)
		go func(opt gowrap.Options) {
			g.generateFromOption(opt)
			wg.Done()
		}(option)
	}

	wg.Wait()
}

func (g *Generator) generateFromOption(option gowrap.Options) {
	wrapper, err := gowrap.NewGenerator(option)
	if err != nil {
		g.logger.Errorf("failed to create generator for option: %v : %s", option, err)

		return
	}

	buf := bytes.NewBuffer([]byte{})
	if err = wrapper.Generate(buf); err != nil {
		g.logger.Errorf("failed to generate wrap: %s", err)

		return
	}

	if err = os.MkdirAll(filepath.Dir(option.OutputFile), os.ModePerm); err != nil {
		g.logger.Errorf("failed to make dirs for output: %s", err)

		return
	}

	if err = os.WriteFile(option.OutputFile, buf.Bytes(), 0664); err != nil {
		g.logger.Errorf("failed to write file: %s", err)

		return
	}

	g.logger.Infof("generated: %s", option.OutputFile)
}
