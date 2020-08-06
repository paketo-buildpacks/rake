package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var settings struct {
	Buildpacks struct {
		Rake struct {
			Online string
		}
		MRI struct {
			Online string
		}
		Bundler struct {
			Online string
		}
		BundleInstall struct {
			Online string
		}
	}

	Config struct {
		MRI           string `json:"mri"`
		Bundler       string `json:"bundler"`
		BundleInstall string `json:"bundle-install"`
	}
}

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	root, err := filepath.Abs("./..")
	Expect(err).NotTo(HaveOccurred())

	file, err := os.Open("../integration.json")
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(file).Decode(&settings.Config)).To(Succeed())
	Expect(file.Close()).To(Succeed())

	buildpackStore := occam.NewBuildpackStore()

	settings.Buildpacks.Rake.Online, err = buildpackStore.Get.
		WithVersion("1.2.3").
		Execute(root)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.MRI.Online, err = buildpackStore.Get.
		Execute(settings.Config.MRI)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Bundler.Online, err = buildpackStore.Get.
		Execute(settings.Config.Bundler)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.BundleInstall.Online, err = buildpackStore.Get.
		Execute(settings.Config.BundleInstall)
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(10 * time.Second)

	suite := spec.New("Integration", spec.Report(report.Terminal{}), spec.Parallel())
	suite("WithoutGemRakeTask", testWithoutGemRakeTask)
	suite("WithGemRakeTask", testWithGemRakeTask)
	suite.Run(t)
}
