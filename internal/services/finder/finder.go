package finder

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/samber/lo"
	"github.com/ythosa/gowrapper/internal/logging"
	"github.com/ythosa/gowrapper/internal/model"
)

type InterfacesFinder struct {
	logger logging.Logger
}

func NewInterfacesFinder() *InterfacesFinder {
	return &InterfacesFinder{logger: logging.New("finder")}
}

func (g *InterfacesFinder) Find(initial string, excludedDirs []string, excludedFiles []string) model.InterfacesByFileMap {
	var (
		files            = g.getFiles(initial, excludedDirs, excludedFiles)
		interfacesByFile = make(model.InterfacesByFileMap, len(files))
	)

	lo.ForEach(files, func(f string, _ int) {
		if interfaces := g.findInterfacesInFile(f); interfaces != nil {
			interfacesByFile[f] = interfaces
		}
	})

	g.logFoundInterfaces(interfacesByFile)

	return interfacesByFile
}

func (g *InterfacesFinder) getFiles(initial string, excludedDirs []string, excludedFiles []string) (paths []string) {
	err := filepath.WalkDir(initial, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			g.logger.Warnf("traversing dirs: %s", err)

			return nil
		}

		if d.IsDir() && g.isExcluded(path, excludedDirs) {
			return filepath.SkipDir
		}

		if !d.IsDir() && !g.isExcluded(path, excludedFiles) {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		g.logger.Fatalf("failed to traverse initial directory: %s", err)
	}

	return paths
}

func (g *InterfacesFinder) isExcluded(path string, paths []string) bool {
	return lo.SomeBy(paths, func(p string) bool { return strings.Contains(path, p) })
}

func (g *InterfacesFinder) findInterfacesInFile(path string) (interfaces []string) {
	fast, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
	if err != nil {
		log.Print(err)
	}

	for _, decl := range fast.Decls {
		stmt, isGenDecl := decl.(*ast.GenDecl)
		if !isGenDecl {
			continue
		}

		if stmt.Tok != token.TYPE {
			continue
		}

		for _, spec := range stmt.Specs {
			tspec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			_, ok = tspec.Type.(*ast.InterfaceType)
			if ok {
				interfaces = append(interfaces, tspec.Name.Name)
			}
		}
	}

	return interfaces
}

func (g *InterfacesFinder) logFoundInterfaces(interfacesByFile model.InterfacesByFileMap) {
	tbl := table.New("#", "Path", "Interfaces").
		WithHeaderFormatter(color.New(color.FgGreen, color.Underline).SprintfFunc()).
		WithFirstColumnFormatter(color.New(color.FgYellow).SprintfFunc())

	var i int
	for file, interfaces := range interfacesByFile {
		i++
		tbl.AddRow(i, file, strings.Join(interfaces, ", "))
	}

	var resultBuffer = bytes.NewBufferString("Found interfaces: \n")
	tbl.WithWriter(resultBuffer)
	tbl.Print()

	g.logger.Info(resultBuffer.String())
}
