package colorscheme

import (
	"gorl/fw/core/logging"
	"image/color"
	"strconv"
	"strings"
)

type Hex string

var Colorscheme = struct {
	Color01 Hex
	Color02 Hex
	Color03 Hex
	Color04 Hex
	Color05 Hex
	Color06 Hex
	Color07 Hex
	Color08 Hex
	Color09 Hex
	Color10 Hex
	Color11 Hex
	Color12 Hex
	Color13 Hex
	Color14 Hex
	Color15 Hex
	Color16 Hex
}{
	"#9ece6a", // Color01
	"#9ece6a", // Color02
	"#9ece6a", // Color03
	"#a9b1d6", // Color04
	"#a9b1d6", // Color05
	"#7dcfff", // Color06
	"#c0c741", // Color07
	"#647d34", // Color08
	"#e4943a", // Color09
	"#db4b4b", // Color010
	"#d26471", // Color011
	"#70377f", // Color012
	"#414868", // Color013
	"#34859d", // Color014
	"#17434b", // Color015
	"#24283b", // Color016
}

// #24283b
// #c0caf5
// #1d202f
// #f7768e
// #9ece6a
// #e0af68
// #7aa2f7
// #bb9af7
// #7dcfff
// #a9b1d6
//
// #c0caf5
// #ff9e64
// #db4b4b

func (h Hex) ToRGBA() color.RGBA {
	rgba, err := hex2RGBA(h)
	if err != nil {
		logging.Error("Failed to convert hex to RGBA: %v", err)
		rgba = color.RGBA{ // A strong pink color for errors
			R: 255,
			G: 51,
			B: 194,
			A: 255,
		}
	}
	return rgba
}

func hex2RGBA(hex Hex) (color.RGBA, error) {
	var rgb color.RGBA
	values, err := strconv.ParseUint(strings.Trim(string(hex), "#"), 16, 32)

	if err != nil {
		return color.RGBA{}, err
	}

	rgb = color.RGBA{
		R: uint8(values >> 16),
		G: uint8((values >> 8) & 0xFF),
		B: uint8(values & 0xFF),
		A: 0xFF,
	}

	return rgb, nil
}
