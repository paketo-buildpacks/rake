package main

import (
	"os"

	"github.com/cloudfoundry/packit/scribe"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-community/rake"
)

func main() {
	parser := rake.NewGemfileParser()
	logger := scribe.NewLogger(os.Stdout)

	packit.Run(
		rake.Detect(parser),
		rake.Build(parser, logger),
	)
}
