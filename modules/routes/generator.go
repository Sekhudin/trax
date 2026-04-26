package routes

import (
	"github.com/sekhudin/trax/internal/fs"
)

type generator struct {
	writer   fs.FileWriter
	template Template
}

type Generator interface {
	Generate(path string) error
}

func NewGenerator(w fs.FileWriter, t Template) Generator {
	return &generator{
		writer:   w,
		template: t,
	}
}

func (g *generator) Generate(path string) error {
	content, err := g.buildContent()
	if err != nil {
		return err
	}

	return g.write(path, content)
}

func (g *generator) buildContent() ([]byte, error) {
	data, err := g.template.Build()
	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

func (g *generator) write(path string, content []byte) error {
	return g.writer.Write(path, content)
}
