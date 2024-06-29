package settings

import (
	"encoding/json"
	"gorl/fw/core/assets"
)

type GameSettings struct {
	// Meta
	Version string `json:"version"` // 0.0.0
	Title   string `json:"title"`   // made with gorl
	// Display
	ScreenWidth     int  `json:"screenWidth"`     // 1920
	ScreenHeight    int  `json:"screenHeight"`    // 1080
	RenderWidth     int  `json:"renderWidth"`     // 1920
	RenderHeight    int  `json:"renderHeight"`    // 1080
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
		Version:          "0.0.0",
		Title:            "made with gorl",
		ScreenWidth:      1920,
		ScreenHeight:     1080,
		RenderWidth:      1920,
		RenderHeight:     1080,
		TargetFps:        144,
		Fullscreen:       false,
		EnableCrtEffect:  true,
		MouseSensitivity: 1.0,
		SoundVolume:      0.5,
		LogPath:          "logs/",
		EnableGamepad:    false,
	}
}

func LoadSettings(path string) error {
	file, err := assets.LoadFile(path)
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

// ============================================================================
// Utilities
// ============================================================================

// ScreenSizeV returns an rl.Vector2 with ScreenWidth and ScreenHeight
func (s *GameSettings) ScreenSizeV() rl.Vector2 {
	return rl.NewVector2(float32(s.ScreenWidth), float32(s.ScreenHeight))
}

// RenderSizeV returns an rl.Vector2 with RenderWidth and RenderHeight
func (s *GameSettings) RenderSizeV() rl.Vector2 {
	return rl.NewVector2(float32(s.RenderWidth), float32(s.RenderHeight))
}
