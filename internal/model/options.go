package model

type Options struct {
	InitialDirectory string
	ExcludedDirs     []string
	ExcludedFiles    []string

	TemplatePath string
	OutputFolder string
	FilePostfix  string
}
