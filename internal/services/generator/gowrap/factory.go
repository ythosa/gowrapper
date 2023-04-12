package gowrap

import (
	"github.com/ythosa/gowrapper/internal/model"
	"github.com/ythosa/gowrapper/internal/services/finder"
	"github.com/ythosa/gowrapper/internal/services/generator"
	"github.com/ythosa/gowrapper/internal/services/options_builder"
)

type Factory struct{}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Construct(opts model.Options) generator.Generator {
	builder := options_builder.NewOptionsBuilder()
	builder.SetOutputFolder(opts.OutputFolder)
	builder.SetFilePostfix(opts.FilePostfix)
	builder.SetTemplate(opts.TemplatePath)

	g := NewGenerator(finder.NewInterfacesFinder(), builder)
	g.SetExcludedDirs(opts.ExcludedDirs)
	g.SetExcludedFiles(opts.ExcludedFiles)
	g.SetInitialDirectory(opts.InitialDirectory)

	return g
}
