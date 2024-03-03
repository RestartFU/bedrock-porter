package porter

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/magiconair/properties"
	"github.com/restartfu/bedrock-porter/porter/frontend"
	"image"
	"image/png"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type SkyProperties struct {
	Rotate bool   `properties:"rotate"`
	Source string `properties:"source"`
}

func (r *ResourcePack) portSky() error {
	path, found := r.resolveDirectory("world0")
	if !found {
		return nil
	}

	skyProperties, err := r.skyProperties(path)
	if err != nil {
		return err
	}

	replacer := strings.NewReplacer("./", "")

	for _, sky := range skyProperties {
		if sky.Source == "" {
			continue
		}
		cubemaps, err := r.Open(path + "/" + replacer.Replace(sky.Source))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
		defer cubemaps.Close()

		img, err := png.Decode(cubemaps)
		if err != nil {
			return err
		}

		fmt.Println(frontend.Style.Render(fmt.Sprintf("\n Porting sky (%s)\n", path+"/"+replacer.Replace(sky.Source))))
		results, err := splitCubeMaps(img)
		if err != nil {
			return err
		}

		_ = os.MkdirAll(r.outputPath+"textures/environment/overworld_cubemap", 0755)
		for i, result := range results {
			var rotation int
			switch i {
			case 0:
				rotation = 5
			case 1:
				rotation = 4
			case 2:
				rotation = 2
			case 3:
				rotation = 3
			case 4:
				rotation = 0
			case 5:
				rotation = 1
			}

			filename := "textures/environment/overworld_cubemap/" + "cubemap_" + strconv.Itoa(rotation) + ".png"

			f, err := os.Create(r.outputPath + filename)
			if err != nil {
				return err
			}

			zipFile, err := r.mcpack.Create(filename)

			_ = png.Encode(zipFile, result)
			_ = png.Encode(f, result)
			_ = f.Close()
		}
		return nil
	}
	return nil
}

func splitCubeMaps(cubemaps image.Image) ([6]*image.RGBA, error) {
	var images [6]*image.RGBA

	cubeMapHeight := cubemaps.Bounds().Dy() / 2
	cubeMapWidth := cubemaps.Bounds().Dx() / 3

	for i := 0; i < 6; i++ {
		images[i] = image.NewRGBA(image.Rect(0, 0, cubeMapWidth, cubeMapHeight))

		for y := 0; y < cubeMapHeight; y++ {
			for x := 0; x < cubeMapWidth; x++ {
				images[i].Set(x, y, cubemaps.At(x+(i%3)*cubeMapWidth, y+(i/3)*cubeMapHeight))
			}
		}

		if rand.Intn(100) > 50 {
			<-time.Tick(time.Second / 2)
		}

		frontend.Spinner, _ = frontend.Spinner.Update(spinner.TickMsg{
			Time: time.Now(),
		})

		fmt.Printf("%s %s %s", frontend.ClearLine, frontend.Spinner.View(), frontend.ProgressBar.ViewAs(float64(i+1)/float64(6)))
	}

	return images, nil
}

func (r *ResourcePack) skyProperties(path string) ([]SkyProperties, error) {
	var skyProperties []SkyProperties

	err := fs.WalkDir(r, path, func(filePath string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(filePath, ".properties") {
			f, err := r.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()

			buf, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			p, err := properties.Load(buf, properties.UTF8)
			var sky SkyProperties
			err = p.Decode(&sky)
			if err != nil {
				return err
			}

			skyProperties = append(skyProperties, sky)
			return nil
		}
		return nil
	})
	return skyProperties, err
}
