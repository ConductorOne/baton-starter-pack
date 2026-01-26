package main

import (
	cfg "github.com/{{ cookiecutter.repo_owner }}/{{ cookiecutter.repo_name }}/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	// The first argument becomes the generated struct name (title-cased).
	// Using "connector" produces struct "Connector" in conf.gen.go.
	config.Generate("connector", cfg.Config)
}
