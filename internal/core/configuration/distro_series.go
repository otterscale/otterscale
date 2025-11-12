package configuration

import "github.com/juju/juju/core/base"

type BootImageSelection struct {
	DisplayName   string
	DistroSeries  string
	Architectures []string
}

var DistroSeriesMap = map[base.SeriesName]BootImageSelection{
	base.Xenial: {
		DisplayName:   "Ubuntu 16.04 LTS Xenial Xerus",
		DistroSeries:  base.Xenial.String(),
		Architectures: []string{"amd64", "arm64", "armhf", "i386", "ppc64el", "s390x"},
	},
	base.Bionic: {
		DisplayName:   "Ubuntu 18.04 LTS Bionic Beaver",
		DistroSeries:  base.Bionic.String(),
		Architectures: []string{"amd64", "arm64", "armhf", "i386", "ppc64el", "s390x"},
	},
	base.Focal: {
		DisplayName:   "Ubuntu 20.04 LTS Focal Fossa",
		DistroSeries:  base.Focal.String(),
		Architectures: []string{"amd64", "arm64", "armhf", "ppc64el", "s390x"},
	},
	base.Jammy: {
		DisplayName:   "Ubuntu 22.04 LTS Jammy Jellyfish",
		DistroSeries:  base.Jammy.String(),
		Architectures: []string{"amd64", "arm64", "armhf", "ppc64el", "s390x"},
	},
	base.Noble: {
		DisplayName:   "Ubuntu 24.04 LTS Noble Numbat",
		DistroSeries:  base.Noble.String(),
		Architectures: []string{"amd64", "arm64", "armhf", "ppc64el", "s390x"},
	},
}
