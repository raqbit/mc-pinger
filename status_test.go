package mcpinger

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	TestDataDir       = "testdata" // Name of testdata directory
	VersionName       = "1.13.2"
	Protocol          = 404
	Description       = "Hello world"
	MaxPlayers        = 100
	OnlinePlayers     = 5
	PlayerSampleCount = 1
	PlayerSampleName  = "Raqbit"
	PlayerSampleUuid  = "09bc745b-3679-4152-b96b-3f9c59c42059"
	Favicon           = "data:image/png;base64,<data>"
)

func TestParseServerInfo(t *testing.T) {
	files := []string{
		"info.json",
		"info_description_string.json",
	}

	for _, f := range files {
		infoJson := GetTestFileContents(t, f)

		info, err := parseServerInfo(infoJson)

		if err != nil {
			t.Fatal(err)
		}

		if info.Version.Name != VersionName {
			parseError(t, f, "version name")
		}

		if info.Version.Protocol != Protocol {
			parseError(t, f, "protocol version")
		}

		if info.Description.Text != Description {
			parseError(t, f, "description")

		}

		if info.Players.Max != MaxPlayers {
			parseError(t, f, "max players")

		}

		if info.Players.Online != OnlinePlayers {
			parseError(t, f, "online players")
		}

		if len(info.Players.Sample) != PlayerSampleCount {
			parseError(t, f, "player sample")
		} else {
			if info.Players.Sample[0].Name != PlayerSampleName {
				parseError(t, f, "player sample name")
			}

			if info.Players.Sample[0].ID != PlayerSampleUuid {
				parseError(t, f, "player sample uuid")
			}

		}

		if info.Favicon != Favicon {
			parseError(t, f, "favicon")
		}
	}
}

func parseError(t *testing.T, file string, name string) {
	t.Errorf("%s: Did not parse %s correctly", file, name)
}

func GetTestFileContents(t *testing.T, name string) []byte {
	path := filepath.Join(TestDataDir, name)

	data, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	return data
}
