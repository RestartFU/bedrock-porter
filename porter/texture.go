package porter

import (
	"archive/zip"
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
)

// Texture represents a Minecraft texture.
type Texture struct {
	image.Image
	name string
}

// NewTexture creates a new Texture from a name and reader.
func NewTexture(name string, reader io.Reader) (*Texture, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	t := &Texture{
		Image: fixAlpha(img),
		name:  name,
	}
	return t, nil
}

func fixAlpha(img image.Image) image.Image {
	newImage := image.NewNRGBA(img.Bounds())

	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			pixel := img.At(x, y)
			r, g, b, a := pixel.RGBA()
			// Convert RGBA to NRGBA
			pixel = color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}

			if a != 0 && a != 65535 {
				// If the alpha is not 0 or 65535, we need to fix the alpha channel.
				pixel = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
			newImage.Set(x, y, pixel)
		}
	}
	return newImage
}

// Copy copies the texture to a file.
func (t Texture) Copy(to string) error {
	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, t)
}

// CopyWriter copies the texture to a zip writer.
func (t Texture) CopyWriter(w *zip.Writer, path string) error {
	writer, err := w.Create(path)
	if err != nil {
		return err
	}

	return png.Encode(writer, t)
}
