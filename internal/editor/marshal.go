package editor

import (
	"encoding/json"

	"github.com/abhigyanwebber/cmd-customizer/internal/config"
)

// marshalTheme converts a Theme struct back to pretty JSON
func marshalTheme(t *config.Theme) ([]byte, error) {
	return json.MarshalIndent(t, "", "  ")
}
