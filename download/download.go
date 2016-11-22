package download

import (
	log "github.com/Sirupsen/logrus"
	"github.com/dainis/go-playstore"
	"io/ioutil"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

type Config struct {
	Email    string
	Password string
	DeviceId string

	OutputFile string
}

func Download(c *Config, pkg string) {

	p, err := playstore.New(c.Email, c.Password, c.DeviceId)

	if err != nil {
		log.WithError(err).Fatal("Failed to connect to playstore")
	}

	details, err := p.PackageDetails(pkg)

	if err != nil {
		log.WithError(err).Fatal("Failed to obtain package details")
	}

	log.Infof("Will download %s with version %s", pkg, details.VersionString)

	b, err := p.DownloadPackage(pkg, details.VersionCode)

	if err != nil {
		log.WithError(err).Fatal("Failed to download apk file")
	}

	outputFile := pkg + "(" + details.VersionString + ")" + ".apk"

	if c.OutputFile != "" {
		outputFile = c.OutputFile
	}

	err = ioutil.WriteFile(outputFile, b, 0644)

	if err != nil {
		log.WithError(err).Fatalf("Failed to write to file %s", outputFile)
	}

	log.Infof("APK downloaded and saved as %s", outputFile)
}
