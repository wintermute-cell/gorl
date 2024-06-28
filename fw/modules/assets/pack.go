package assets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

type Packfile map[string][]byte

const packFilePath = "data.pack"

// LoadPackfile loads the asset pack file into memory.
func LoadPackfile() (Packfile, error) {
	data, err := os.ReadFile(packFilePath)
	if err != nil {
		return nil, err
	}

	packfile := make(Packfile)

	buffer := bytes.NewReader(data)
	for {
		var pathLength int32
		err := binary.Read(buffer, binary.LittleEndian, &pathLength)
		if err != nil {
			break
		}

		path := make([]byte, pathLength)
		_, err = buffer.Read(path)
		if err != nil {
			return nil, err
		}

		var dataLength int32
		err = binary.Read(buffer, binary.LittleEndian, &dataLength)
		if err != nil {
			return nil, err
		}

		data := make([]byte, dataLength)
		_, err = buffer.Read(data)
		if err != nil {
			return nil, err
		}

		packfile[string(path)] = data
	}

	return packfile, nil
}

// GetAssetBytes retrieves the asset data bytes for the given path.
func (pf Packfile) GetAssetBytes(path string) ([]byte, error) {
	data, ok := pf[path]
	if !ok {
		msg := "asset not found in packfile: " + path
		return nil, errors.New(msg)
	}

	return data, nil
}
