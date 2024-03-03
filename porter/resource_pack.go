package porter

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/restartfu/bedrock-porter/porter/frontend"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ResourcePack represents a Minecraft resource pack.
type ResourcePack struct {
	*zip.Reader
	mcpack *zip.Writer

	icon *Texture

	items  []*Texture
	blocks []*Texture

	armor []*Texture

	packName     string
	outputPath   string
	texturesPath string
}

// NewResourcePack creates a new ResourcePack from a path.
func NewResourcePack(path string) (*ResourcePack, error) {
	packZip, err := openZip(path)
	packName := fileNameFromPath(path)
	outputPath := "output/" + packName + "/" + time.Now().Format("2006-01-02-15-04-05") + "/"

	_ = os.MkdirAll(outputPath, 0755)
	mcpack, err := createZip(outputPath + packName + ".mcpack")
	if err != nil {
		return nil, err
	}

	rp := &ResourcePack{
		Reader: packZip,
		mcpack: mcpack,

		packName:   packName,
		outputPath: outputPath,
	}

	textures, ok := rp.resolveDirectory("textures")
	if !ok {
		return nil, errors.New("textures directory not found")
	}
	rp.texturesPath = textures

	items, err := rp.resolveDirectoryTextures("items")
	if err != nil {
		return nil, err
	}

	blocks, err := rp.resolveDirectoryTextures("blocks")
	if err != nil {
		return nil, err
	}

	armor, err := rp.resolveDirectoryTextures("models/armor")
	if err != nil {
		return nil, err
	}

	rp.items = items
	rp.blocks = blocks
	rp.armor = armor

	return rp, err
}

// Port ports the resource pack to the output directory.
func (r *ResourcePack) Port() {
	_ = r.generateManifest()

	err := r.portSky()
	if err != nil {
		fmt.Println(err)

	}
	_ = r.portSingleTexture("pack.png", "pack_icon.png")
	_ = r.portSingleTexture(r.texturesPath+"/particle/particles.png", "textures/particle/particles.png")
	//_ = r.portSingleTexture(r.texturesPath+"/gui/icons.png", "textures/gui/icons.png")

	r.portTextures("textures/items", r.items, ignoreNone)
	r.portTextures("textures/blocks", r.blocks, ignoreNone)
	r.portTextures("textures/armor", r.armor, func(s string) bool {
		return strings.Contains(s, "overlay")
	})

	_ = r.mcpack.Close()

	fmt.Println(frontend.Style.Render("\n Porting complete."))
	<-time.Tick(10 * time.Second)
}

// portTextures ports textures to the output directory.
func (r *ResourcePack) portTextures(path string, textures []*Texture, ignore func(string) bool) {
	_ = os.MkdirAll(r.outputPath+path, 0755)
	count := len(textures)
	fmt.Println(frontend.Style.Render("\n Porting " + path + "\n"))

	for i, t := range textures {
		if ignore(t.name) {
			count--
			continue
		}

		target := replacer.Replace(t.name)
		_ = t.Copy(r.outputPath + path + "/" + target)
		_ = t.CopyWriter(r.mcpack, path+"/"+target)

		if rand.Intn(100) > 50 {
			<-time.Tick(time.Millisecond)
		}

		frontend.Spinner, _ = frontend.Spinner.Update(spinner.TickMsg{
			Time: time.Now(),
		})

		fmt.Printf("%s %s %s", frontend.ClearLine, frontend.Spinner.View(), frontend.ProgressBar.ViewAs(float64(i+1)/float64(count)))
		if i == count-1 {
			fmt.Println()
		}

	}
}

// Texture represents a Minecraft texture.
func (r *ResourcePack) resolveDirectoryTextures(path string) ([]*Texture, error) {
	var textures []*Texture

	err := fs.WalkDir(r, r.texturesPath+"/"+path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".png") {
			return nil
		}

		reader, err := r.Open(path)
		if err != nil {
			return err
		}

		texture, err := NewTexture(d.Name(), reader)
		if err != nil {
			return err
		}

		textures = append(textures, texture)
		return nil
	})

	return textures, err
}

// portSingleTexture resolves the icon of the resource pack.
func (r *ResourcePack) portSingleTexture(path, newPath string) error {
	reader, err := r.Open(path)
	if err != nil {
		return err
	}

	icon, err := NewTexture(path, reader)
	if err != nil {
		return err
	}

	_ = os.MkdirAll(r.outputPath+directoryPathFromPath(newPath), 0755)
	_ = icon.Copy(r.outputPath + newPath)
	_ = icon.CopyWriter(r.mcpack, newPath)
	return nil
}

// resolveDirectory resolves a directory in the resource pack.
func (r *ResourcePack) resolveDirectory(name string) (dirPath string, found bool) {
	_ = fs.WalkDir(r, ".", func(path string, d fs.DirEntry, err error) error {
		if d.Name() == name && d.IsDir() {
			dirPath = path
			found = true
		}
		return nil
	})
	return
}

// generateManifest generates the manifest for the resource pack.
func (r *ResourcePack) generateManifest() error {
	manifest := NewManifest(r.packName, "Ported by Bedrock Porter")

	manifestFile, err := os.Create(r.outputPath + "manifest.json")
	if err != nil {
		return err
	}
	defer manifestFile.Close()

	manifestMcpack, err := r.mcpack.Create("manifest.json")
	if err != nil {
		return err
	}

	return writeManifest(manifest, manifestFile, manifestMcpack)

}

// writeManifest writes the manifest to the output directory.
func writeManifest(m *Manifest, w ...io.Writer) error {
	for _, writer := range w {
		enc := json.NewEncoder(writer)
		enc.SetIndent("", "  ")
		if err := enc.Encode(m); err != nil {
			return err
		}
	}
	return nil
}
