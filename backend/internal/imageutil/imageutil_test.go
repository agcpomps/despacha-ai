package imageutil

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// jpegOf gera um JPEG de teste com as dimensões dadas.
func jpegOf(t *testing.T, w, h int) *bytes.Buffer {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 120, 255})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
		t.Fatalf("encode fixture: %v", err)
	}
	return &buf
}

func savedDimensions(t *testing.T, path string) (int, int) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open saved: %v", err)
	}
	defer f.Close()
	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		t.Fatalf("decode saved: %v", err)
	}
	if format != "jpeg" {
		t.Errorf("formato gravado = %q, esperado jpeg", format)
	}
	return cfg.Width, cfg.Height
}

func TestSaveCompressedDownscalesLargeImage(t *testing.T) {
	dir := t.TempDir()
	src := jpegOf(t, 3000, 2000)
	originalSize := src.Len()

	name, err := SaveCompressed(src, dir, "large")
	if err != nil {
		t.Fatalf("SaveCompressed: %v", err)
	}
	if !strings.HasSuffix(name, ".jpg") {
		t.Errorf("nome = %q, esperado terminar em .jpg", name)
	}

	w, h := savedDimensions(t, filepath.Join(dir, name))
	if w > MaxDimension || h > MaxDimension {
		t.Errorf("não redimensionou: %dx%d (máx %d)", w, h, MaxDimension)
	}
	// rácio preservado: 3000x2000 (3:2) -> 1280x853
	if w != MaxDimension {
		t.Errorf("lado maior = %d, esperado %d", w, MaxDimension)
	}

	info, _ := os.Stat(filepath.Join(dir, name))
	if info.Size() >= int64(originalSize) {
		t.Errorf("comprimido (%d) não é menor que o original (%d)", info.Size(), originalSize)
	}
}

func TestSaveCompressedKeepsSmallImage(t *testing.T) {
	dir := t.TempDir()
	src := jpegOf(t, 400, 300)

	name, err := SaveCompressed(src, dir, "small")
	if err != nil {
		t.Fatalf("SaveCompressed: %v", err)
	}

	w, h := savedDimensions(t, filepath.Join(dir, name))
	if w != 400 || h != 300 {
		t.Errorf("imagem pequena foi alterada: %dx%d, esperado 400x300", w, h)
	}
}

func TestSaveCompressedRejectsNonImage(t *testing.T) {
	dir := t.TempDir()
	src := bytes.NewBufferString("isto não é uma imagem")

	if _, err := SaveCompressed(src, dir, "bad"); err == nil {
		t.Fatal("esperado erro para conteúdo não-imagem, obtido nil")
	}
}
