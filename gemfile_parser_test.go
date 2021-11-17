package rake_test

import (
	"os"
	"testing"

	"github.com/paketo-buildpacks/rake"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testGemfileParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		path   string
		parser rake.GemfileParser
	)

	it.Before(func() {
		file, err := os.CreateTemp("", "Gemfile")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		path = file.Name()
		parser = rake.NewGemfileParser()
	})

	it.After(func() {
		Expect(os.RemoveAll(path)).To(Succeed())
	})

	context("Parse", func() {
		context("when using rake", func() {
			it("parses correctly", func() {
				Expect(os.WriteFile(path, []byte(`source 'https://rubygems.org'

gem 'rake'`), 0600)).To(Succeed())

				hasRake, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasRake).To(BeTrue())
			})
		})

		context("when not using rake", func() {
			it("parses correctly", func() {
				Expect(os.WriteFile(path, []byte(`source 'https://rubygems.org'
ruby '~> 2.0'`), 0600)).To(Succeed())

				hasRake, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasRake).To(BeFalse())
			})
		})

		context("when the Gemfile file does not exist", func() {
			it.Before(func() {
				Expect(os.Remove(path)).To(Succeed())
			})

			it("returns all false", func() {
				hasRake, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasRake).To(BeFalse())
			})
		})

		context("failure cases", func() {
			context("when the Gemfile cannot be opened", func() {
				it.Before(func() {
					Expect(os.Chmod(path, 0000)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := parser.Parse(path)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failed to parse Gemfile:")))
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})
		})
	})
}
