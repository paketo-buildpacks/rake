package rake

import (
	"fmt"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

//go:generate faux --interface Parser --output fakes/parser.go
type Parser interface {
	Parse(path string) (hasRakeGem bool, err error)
}

func Build(gemfileParser Parser, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		logger.Debug.Process("Checking Gemfile for rake")
		hasRakeGem, err := gemfileParser.Parse(filepath.Join(context.WorkingDir, "Gemfile"))
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to parse Gemfile: %w", err)
		}

		command := "rake"

		args := []string{}
		if hasRakeGem {
			logger.Debug.Subprocess("Gemfile contains rake gem")
			command = "bundle"
			args = []string{"exec rake"}
		}
		logger.Debug.Break()

		processes := []packit.Process{
			{
				Type:    "web",
				Command: command,
				Args:    args,
				Default: true,
				Direct:  true,
			},
		}
		logger.LaunchProcesses(processes)

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: processes,
			},
		}, nil
	}
}
