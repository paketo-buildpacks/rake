package rake_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/paketo-community/rake"
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
		file, err := ioutil.TempFile("", "Gemfile")
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
				const GEMFILE_CONTENTS = `source 'https://rubygems.org'

gem 'rake'`

				Expect(ioutil.WriteFile(path, []byte(GEMFILE_CONTENTS), 0644)).To(Succeed())

				hasRake, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasRake).To(Equal(true))
			})
		})

		context("when not using rake", func() {
			it("parses correctly", func() {
				const GEMFILE_CONTENTS = `source 'https://rubygems.org'
ruby '~> 2.0'`

				Expect(ioutil.WriteFile(path, []byte(GEMFILE_CONTENTS), 0644)).To(Succeed())

				hasRake, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasRake).To(Equal(false))
			})
		})

		context("when the Gemfile file does not exist", func() {
			it.Before(func() {
				Expect(os.Remove(path)).To(Succeed())
			})

			it("returns all false", func() {
				hasRake, err := parser.Parse(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(hasRake).To(Equal(false))
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
