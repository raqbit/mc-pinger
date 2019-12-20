package mcpinger

import (
	"encoding/json"
)

// Minecraft chat component
// See: https://wiki.vg/Chat#Current_system_.28JSON_Chat.29
type RegularChatComponent struct {
	Text          string                 `json:"text"`          // Text content
	Bold          bool                   `json:"bold"`          // Component is emboldened
	Italic        bool                   `json:"italic"`        // Component is italicized
	Underlined    bool                   `json:"underlined"`    // Component is underlined
	Strikethrough bool                   `json:"strikethrough"` // Component is struck out
	Obfuscated    bool                   `json:"obfuscated"`    // Component randomly switches between characters of the same width
	Color         string                 `json:"color"`         // Contains the color for the component
	Extra         []RegularChatComponent `json:"extra"`         // RegularChatComponent siblings
}

// Wraps a RegularChatComponent for parsing both regular & string-only MOTD's
type ChatComponent struct {
	RegularChatComponent
}

func (c *ChatComponent) UnmarshalJSON(data []byte) error {
	var regular RegularChatComponent

	// The data starts with quotes which means it's a string, not an object
	if data[0] == '"' {
		var text string
		if err := json.Unmarshal(data, &text); err != nil {
			return err
		}

		regular.Text = text
	} else {
		if err := json.Unmarshal(data, &regular); err != nil {
			return err
		}
	}

	c.RegularChatComponent = regular

	return nil
}
