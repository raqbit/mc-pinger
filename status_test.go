package mcpinger

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	TestDataDir = "testdata" // Name of testdata directory
)

func TestParseServerInfo(t *testing.T) {
	tests := []struct {
		File              string
		VersionName       string
		Protocol          int32
		Description       string
		MaxPlayers        int32
		OnlinePlayers     int32
		PlayerSampleCount int
		PlayerSampleName  string
		PlayerSampleUuid  string
		Favicon           string
	}{
		{
			File:              "info.json",
			VersionName:       "1.13.2",
			Protocol:          404,
			Description:       "Hello world",
			MaxPlayers:        100,
			OnlinePlayers:     5,
			PlayerSampleCount: 1,
			PlayerSampleName:  "Raqbit",
			PlayerSampleUuid:  "09bc745b-3679-4152-b96b-3f9c59c42059",
			Favicon:           "data:image/png;base64,<data>",
		},
		{
			File:              "info_description_string.json",
			VersionName:       "1.13.2",
			Protocol:          404,
			Description:       "Hello world",
			MaxPlayers:        100,
			OnlinePlayers:     5,
			PlayerSampleCount: 1,
			PlayerSampleName:  "Raqbit",
			PlayerSampleUuid:  "09bc745b-3679-4152-b96b-3f9c59c42059",
			Favicon:           "data:image/png;base64,<data>",
		},
		{
			File:          "info_description_1_20_3.json",
			VersionName:   "Paper 1.20.4",
			Protocol:      765,
			Description:   "Foo",
			MaxPlayers:    20,
			OnlinePlayers: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.File, func(t *testing.T) {
			infoJson := GetTestFileContents(t, test.File)

			info, err := parseServerInfo(infoJson)

			if err != nil {
				t.Fatal(err)
			}

			if info.Version.Name != test.VersionName {
				parseError(t, test.File, "version name")
			}

			if info.Version.Protocol != test.Protocol {
				parseError(t, test.File, "protocol version")
			}

			if info.Description.Text != test.Description {
				parseError(t, test.File, "description")

			}

			if info.Players.Max != test.MaxPlayers {
				parseError(t, test.File, "max players")

			}

			if info.Players.Online != test.OnlinePlayers {
				parseError(t, test.File, "online players")
			}

			if len(info.Players.Sample) != test.PlayerSampleCount {
				parseError(t, test.File, "player sample")
			} else if test.PlayerSampleCount > 0 {
				if info.Players.Sample[0].Name != test.PlayerSampleName {
					parseError(t, test.File, "player sample name")
				}

				if info.Players.Sample[0].ID != test.PlayerSampleUuid {
					parseError(t, test.File, "player sample uuid")
				}

			}

			if info.Favicon != test.Favicon {
				parseError(t, test.File, "favicon")
			}

		})
	}
}

func parseError(t *testing.T, file string, name string) {
	t.Errorf("%s: Did not parse %s correctly", file, name)
}

func GetTestFileContents(t *testing.T, name string) []byte {
	path := filepath.Join(TestDataDir, name)

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	return data
}
