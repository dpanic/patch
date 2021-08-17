package generator

import (
	"bytes"
	"os"
	"strings"
	"text/template"
	"time"

	"go.uber.org/zap"
)

const (
	templateHeader = `#!/usr/bin/env bash
	#
	# PATCHFILES SCRIPT FOR {{.ScriptFor}}
	# 
	# author: {{.Author}}
	# version: {{.Version}}
	# environment: {{.Environment}}
	# built: {{.Built}}
	#
	#

	args=("$@")
	category="${args[0]}"

	{{ if eq .ScriptFor "PATCHING" }}
		if test -f "{{.PatchFilesControlFile}}"; then
			echo "System already patched exiting"
			exit 0
		fi
	{{ end }}	

	{{ if eq .ScriptFor "REVERTING" }}
		if test ! -f "{{.PatchFilesControlFile}}"; then
			echo "System is not patched. Exiting."
			exit 0
		fi
	{{ end }}	

	`
)

type Header struct {
	ScriptFor             string
	Author                string
	Version               string
	Environment           string
	Built                 string
	PatchFilesControlFile string
}

// generateHeader generates header based on input parameters
func (generator *Generator) writeHeader(fd *os.File, scriptFor string) (err error) {
	logger := generator.Log.WithOptions(zap.Fields())
	logger.Debug("attempt to write footer",
		zap.String("scriptFor", scriptFor),
	)

	built := time.Now().UTC().Format("2006-01-02 15:04:05 -07:00")

	author := os.Getenv("AUTHOR")
	author = strings.ToLower(author)
	author = strings.Trim(author, " ")

	version := os.Getenv("VERSION")
	version = strings.ToLower(version)
	version = strings.Trim(version, " ")

	data := Header{
		Author:                author,
		Version:               version,
		Built:                 built,
		ScriptFor:             scriptFor,
		Environment:           generator.Environment,
		PatchFilesControlFile: patchFilesControlFile,
	}

	var (
		buf = new(bytes.Buffer)
	)

	tpl, err := template.New("template").Parse(templateHeader)

	t := template.Must(tpl, err)
	err = t.Execute(buf, data)
	if err != nil {
		return
	}

	res := buf.String()
	res = strings.ReplaceAll(res, "\t", "")

	fd.WriteString(res)
	fd.Sync()

	return
}
