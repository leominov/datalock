package handlers

import (
	"html/template"
	"path"

	"github.com/leominov/datalock/server"
)

var (
	SecuredPlayerTemplate *template.Template
	PlayerPageTemplate    *template.Template
)

func ParseTemplates(config *server.Config) error {
	var err error
	if PlayerPageTemplate, err = template.New("master").ParseFiles(path.Join(config.TemplatesDir, "standard_player.html")); err != nil {
		return err
	}
	if SecuredPlayerTemplate, err = template.New("master").ParseFiles(path.Join(config.TemplatesDir, "secured_player.html")); err != nil {
		return err
	}
	return nil
}
