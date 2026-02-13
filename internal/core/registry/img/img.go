package img

import ocispec "github.com/opencontainers/image-spec/specs-go/v1"

type (
	// Image represents a OCI Image resource.
	Image = ocispec.Image

	// ImageConfig represents a OCI ImageConfig resource.
	//
	//nolint:revive // allows this type alias for brevity.
	ImageConfig = ocispec.ImageConfig

	// Platform represents a OCI Platform resource.
	Platform = ocispec.Platform

	// RootFS represents a OCI RootFS resource.
	RootFS = ocispec.RootFS
)
