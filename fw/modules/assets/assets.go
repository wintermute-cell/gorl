package assets

import (
	"bytes"
	"errors"
	"gorl/fw/core/logging"
	"io"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var usingPackfile bool
var packfile Packfile

// UsePackfile enables the use of a packfile.
func UsePackfile() {
	var err error
	packfile, err = LoadPackfile()
	if err != nil {
		logging.Fatal("Failed to load packfile: %v", err)
	}
	usingPackfile = true
}

// LoadFile returns a file handle, either from the packfile or from disk.
// Can be used as a drop-in replacement for os.Open.
func LoadFile(path string) (io.ReadCloser, error) {
	if usingPackfile {
		data, err := packfile.GetAssetBytes(path)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(bytes.NewReader(data)), nil
	}
	return os.Open(path)
}

// LoadTexture loads a Texture2D.
func LoadTexture(path string) (rl.Texture2D, error) {

	// If using packfile, load from packfile.
	if usingPackfile {
		data, err := packfile.GetAssetBytes(path)
		if err != nil {
			return rl.Texture2D{}, err
		}

		image := rl.LoadImageFromMemory(".png", data, int32(len(data)))
		if image == nil {
			return rl.Texture2D{}, errors.New("failed to load image")
		}

		tex := rl.LoadTextureFromImage(image)
		if tex.ID == 0 {
			return rl.Texture2D{}, errors.New("failed to load texture")
		}
		return tex, nil
	}

	// Otherwise, load from disk using native raylib.
	tex := rl.LoadTexture(path)
	if tex.ID == 0 {
		return rl.Texture2D{}, errors.New("failed to load texture")
	}
	return tex, nil
}
