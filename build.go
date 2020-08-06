package rake

import (
	"fmt"
	"path/filepath"

	"github.com/cloudfoundry/packit/scribe"
	"github.com/paketo-buildpacks/packit"
)

//go:generate faux --interface Parser --output fakes/parser.go
type Parser interface {
	Parse(path string) (hasRakeGem bool, err error)
}

func Build(gemfileParser Parser, logger scribe.Logger) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		hasRakeGem, err := gemfileParser.Parse(filepath.Join(context.WorkingDir, "Gemfile"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to parse Gemfile: %w", err)
		}

		var command string
		if hasRakeGem {
			command = "bundle exec rake"
		} else {
			command = "rake"
		}

		logger.Process("Writing rake command")
		logger.Subprocess(command)

		return packit.BuildResult{
			Processes: []packit.Process{
				{
					Type:    "web",
					Command: command,
				},
			},
		}, nil
	}
}
