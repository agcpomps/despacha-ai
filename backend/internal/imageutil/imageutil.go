// Package imageutil normalises uploaded images: it auto-orients them using
// their EXIF metadata, downscales oversized photos and re-encodes them as
// compressed JPEGs. Phone cameras produce 3-5MB files; after processing each
// image is typically well under 300KB, which keeps storage and bandwidth low.
package imageutil

import (
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

const (
	// MaxDimension is the longest side (px) kept after downscaling.
	MaxDimension = 1280
	// JPEGQuality balances visual quality against file size (0-100).
	JPEGQuality = 82
)

// SaveCompressed decodes the image in src, applies its EXIF orientation,
// fits it within MaxDimension x MaxDimension (preserving aspect ratio, never
// upscaling) and writes it to dir as a JPEG. It returns the generated file
// name. An error is returned when src is not a decodable image (e.g. HEIC),
// in which case nothing is written.
func SaveCompressed(src io.Reader, dir, baseName string) (string, error) {
	img, err := imaging.Decode(src, imaging.AutoOrientation(true))
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	if bounds.Dx() > MaxDimension || bounds.Dy() > MaxDimension {
		img = imaging.Fit(img, MaxDimension, MaxDimension, imaging.Lanczos)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	fileName := baseName + ".jpg"
	destination := filepath.Join(dir, fileName)
	if err := imaging.Save(img, destination, imaging.JPEGQuality(JPEGQuality)); err != nil {
		return "", err
	}

	return fileName, nil
}
