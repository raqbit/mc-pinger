package mcpinger

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	TestDataDir = "testdata" // Name of testdata directory
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

		if info.Version.Name != "1.13.2" {
			parseError(t, f, "version name")
		}

		if info.Version.Protocol != 404 {
			parseError(t, f, "protocol version")
		}

		if info.Description.Text != "Hello world" {
			parseError(t, f, "description")

		}

		if info.Players.Max != 100 {
			parseError(t, f, "max players")

		}

		if info.Players.Online != 5 {
			parseError(t, f, "online players")
		}

		if len(info.Players.Sample) != 1 {
			parseError(t, f, "player sample")
		} else {
			if info.Players.Sample[0].Name != "Raqbit" {
				parseError(t, f, "player sample name")
			}

			if info.Players.Sample[0].ID != "09bc745b-3679-4152-b96b-3f9c59c42059" {
				parseError(t, f, "player sample uuid")
			}

		}

		if info.Favicon != "data:image/png;base64,<data>" {
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
