package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/net/html/charset"
)

// Root FOMOD configuration structure
type ModuleConfig struct {
	XMLName                 xml.Name                 `xml:"config"`
	ModuleName              string                   `xml:"moduleName"`
	ModuleImage             string                   `xml:"moduleImage"`
	ModuleDependencies      *ModuleDependencies      `xml:"moduleDependencies"`
	RequiredInstallFiles    *RequiredInstallFiles    `xml:"requiredInstallFiles"`
	InstallSteps            []InstallStep            `xml:"installSteps>installStep"`
	ConditionalFileInstalls []ConditionalFileInstall `xml:"conditionalFileInstalls>patterns>pattern"`
}

// Module dependencies structure
type ModuleDependencies struct {
	Operator     string       `xml:"operator,attr"` // "And", "Or"
	Dependencies []Dependency `xml:",any"`
}

// Required files installed regardless of user choice
type RequiredInstallFiles struct {
	Files   []FileInstall   `xml:"file"`
	Folders []FolderInstall `xml:"folder"`
}

// Install step structure with full conditional support
type InstallStep struct {
	Name    string               `xml:"name,attr"`
	Order   string               `xml:"order,attr"` // "Explicit", "Alphabetical"
	Visible *VisibilityCondition `xml:"visible"`
	Groups  []PluginGroup        `xml:"optionalFileGroups>group"`
}

// Visibility conditions for steps
type VisibilityCondition struct {
	Operator     string       `xml:"operator,attr"` // "And", "Or"
	Dependencies []Dependency `xml:",any"`
}

