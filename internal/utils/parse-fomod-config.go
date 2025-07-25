package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html/charset"
)

type ModuleConfig struct {
	InstallSteps []InstallStep `xml:"installSteps>installStep"`
}

type InstallStep struct {
	Groups []PluginGroup `xml:"optionalFileGroups>group"`
}

type PluginGroup struct {
	Plugins []Plugin `xml:"plugins>plugin"`
}

type Plugin struct {
	Name    string       `xml:"name,attr"`
	Folders []PluginFile `xml:"files>folder"`
	Files   []PluginFile `xml:"files>file"`
}

type PluginFile struct {
	Source      string `xml:"source,attr"`
	Destination string `xml:"destination,attr"`
}

func ParseFomodConfig(path string) (*ModuleConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read ModuleConfig.xml: %w", err)
	}

	var reader io.Reader
	switch {
	case bytes.HasPrefix(raw, []byte{0xFF, 0xFE}):
		reader = bytes.NewReader(raw[2:])
		reader, err = charset.NewReaderLabel("utf-16le", reader)
	case bytes.HasPrefix(raw, []byte{0xFE, 0xFF}):
		reader = bytes.NewReader(raw[2:])
		reader, err = charset.NewReaderLabel("utf-16be", reader)
	default:
		reader = bytes.NewReader(raw)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode FOMOD XML: %w", err)
	}

	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var config ModuleConfig
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to parse FOMOD config XML: %w", err)
	}

	return &config, nil
}
