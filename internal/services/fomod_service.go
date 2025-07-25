package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"scrolljack/internal/db/dtos"
	"scrolljack/internal/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// DetectFomodOptions detects which FOMOD options were selected based on installed files.
func DetectFomodOptions(ctx context.Context, db *sql.DB, modId string) (string, error) {
	// Open file dialog to select mod archive
	log.Println("Opening file dialog to select mod archive")
	result, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select a mod archive (zip, rar, 7z)",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Mod Archive",
				Pattern:     "*.zip;*.rar;*.7z",
			},
		},
	})
	if err != nil {
		log.Printf("Error opening file dialog: %v", err)
		return "", fmt.Errorf("failed to open file dialog: %w", err)
	}
	if result == "" {
		log.Println("File dialog canceled by user")
		return "", nil // User canceled the dialog
	}
	log.Printf("Selected archive: %s", result)

	// Get app directory and set up temp directory
	log.Println("Getting app directory")
	appDir, err := utils.GetAppDir()
	if err != nil {
		log.Printf("Error getting app directory: %v", err)
		return "", fmt.Errorf("failed to get app directory: %w", err)
	}
	tempDir := filepath.Join(appDir, "temp")
	log.Printf("Extracting archive to temporary directory: %s", tempDir)
	if err := utils.ExtractArchive(result, tempDir); err != nil {
		log.Printf("Error extracting archive: %v", err)
		return "", fmt.Errorf("failed to extract file: %w", err)
	}
	defer cleanupTempDir(tempDir) // Ensure cleanup after function exits

	// Log contents of temp directory
	log.Println("Listing contents of temporary directory")
	if err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing %s: %v", path, err)
			return nil
		}
		log.Printf("  %s (%s)", path, map[bool]string{true: "dir", false: "file"}[!info.IsDir()])
		return nil
	}); err != nil {
		log.Printf("Error walking temp directory %s: %v", tempDir, err)
	}

	// Get mod files from database
	log.Printf("Fetching mod files for mod ID: %s", modId)
	modFiles, err := GetModFilesByModId(ctx, db, modId)
	if err != nil {
		log.Printf("Error fetching mod files for mod ID %s: %v", modId, err)
		return "", fmt.Errorf("failed to get mod files for mod ID %s: %w", modId, err)
	}
	log.Printf("Found %d mod files in database", len(modFiles))

	// Find and parse FOMOD configuration
	log.Println("Searching for FOMOD directory and ModuleConfig.xml")
	fomodDir, moduleConfigPath, err := utils.FindFomodDirectory(tempDir)
	if err != nil {
		log.Printf("Error finding FOMOD directory: %v", err)
		return "", fmt.Errorf("failed to find FOMOD directory: %w", err)
	}
	if fomodDir == "" || moduleConfigPath == "" {
		log.Println("No FOMOD directory or ModuleConfig.xml found in archive")
		return "No FOMOD configuration found in this mod archive", nil
	}
	log.Printf("Found FOMOD config at: %s", moduleConfigPath)

	// Log contents of fomod directory
	log.Printf("Listing contents of FOMOD directory: %s", fomodDir)
	if err := filepath.Walk(fomodDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing %s: %v", path, err)
			return nil
		}
		log.Printf("  %s (%s)", path, map[bool]string{true: "dir", false: "file"}[!info.IsDir()])
		return nil
	}); err != nil {
		log.Printf("Error walking FOMOD directory %s: %v", fomodDir, err)
	}

	log.Println("Parsing FOMOD configuration")
	config, err := utils.ParseFomodConfig(moduleConfigPath)
	if err != nil {
		log.Printf("Error parsing FOMOD config: %v", err)
		return "", fmt.Errorf("failed to parse FOMOD module config: %w", err)
	}

	// Log FOMOD configuration
	configJSON, _ := json.MarshalIndent(config, "", "  ")
	log.Printf("FOMOD Configuration:\n%s", string(configJSON))

	// Detect selected plugins
	log.Println("Detecting selected FOMOD plugins")
	selectedPlugins, err := detectSelectedPlugins(config, modFiles, tempDir, fomodDir)
	if err != nil {
		log.Printf("Error detecting selected plugins: %v", err)
		return "", fmt.Errorf("failed to detect selected plugins: %w", err)
	}

	// Create comma-separated string of selected plugins
	resultStr := strings.Join(selectedPlugins, ",")
	if resultStr == "" {
		log.Println("No plugins matched the mod files")
	} else {
		log.Printf("Selected FOMOD plugins: %s", resultStr)
	}
	return resultStr, nil
}