// Plugin group with selection type
type PluginGroup struct {
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`  // "SelectExactlyOne", "SelectAll", "SelectAtMostOne", "SelectAtLeastOne", "SelectAny"
	Order   string   `xml:"order,attr"` // "Explicit", "Alphabetical"
	Plugins []Plugin `xml:"plugins>plugin"`
}

// Plugin with comprehensive structure
type Plugin struct {
	Name           string          `xml:"name,attr"`
	Description    string          `xml:"description"`
	Image          string          `xml:"image"`
	Files          *FileList       `xml:"files"`
	TypeDescriptor *TypeDescriptor `xml:"typeDescriptor"`
	ConditionFlags *ConditionFlags `xml:"conditionFlags"`
}

// File list structure
type FileList struct {
	Files   []FileInstall   `xml:"file"`
	Folders []FolderInstall `xml:"folder"`
}

// File installation structure with all attributes
type FileInstall struct {
	Source          string `xml:"source,attr"`
	Destination     string `xml:"destination,attr"`
	AlwaysInstall   string `xml:"alwaysInstall,attr"`   // "true"/"false"
	InstallIfUsable string `xml:"installIfUsable,attr"` // "true"/"false"
	Priority        string `xml:"priority,attr"`
}

// Folder installation structure with all attributes
type FolderInstall struct {
	Source          string `xml:"source,attr"`
	Destination     string `xml:"destination,attr"`
	AlwaysInstall   string `xml:"alwaysInstall,attr"`
	InstallIfUsable string `xml:"installIfUsable,attr"`
	Priority        string `xml:"priority,attr"`
}

// Type descriptor for plugin behavior
type TypeDescriptor struct {
	Type           *PluginType     `xml:"type"`
	DependencyType *DependencyType `xml:"dependencyType"`
}

// Plugin type structure
type PluginType struct {
	Name string `xml:"name,attr"` // "Required", "Optional", "Recommended", "NotUsable", "CouldBeUsable"
}

// Dependency-based type with patterns
type DependencyType struct {
	DefaultType *PluginType         `xml:"defaultType"`
	Patterns    []DependencyPattern `xml:"patterns>pattern"`
}

// Dependency pattern structure
type DependencyPattern struct {
	Dependencies *CompositeDependency `xml:"dependencies"`
	Type         *PluginType          `xml:"type"`
}

// Condition flags set by plugins
type ConditionFlags struct {
	Flags []ConditionFlag `xml:"flag"`
}

// Individual condition flag
type ConditionFlag struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// Conditional file installation patterns
type ConditionalFileInstall struct {
	Dependencies *CompositeDependency `xml:"dependencies"`
	Files        *FileList            `xml:"files"`
}

// Composite dependency structure (recursive)
type CompositeDependency struct {
	Operator     string       `xml:"operator,attr"` // "And", "Or"
	Dependencies []Dependency `xml:",any"`
}

// Union type for all dependency types
type Dependency struct {
	XMLName xml.Name

	// File dependency
	File  string `xml:"file,attr"`
	State string `xml:"state,attr"` // "Active", "Inactive", "Missing"

	// Flag dependency
	Flag  string `xml:"flag,attr"`
	Value string `xml:"value,attr"`

	// Game version dependency
	Version string `xml:"version,attr"`

	// Nested composite dependency
	Operator     string       `xml:"operator,attr"`
	Dependencies []Dependency `xml:",any"`
}

// Legacy support structure for older FOMOD formats
type PluginFile struct {
	Source      string `xml:"source,attr"`
	Destination string `xml:"destination,attr"`
	Priority    string `xml:"priority,attr"`
}

// Enhanced parsing function with full pattern support
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
	case bytes.HasPrefix(raw, []byte{0xEF, 0xBB, 0xBF}):
		// UTF-8 BOM
		reader = bytes.NewReader(raw[3:])
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

// Helper methods for working with the parsed structure

// GetFileList returns files from both new and legacy structures
func (p *Plugin) GetFileList() []FileInstall {
	var files []FileInstall

	if p.Files != nil {
		files = append(files, p.Files.Files...)
	}

	return files
}

// GetFolderList returns folders from both new and legacy structures
func (p *Plugin) GetFolderList() []FolderInstall {
	var folders []FolderInstall

	if p.Files != nil {
		folders = append(folders, p.Files.Folders...)
	}

	return folders
}

// IsVisible checks if a step should be visible based on conditions
func (step *InstallStep) IsVisible(flags map[string]string, installedFiles map[string]bool) bool {
	if step.Visible == nil {
		return true // No visibility conditions means always visible
	}

	return evaluateCompositeDependency(step.Visible.Operator, step.Visible.Dependencies, flags, installedFiles)
}

// GetPluginType returns the effective plugin type based on dependencies
func (p *Plugin) GetPluginType(flags map[string]string, installedFiles map[string]bool) string {
	if p.TypeDescriptor == nil {
		return "Optional" // Default type
	}

	if p.TypeDescriptor.Type != nil {
		return p.TypeDescriptor.Type.Name
	}

	if p.TypeDescriptor.DependencyType != nil {
		// Check dependency patterns
		for _, pattern := range p.TypeDescriptor.DependencyType.Patterns {
			if pattern.Dependencies != nil &&
				evaluateCompositeDependency(pattern.Dependencies.Operator, pattern.Dependencies.Dependencies, flags, installedFiles) {
				return pattern.Type.Name
			}
		}

		// Return default type if no patterns match
		if p.TypeDescriptor.DependencyType.DefaultType != nil {
			return p.TypeDescriptor.DependencyType.DefaultType.Name
		}
	}

	return "Optional"
}

// GetPriority returns the priority value as integer
func (f *FileInstall) GetPriority() int {
	if f.Priority == "" {
		return 0
	}
	priority, _ := strconv.Atoi(f.Priority)
	return priority
}

// GetPriority returns the priority value as integer
func (f *FolderInstall) GetPriority() int {
	if f.Priority == "" {
		return 0
	}
	priority, _ := strconv.Atoi(f.Priority)
	return priority
}

// ShouldAlwaysInstall checks if file should always be installed
func (f *FileInstall) ShouldAlwaysInstall() bool {
	return f.AlwaysInstall == "true"
}

// ShouldAlwaysInstall checks if folder should always be installed
func (f *FolderInstall) ShouldAlwaysInstall() bool {
	return f.AlwaysInstall == "true"
}

// evaluateCompositeDependency evaluates dependency conditions recursively
func evaluateCompositeDependency(operator string, dependencies []Dependency, flags map[string]string, installedFiles map[string]bool) bool {
	if len(dependencies) == 0 {
		return true
	}

	results := make([]bool, len(dependencies))

	for i, dep := range dependencies {
		switch dep.XMLName.Local {
		case "fileDependency":
			if state, exists := installedFiles[dep.File]; exists {
				switch dep.State {
				case "Active":
					results[i] = state
				case "Inactive":
					results[i] = !state
				case "Missing":
					results[i] = false
				default:
					results[i] = state
				}
			} else {
				results[i] = dep.State == "Missing"
			}

		case "flagDependency":
			if value, exists := flags[dep.Flag]; exists {
				results[i] = value == dep.Value
			} else {
				results[i] = false
			}

		case "gameDependency":
			// Game version check would need to be implemented based on your needs
			results[i] = true // Placeholder

		case "dependencies":
			// Nested composite dependency
			results[i] = evaluateCompositeDependency(dep.Operator, dep.Dependencies, flags, installedFiles)

		default:
			results[i] = false
		}
	}

	// Apply operator logic
	switch operator {
	case "And":
		for _, result := range results {
			if !result {
				return false
			}
		}
		return true

	case "Or":
		for _, result := range results {
			if result {
				return true
			}
		}
		return false

	default:
		// Default to "And" behavior
		for _, result := range results {
			if !result {
				return false
			}
		}
		return true
	}
}
