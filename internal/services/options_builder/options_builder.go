package options_builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	gowrap "github.com/hexdigest/gowrap/generator"
	"github.com/hexdigest/gowrap/pkg"
	"github.com/samber/lo"
	"github.com/schollz/progressbar/v3"
	"github.com/ythosa/gowrapper/internal/logging"
	"github.com/ythosa/gowrapper/internal/model"
)

const headerTemplate = `// Code generated by gowrap. DO NOT EDIT.
// template: {{.Options.HeaderVars.Template}}
// gowrap: https://github.com/hexdigest/gowrap

package {{.Package.Name}}

{{if (not .Options.HeaderVars.DisableGoGenerate)}}
//{{"go:generate"}} gowrap gen -p {{.SourcePackage.PkgPath}} -i {{.Options.InterfaceName}} -t {{.Options.HeaderVars.Template}} -o {{.Options.HeaderVars.OutputFileName}}{{.Options.HeaderVars.VarsArgs}} -l "{{.Options.LocalPrefix}}"
{{end}}

`

type OptionsBuilder struct {
	logger logging.Logger

	template     []byte
	templatePath string
	outputFolder string
	filePostfix  string
}

func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{logger: logging.New("options_builder")}
}

func (builder *OptionsBuilder) SetTemplate(templatePath string) {
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		builder.logger.Fatalf("template file(%s) is not exist", templatePath)
	}

	body, err := os.ReadFile(templatePath)
	if err != nil {
		builder.logger.Fatalf("failed to load template file: %s", err)
	}

	builder.template = body
	builder.templatePath = templatePath
}

func (builder *OptionsBuilder) SetOutputFolder(outputFolder string) {
	builder.outputFolder = outputFolder
}

func (builder *OptionsBuilder) SetFilePostfix(filePostfix string) {
	builder.filePostfix = filePostfix
}

func (builder *OptionsBuilder) Build(interfacesByFile model.InterfacesByFileMap) []gowrap.Options {
	var options = make([]gowrap.Options, 0, len(interfacesByFile))

	builder.logger.Info("Processing files...")
	progressBar := progressbar.Default(int64(len(interfacesByFile)))

	for file, interfaces := range interfacesByFile {
		currentOpts := lo.Map(interfaces, func(i string, _ int) gowrap.Options { return builder.buildOptions(file, i) })
		options = append(options, currentOpts...)
		_ = progressBar.Add(1)
	}

	return options
}

func (builder *OptionsBuilder) buildOptions(file string, infs string) gowrap.Options {
	dir := filepath.Dir(file)

	sourcePackage, err := pkg.Load(dir)
	if err != nil {
		builder.logger.Errorf("failed to load source package: %s", err)

		return gowrap.Options{}
	}

	outputFile := builder.buildOutput(dir, infs)
	_, outputFileName := filepath.Split(outputFile)

	return gowrap.Options{
		InterfaceName:  infs,
		SourcePackage:  sourcePackage.PkgPath,
		OutputFile:     outputFile,
		HeaderTemplate: headerTemplate,
		BodyTemplate:   string(builder.template),
		HeaderVars: map[string]interface{}{
			"Template":       builder.templatePath,
			"OutputFileName": outputFileName,
			"VarsArgs":       "",
		},
	}
}

func (builder *OptionsBuilder) buildOutput(dir, infs string) string {
	if builder.outputFolder != "" {
		dir += "/" + builder.outputFolder
	}

	return fmt.Sprintf("%s/%s_%s.go", dir, builder.buildFileName(infs), builder.filePostfix)
}

func (builder *OptionsBuilder) buildFileName(interfaceName string) string {
	var result strings.Builder

	for i, r := range interfaceName {
		if unicode.IsUpper(r) && i >= 1 {
			result.WriteRune('_')
		}

		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}
