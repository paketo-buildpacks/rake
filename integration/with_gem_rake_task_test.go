package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testWithGemRakeTask(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when building a simple rake exec container with the rake gem", func() {
		var (
			image     occam.Image
			container occam.Container

			name   string
			source string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		it("creates a working OCI image that can run the default rake task", func() {
			var err error
			source, err = occam.Source(filepath.Join("testdata", "with_gem_rake_task"))
			Expect(err).NotTo(HaveOccurred())

			var logs fmt.Stringer
			image, logs, err = pack.WithNoColor().Build.
				WithBuildpacks(
					settings.Buildpacks.MRI.Online,
					settings.Buildpacks.Bundler.Online,
					settings.Buildpacks.BundleInstall.Online,
					settings.Buildpacks.Rake.Online,
				).
				WithEnv(map[string]string{"BP_LOG_LEVEL": "DEBUG"}).
				WithPullPolicy("never").
				Execute(name, source)
			Expect(err).NotTo(HaveOccurred(), logs.String())

			container, err = docker.Container.Run.Execute(image.ID)
			Expect(err).NotTo(HaveOccurred())

			rLogs := func() string {
				rakeLogs, err := docker.Container.Logs.Execute(container.ID)
				Expect(err).NotTo(HaveOccurred())
				return rakeLogs.String()
			}

			Eventually(rLogs).Should(ContainSubstring("I am a rake task"))

			Expect(logs).To(ContainLines(
				MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, settings.Buildpack.Name)),
				"  Checking Gemfile for rake",
				"    Gemfile contains rake gem",
				"",
				"  Assigning launch processes:",
				`    web (default): bundle exec rake`,
			))
		})
	})
}
