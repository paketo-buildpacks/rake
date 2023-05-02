package rake_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/rake"
	"github.com/paketo-buildpacks/rake/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir    string
		gemfileParser *fakes.Parser
		detect        packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		gemfileParser = &fakes.Parser{}

		detect = rake.Detect(gemfileParser)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("when there is a rake gem installed", func() {
		it.Before(func() {
			gemfileParser.ParseCall.Returns.HasRakeGem = true
		})

		context("when there is a Rakefile", func() {
			it.Before(func() {
				err := os.WriteFile(filepath.Join(workingDir, "Rakefile"), nil, 0600)
				Expect(err).NotTo(HaveOccurred())
			})

			it("detects and requires gems, bundler, and mri", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						{
							Name: "mri",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "gems",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "bundler",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
					},
				}))
			})
		})

		context("when there is a rakefile", func() {
			it.Before(func() {
				err := os.WriteFile(filepath.Join(workingDir, "rakefile"), nil, 0600)
				Expect(err).NotTo(HaveOccurred())
			})

			it("detects and requires gems, bundler, and mri", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						{
							Name: "mri",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "gems",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "bundler",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
					},
				}))
			})
		})

		context("when there is a rakefile.rb", func() {
			it.Before(func() {
				err := os.WriteFile(filepath.Join(workingDir, "rakefile.rb"), nil, 0600)
				Expect(err).NotTo(HaveOccurred())
			})

			it("detects and requires gems, bundler, and mri", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						{
							Name: "mri",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "gems",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "bundler",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
					},
				}))
			})
		})

		context("when there is a Rakefile.rb", func() {
			it.Before(func() {
				err := os.WriteFile(filepath.Join(workingDir, "Rakefile.rb"), nil, 0600)
				Expect(err).NotTo(HaveOccurred())
			})

			it("detects and requires gems, bundler, and mri", func() {
				result, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Plan).To(Equal(packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{},
					Requires: []packit.BuildPlanRequirement{
						{
							Name: "mri",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "gems",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
						{
							Name: "bundler",
							Metadata: rake.BuildPlanMetadata{
								Launch: true,
							},
						},
					},
				}))
			})
		})
	})

	context("when the Gemfile does not list rake and there is a Rakefile", func() {
		it.Before(func() {
			gemfileParser.ParseCall.Returns.HasRakeGem = false

			err := os.WriteFile(filepath.Join(workingDir, "Rakefile"), nil, 0600)
			Expect(err).NotTo(HaveOccurred())
		})

		it("detects and only requires mri", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "mri",
						Metadata: rake.BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			}))
		})
	})

	context("when there is not a rakefile", func() {
		it("fails detection", func() {
			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).To(MatchError(packit.Fail.WithMessage("no 'Rakefile', 'Rakefile.rb', 'rakefile', or 'rakefile.rb' file found")))
		})
	})

	context("failure cases", func() {
		context("the workspace directory cannot be accessed", func() {
			it.Before(func() {
				Expect(os.Chmod(workingDir, 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to stat Rakefile:")))
			})
		})

		context("when the gemfile parser fails", func() {
			it.Before(func() {
				gemfileParser.ParseCall.Returns.Err = errors.New("some-error")

				err := os.WriteFile(filepath.Join(workingDir, "Rakefile"), nil, 0600)
				Expect(err).NotTo(HaveOccurred())
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError("failed to parse Gemfile: some-error"))
			})
		})
	})
}