// detectSelectedPlugins matches mod files to FOMOD plugins using paths and hashes.
func detectSelectedPlugins(config *utils.ModuleConfig, modFiles []dtos.ModFileDTO, tempDir, fomodDir string) ([]string, error) {
	var selectedPlugins []string
	modFileMap := make(map[string]dtos.ModFileDTO) // Map normalized path to ModFileDTO
	for _, mf := range modFiles {
		normalizedPath := strings.ReplaceAll(mf.Path, "\\", "/")
		modFileMap[normalizedPath] = mf
		log.Printf("Mod file: %s (Hash: %s)", normalizedPath, mf.Hash)
	}

	// Find all folders in tempDir for case-insensitive matching
	folderMap := make(map[string]string) // Map lowercase folder name to actual path
	if err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing %s: %v", path, err)
			return nil
		}
		if info.IsDir() {
			folderName := strings.ToLower(filepath.Base(path))
			folderMap[folderName] = path
			log.Printf("Found folder: %s", path)
		}
		return nil
	}); err != nil {
		log.Printf("Error walking temp directory for folders: %v", err)
	}

	for stepIdx, step := range config.InstallSteps {
		log.Printf("Processing Install Step %d", stepIdx+1)
		for groupIdx, group := range step.Groups {
			log.Printf("  Processing Group %d", groupIdx+1)
			// Track plugins that write to the same destination to resolve ambiguities
			destToPlugins := make(map[string][]struct {
				PluginName string
				SourcePath string
				Hash       string
			})
			var stepMatches []string

			for _, plugin := range group.Plugins {
				log.Printf("    Checking plugin: %s", plugin.Name)
				matches := false
				if plugin.Files != nil {
					for _, file := range plugin.Files {
						normalizedDest := strings.ReplaceAll(file.Destination, "\\", "/")
						normalizedSource := strings.ReplaceAll(file.Source, "\\", "/")
						sourcePath := filepath.Join(fomodDir, file.Source)
						log.Printf("      File: %s -> %s", normalizedSource, normalizedDest)

						// Check if destination path exists in mod files
						if modFile, exists := modFileMap[normalizedDest]; exists {
							// Compute hash of source file in archive
							hash, err := utils.HashFile(sourcePath)
							if err != nil {
								log.Printf("      Error hashing source file %s: %v", sourcePath, err)
								continue
							}
							log.Printf("      Source file hash: %s, Mod file hash: %s", hash, modFile.Hash)

							// Compare hashes
							if hash == modFile.Hash {
								log.Printf("      Match found for %s (hash match)", normalizedDest)
								matches = true
							} else {
								log.Printf("      Hash mismatch for %s", normalizedDest)
							}

							// Track plugins that write to this destination
							destToPlugins[normalizedDest] = append(destToPlugins[normalizedDest], struct {
								PluginName string
								SourcePath string
								Hash       string
							}{PluginName: plugin.Name, SourcePath: normalizedSource, Hash: hash})
						} else {
							log.Printf("      No mod file found for destination %s", normalizedDest)
						}
					}
				}
				if plugin.Folders != nil && !matches {
					for _, folder := range plugin.Folders {
						normalizedDest := strings.ReplaceAll(folder.Destination, "\\", "/")
						normalizedSource := strings.ReplaceAll(folder.Source, "\\", "/")
						if normalizedDest == "" {
							normalizedDest = normalizedSource
						}
						log.Printf("      Folder: %s -> %s", normalizedSource, normalizedDest)

						// Try case-insensitive folder lookup
						sourceFolderPath := filepath.Join(fomodDir, folder.Source)
						lowerSource := strings.ToLower(normalizedSource)
						if mappedPath, exists := folderMap[lowerSource]; exists && mappedPath != sourceFolderPath {
							log.Printf("      Case-insensitive match for folder %s: %s", normalizedSource, mappedPath)
							sourceFolderPath = mappedPath
						}

						// Walk the source folder to find files and compare with mod files
						err := filepath.Walk(sourceFolderPath, func(path string, info os.FileInfo, err error) error {
							if err != nil {
								log.Printf("      Error walking folder %s: %v", sourceFolderPath, err)
								return nil // Continue walking despite errors
							}
							if info.IsDir() {
								return nil // Skip directories
							}
							relativePath, err := filepath.Rel(fomodDir, path)
							if err != nil {
								log.Printf("      Error computing relative path for %s: %v", path, err)
								return nil
							}
							normalizedRelativePath := strings.ReplaceAll(relativePath, "\\", "/")
							// Compute expected destination path
							destPath := normalizedRelativePath
							if normalizedDest != normalizedSource {
								if strings.HasPrefix(normalizedRelativePath, normalizedSource) {
									destPath = filepath.Join(normalizedDest, strings.TrimPrefix(normalizedRelativePath, normalizedSource))
									destPath = strings.ReplaceAll(destPath, "\\", "/")
								}
							}
							log.Printf("      Checking file in folder: %s -> %s", normalizedRelativePath, destPath)

							// Check if destination path exists in mod files
							if modFile, exists := modFileMap[destPath]; exists {
								// Compute hash of source file
								hash, err := utils.HashFile(path)
								if err != nil {
									log.Printf("      Error hashing file %s: %v", path, err)
									return nil
								}
								log.Printf("      File %s hash: %s, Mod file hash: %s", destPath, hash, modFile.Hash)
								if hash == modFile.Hash {
									log.Printf("      Match found for %s (hash match)", destPath)
									matches = true
								} else {
									log.Printf("      Hash mismatch for %s", destPath)
								}
							}
							return nil
						})
						if err != nil {
							log.Printf("      Error walking folder %s: %v", sourceFolderPath, err)
						}
						if matches {
							log.Printf("      Match found for folder %s", normalizedDest)
							break
						}
					}
				}
				if matches {
					log.Printf("    Plugin %s selected", plugin.Name)
					stepMatches = append(stepMatches, plugin.Name)
				}
			}

			// Resolve ambiguities for destinations with multiple plugins
			for dest, plugins := range destToPlugins {
				if len(plugins) > 1 {
					log.Printf("    Ambiguity detected for destination %s: %v", dest, plugins)
					// Find the plugin with matching hash
					if modFile, exists := modFileMap[dest]; exists {
						for _, p := range plugins {
							if p.Hash == modFile.Hash {
								log.Printf("    Resolved ambiguity: %s selected for %s (hash match)", p.PluginName, dest)
								if !contains(stepMatches, p.PluginName) {
									stepMatches = append(stepMatches, p.PluginName)
								}
							}
						}
					}
				} else if len(plugins) == 1 && !contains(stepMatches, plugins[0].PluginName) {
					// Single plugin for destination, ensure it's included if hash matches
					if modFile, exists := modFileMap[dest]; exists && plugins[0].Hash == modFile.Hash {
						log.Printf("    Single plugin %s confirmed for %s (hash match)", plugins[0].PluginName, dest)
						stepMatches = append(stepMatches, plugins[0].PluginName)
					}
				}
			}

			// Add unique matches from this step
			for _, pluginName := range stepMatches {
				if !contains(selectedPlugins, pluginName) {
					selectedPlugins = append(selectedPlugins, pluginName)
				}
			}
		}
	}

	// Fallback: Try direct file matching across entire tempDir
	if len(selectedPlugins) == 0 {
		log.Println("No plugins matched via FOMOD config, attempting direct file matching")
		for _, modFile := range modFiles {
			normalizedModPath := strings.ReplaceAll(modFile.Path, "\\", "/")
			err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Printf("Error walking tempDir for %s: %v", path, err)
					return nil
				}
				if info.IsDir() {
					return nil
				}
				hash, err := utils.HashFile(path)
				if err != nil {
					log.Printf("Error hashing file %s: %v", path, err)
					return nil
				}
				if hash == modFile.Hash {
					relativePath, err := filepath.Rel(tempDir, path)
					if err != nil {
						log.Printf("Error computing relative path for %s: %v", path, err)
						return nil
					}
					normalizedRelativePath := strings.ReplaceAll(relativePath, "\\", "/")
					log.Printf("Direct match found: %s (hash: %s) matches mod file %s", normalizedRelativePath, hash, normalizedModPath)
					// Try to find plugin with matching source path
					for _, step := range config.InstallSteps {
						for _, group := range step.Groups {
							for _, plugin := range group.Plugins {
								for _, file := range plugin.Files {
									if strings.ReplaceAll(file.Source, "\\", "/") == normalizedRelativePath {
										log.Printf("    Plugin %s matched via direct file %s", plugin.Name, normalizedRelativePath)
										if !contains(selectedPlugins, plugin.Name) {
											selectedPlugins = append(selectedPlugins, plugin.Name)
										}
										return nil
									}
								}
								for _, folder := range plugin.Folders {
									normalizedSource := strings.ReplaceAll(folder.Source, "\\", "/")
									if strings.HasPrefix(normalizedRelativePath, normalizedSource) {
										log.Printf("    Plugin %s matched via folder %s for file %s", plugin.Name, normalizedSource, normalizedRelativePath)
										if !contains(selectedPlugins, plugin.Name) {
											selectedPlugins = append(selectedPlugins, plugin.Name)
										}
										return nil
									}
								}
							}
						}
					}
				}
				return nil
			})
			if err != nil {
				log.Printf("Error walking tempDir for direct matching: %v", err)
			}
		}
	}

	if len(selectedPlugins) == 0 {
		log.Println("No plugins matched the mod files after direct matching")
	}
	return selectedPlugins, nil
}

// contains checks if a string exists in a slice.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// cleanupTempDir removes the temporary directory.
func cleanupTempDir(tempDir string) {
	log.Printf("Cleaning up temporary directory: %s", tempDir)
	if err := os.RemoveAll(tempDir); err != nil {
		log.Printf("Failed to delete temporary directory %s: %v", tempDir, err)
	}
}
