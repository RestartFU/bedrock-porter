package porter

import "github.com/google/uuid"

// NewManifest creates a new manifest.
func NewManifest(name, description string) *Manifest {
	header := ManifestHeader{
		Description:    description,
		Name:           name,
		PlatformLocked: false,
		UUID:           uuid.NewString(),
		Version:        []int{0, 0, 1},
	}

	modules := []ManifestModule{
		{
			Description: description,
			Type:        "resources",
			UUID:        uuid.NewString(),
			Version:     []int{0, 0, 1},
		},
	}

	manifest := &Manifest{
		FormatVersion:    1,
		Header:           header,
		MinEngineVersion: []int{0, 0, 1, 1, 2, 5},
		Modules:          modules,
	}

	return manifest
}

// Manifest represents a Minecraft resource pack resource pack.
type Manifest struct {
	// FormatVersion is the format version of the manifest.
	FormatVersion int `json:"format_version"`
	// Header is the header of the manifest.
	Header ManifestHeader `json:"header"`
	// MinEngineVersion is the minimum engine version of the resource pack.
	MinEngineVersion []int `json:"min_engine_version"`
	// Modules is the modules of the resource pack.
	Modules []ManifestModule `json:"modules"`
}

// ManifestHeader represents a Minecraft resource pack manifest header.
type ManifestHeader struct {
	// Description is the description of the resource pack.
	Description string `json:"description"`
	// Name is the name of the resource pack.
	Name string `json:"name"`
	// PlatformLocked is the platform locked of the resource pack.
	PlatformLocked bool `json:"platform_locked"`
	// 		UUID is the UUID of the resource pack.
	UUID string `json:"uuid"`
	// Version is the version of the resource pack.
	Version []int `json:"version"`
}

// ManifestModule represents a Minecraft resource pack manifest module.
type ManifestModule struct {
	// Description is the description of the module.
	Description string `json:"description"`
	// Type is the type of the module.
	Type string `json:"type"`
	// UUID is the UUID of the module.
	UUID string `json:"uuid"`
	// Version is the version of the module.
	Version []int `json:"version"`
}
