package tool

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

//
// The `packer` command is used to pack all the assets within a directory into a single packfile.
// The packfile format is as follows:
// - Each file is identified by its path relative to the input directory.
// - The path is written as a fixed length int32 followed by the path string.
// - The data length is written as a fixed length int32 followed by the file data.
// - E.g. <path_length><path><data_length><data><path_length><path><data_length><data>...
//

// packerCmd represents the packer command
var packerCmd = &cobra.Command{
	Use:   "packer <in_dir> <out_file>",
	Short: "Pack the <in_dir> directory into <out_file> packfile",
	Long: `Pack all the assets within the <in_dir> directory into a single <out_file> packfile.
They are identified with their path relative to <in_dir> excluding <in_dir> itself.
Loading from such a file is supported through the assets module.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inputDir := args[0]
		outputFile := args[1]

		var buffer bytes.Buffer

		// Walk the input directory and write the relative path and file data to the buffer.
		err := filepath.Walk(inputDir, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(inputDir, filePath)
			if err != nil {
				return err
			}

			// Write the relative path and data length to the buffer, as fixed length int32s.
			// This limits the maximum size of a single file to 2^31-1 bytes or 2GB.
			pathLength := int32(len(relPath))
			dataLength := int32(len(data))

			err = binary.Write(&buffer, binary.LittleEndian, pathLength)
			if err != nil {
				return err
			}
			_, err = buffer.WriteString(relPath)
			if err != nil {
				return err
			}
			err = binary.Write(&buffer, binary.LittleEndian, dataLength)
			if err != nil {
				return err
			}
			_, err = buffer.Write(data)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			panic(err)
		}

		err = os.WriteFile(outputFile, buffer.Bytes(), 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Assets packed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(packerCmd)
}
