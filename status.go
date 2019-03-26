package mcpinger

import "encoding/json"

// Minecraft chat component
// See: https://wiki.vg/Chat#Current_system_.28JSON_Chat.29
//
// Because the Notchian server does not provide actual chat component data,
// we only parse "text". This will use the old system, see: https://wiki.vg/Chat#Control_Sequences
type ChatComponent struct {
	Text string `json:"text"` // Text content
}

// Server info version
type Version struct {
	Name     string `json:"name"`     // Version name
	Protocol uint   `json:"protocol"` // Version protocol number
}

// Server info player
type Player struct {
	Name string `json:"name"` // Player name
	ID   string `json:"id"`   // Player UUID
}

// Server info players
type Players struct {
	Max    uint `json:"max"`    // Max amount of players allowed
	Online uint `json:"online"` // Amount of players online
	Sample []Player             // Sample of online players
}

// Server ping response
// https://wiki.vg/Server_List_Ping#Response
type ServerInfo struct {
	Version     Version       `json:"version"`     // Server version info
	Players     Players       `json:"players"`     // Server player info
	Description ChatComponent `json:"description"` // Server description
	Favicon     string        `json:"favicon"`     // Server favicon
}

// Parses the provided json byte array into a ServerInfo struct
func parseServerInfo(infoJson []byte) (*ServerInfo, error) {
	info := new(ServerInfo)
	err := json.Unmarshal(infoJson, info)
	return info, err
}
