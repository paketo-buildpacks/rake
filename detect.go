package rake

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
)

type BuildPlanMetadata struct {
	Launch bool `toml:"launch"`
}

func Detect(gemfileParser Parser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		rakefiles := []string{"Rakefile", "Rakefile.rb", "rakefile", "rakefile.rb"}
		rakeFileExists := false
		for _, file := range rakefiles {
			_, err := os.Stat(filepath.Join(context.WorkingDir, file))
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return packit.DetectResult{}, fmt.Errorf("failed to stat %s: %w", file, err)
			} else {
				rakeFileExists = true
				break
			}
		}

		if !rakeFileExists {
			return packit.DetectResult{}, packit.Fail
		}

		hasRakeGem, err := gemfileParser.Parse(filepath.Join(context.WorkingDir, "Gemfile"))
		if err != nil {
			return packit.DetectResult{}, fmt.Errorf("failed to parse Gemfile: %w", err)
		}

		requirements := []packit.BuildPlanRequirement{}

		requirements = append(requirements, packit.BuildPlanRequirement{
			Name: "mri",
			Metadata: BuildPlanMetadata{
				Launch: true,
			},
		})

		if hasRakeGem {
			requirements = append(
				requirements,
				packit.BuildPlanRequirement{
					Name: "gems",
					Metadata: BuildPlanMetadata{
						Launch: true,
					},
				},
				packit.BuildPlanRequirement{
					Name: "bundler",
					Metadata: BuildPlanMetadata{
						Launch: true,
					},
				})
		}

		fmt.Println("REQUIREMENTS: ", requirements)

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: requirements,
			},
		}, nil
	}

}
