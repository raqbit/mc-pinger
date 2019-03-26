package mcpinger

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	TestDataDir = "testdata" // Name of testdata directory
)

func TestParseServerInfo(t *testing.T) {
	infoJson := GetTestFileContents(t, "info.json")

	info, err := parseServerInfo(infoJson)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", info)
}

func GetTestFileContents(t *testing.T, name string) []byte {
	path := filepath.Join(TestDataDir, name)

	data, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	return data
}
