package porter

import (
	"archive/zip"
	"bytes"
	"image"
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
		Image: img,
		name:  name,
	}
	return t, nil
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
