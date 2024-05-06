package settings

import (
	"encoding/json"
	"os"
)

type GameSettings struct {
	// Meta
	Version string `json:"version"` // 0.0.0
	Title   string `json:"title"`   // made with gorl
	// Display
	ScreenWidth     int  `json:"screenWidth"`     // 960
	ScreenHeight    int  `json:"screenHeight"`    // 540
	RenderWidth     int  `json:"renderWidth"`     // 960
	RenderHeight    int  `json:"renderHeight"`    // 540
	TargetFps       int  `json:"targetFps"`       // 144
	Fullscreen      bool `json:"fullscreen"`      // false
	EnableCrtEffect bool `json:"enableCrtEffect"` // true
	// Gameplay
	MouseSensitivity float32 `json:"mouseSensitivity"` // 1.0
	// Audio
	SoundVolume float32 `json:"soundVolume"` // 0.5
	// Logging
	LogPath string `json:"logPath"` // logs/
	// Controls
	EnableGamepad bool `json:"enableGamepad"` // false
}

var (
	settings *GameSettings
)

// Get the current settings
func CurrentSettings() *GameSettings {
	return settings
}

func UseFallbackSettings() {
	settings = &GameSettings{
		Version:  "0.0.0",
		Title:  "made with gorl",
		ScreenWidth:  960,
		ScreenHeight:  540,
		RenderWidth:  960,
		RenderHeight:  540,
		TargetFps:  144,
		Fullscreen:  false,
		EnableCrtEffect:  true,
		MouseSensitivity:  1.0,
		SoundVolume:  0.5,
		LogPath:  "logs/",
		EnableGamepad:  false,
	}
}

func LoadSettings(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	settings = new(GameSettings)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(settings)
	if err != nil {
		return err
	}

	return nil
}
