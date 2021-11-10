package rake

import (
	"fmt"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
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

		command := "rake"

		if hasRakeGem {
			command = "bundle exec rake"
		}

		logger.Process("Assigning launch processes")
		logger.Subprocess("web: %s", command)

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					{
						Type:    "web",
						Command: command,
						Default: true,
					},
				},
			},
		}, nil
	}
}
