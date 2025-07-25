package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"scrolljack/internal/db/dtos"
	"scrolljack/internal/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ArchiveFile represents a file in the extracted archive (from original code)
type ArchiveFile struct {
	RelativePath string
	FullPath     string
	Hash         string
	Size         int64
}

// Enhanced structures for complex FOMOD detection
type EnhancedDetectionResult struct {
	StepIndex            int
	StepName             string
	IsVisible            bool
	BestPlugin           string
	Confidence           float64
	MatchDetails         string
	AlternativePlugins   []PluginMatch
	RequiredDependencies []string
	ConflictingChoices   []string
	GroupType            string
}

type PluginMatch struct {
	Name       string
	Confidence float64
	Reason     string
	Type       string // Required, Optional, Recommended, etc.
}

type FomodState struct {
	Flags           map[string]string
	SelectedPlugins map[int]string // step index -> plugin name
	InstalledFiles  map[string]bool
}

type EnhancedFomodAnalysis struct {
	ModuleName             string
	TotalSteps             int
	StepResults            []EnhancedDetectionResult
	SelectedCount          int
	OverallSuccess         bool
	ConditionalFileMatches int
	RequiredFileMatches    int
	DetectionQuality       string // "High", "Medium", "Low"
	RecommendedChoices     []string
	PotentialConflicts     []string
	MissingDependencies    []string
}

// Enhanced detection function with complex case handling
func EnhancedDetectFomodOptions(ctx context.Context, db *sql.DB, modId string) (string, error) {
	// File dialog and extraction (same as before)
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
		return "", fmt.Errorf("failed to open file dialog: %w", err)
	}
	if result == "" {
		return "", nil
	}

	appDir, err := utils.GetAppDir()
	if err != nil {
		return "", fmt.Errorf("failed to get app directory: %w", err)
	}
	tempDir := filepath.Join(appDir, "temp")

	log.Printf("Extracting archive: %s", filepath.Base(result))
	if err := utils.ExtractArchive(result, tempDir); err != nil {
		return "", fmt.Errorf("failed to extract file: %w", err)
	}
	defer cleanupTempDir(tempDir)

	// Get mod files from database
	modFiles, err := GetModFilesByModId(ctx, db, modId)
	if err != nil {
		return "", fmt.Errorf("failed to get mod files: %w", err)
	}
	log.Printf("Found %d installed mod files", len(modFiles))

	// Find and parse FOMOD configuration
	fomodDir, moduleConfigPath, err := utils.FindFomodDirectory(tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to find FOMOD directory: %w", err)
	}
	if fomodDir == "" || moduleConfigPath == "" {
		return "No FOMOD configuration found", nil
	}
	log.Printf("Found FOMOD config at: %s", moduleConfigPath)

	// Parse FOMOD configuration
	config, err := utils.ParseFomodConfig(moduleConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse FOMOD config: %w", err)
	}

	// Build enhanced file maps
	archiveFiles, err := buildArchiveFileMap(tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to build archive file map: %w", err)
	}
	modFileMap := buildModFileMap(modFiles)

	// Perform enhanced detection with complex case handling
	analysis := performEnhancedFomodDetection(config, archiveFiles, modFileMap, fomodDir)

	return formatEnhancedDetectionResults(analysis), nil
}

// Enhanced detection with complex case handling
func performEnhancedFomodDetection(config *utils.ModuleConfig, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO, fomodDir string) *EnhancedFomodAnalysis {
	analysis := &EnhancedFomodAnalysis{
		ModuleName:  config.ModuleName,
		TotalSteps:  len(config.InstallSteps),
		StepResults: make([]EnhancedDetectionResult, 0),
	}

	log.Printf("🔍 Enhanced FOMOD Analysis: %s with %d steps", analysis.ModuleName, analysis.TotalSteps)

	// Initialize state tracking
	state := &FomodState{
		Flags:           make(map[string]string),
		SelectedPlugins: make(map[int]string),
		InstalledFiles:  buildInstalledFileMap(modFileMap),
	}

	// First pass: Analyze required files
	if config.RequiredInstallFiles != nil {
		analysis.RequiredFileMatches = analyzeRequiredFiles(config.RequiredInstallFiles, archiveFiles, modFileMap)
		log.Printf("📋 Required files analysis: %d matches found", analysis.RequiredFileMatches)
	}

	// Second pass: Multi-step analysis with dependency tracking
	for stepIdx, step := range config.InstallSteps {
		log.Printf("🎯 Processing step %d: %s", stepIdx+1, step.Name)

		result := analyzeEnhancedInstallStep(stepIdx, step, archiveFiles, modFileMap, state, fomodDir)
		analysis.StepResults = append(analysis.StepResults, result)

		// Update state based on best choice
		if result.BestPlugin != "" {
			analysis.SelectedCount++
			state.SelectedPlugins[stepIdx] = result.BestPlugin

			// Update flags for selected plugin
			updateStateFlags(stepIdx, step, result.BestPlugin, state)
		}
	}

	// Third pass: Analyze conditional file installs
	if len(config.ConditionalFileInstalls) > 0 {
		analysis.ConditionalFileMatches = analyzeConditionalFileInstalls(config.ConditionalFileInstalls, archiveFiles, modFileMap, state)
		log.Printf("🔄 Conditional files analysis: %d matches found", analysis.ConditionalFileMatches)
	}

	// Fourth pass: Quality assessment and recommendations
	analysis.DetectionQuality = assessDetectionQuality(analysis)
	analysis.RecommendedChoices = generateRecommendations(analysis, config, state)
	analysis.PotentialConflicts = detectPotentialConflicts(analysis, config, state)
	analysis.MissingDependencies = detectMissingDependencies(config, state)

	analysis.OverallSuccess = analysis.SelectedCount > 0 && analysis.DetectionQuality != "Low"

	log.Printf("✅ Enhanced detection complete: %s quality, %d/%d steps detected",
		analysis.DetectionQuality, analysis.SelectedCount, analysis.TotalSteps)

	return analysis
}

// Enhanced step analysis with alternative plugin detection
func analyzeEnhancedInstallStep(stepIdx int, step utils.InstallStep, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO, state *FomodState, fomodDir string) EnhancedDetectionResult {
	result := EnhancedDetectionResult{
		StepIndex:          stepIdx,
		StepName:           step.Name,
		IsVisible:          step.IsVisible(state.Flags, state.InstalledFiles),
		AlternativePlugins: make([]PluginMatch, 0),
	}

	log.Printf("🔍 Enhanced Step Analysis %d: '%s'", stepIdx+1, step.Name)
	log.Printf("   Visible: %t, Groups: %d", result.IsVisible, len(step.Groups))

	if !result.IsVisible {
		result.MatchDetails = "Step not visible based on current conditions"
		return result
	}

	var allMatches []PluginMatch
	var groupType string

	// Analyze each plugin group with enhanced scoring
	for groupIdx, group := range step.Groups {
		groupType = group.Type
		result.GroupType = group.Type

		log.Printf("   📦 Group %d (%s): %d plugins, type: %s", groupIdx+1, group.Name, len(group.Plugins), group.Type)

		for pluginIdx, plugin := range group.Plugins {
			log.Printf("     🔌 Plugin %d: '%s'", pluginIdx+1, plugin.Name)

			confidence, details, pluginType := analyzeEnhancedPlugin(plugin, archiveFiles, modFileMap, state, fomodDir)

			match := PluginMatch{
				Name:       plugin.Name,
				Confidence: confidence,
				Reason:     strings.Join(details, "; "),
				Type:       pluginType,
			}

			allMatches = append(allMatches, match)
			log.Printf("     📊 Confidence: %.1f%% (%s, type: %s)", confidence*100, match.Reason, pluginType)
		}
	}

	// Sort matches by confidence and apply group type logic
	sort.Slice(allMatches, func(i, j int) bool {
		return allMatches[i].Confidence > allMatches[j].Confidence
	})

	// Apply group-specific selection logic
	selectedMatches := applyGroupSelectionLogic(allMatches, groupType)

	if len(selectedMatches) > 0 {
		result.BestPlugin = selectedMatches[0].Name
		result.Confidence = selectedMatches[0].Confidence
		result.MatchDetails = selectedMatches[0].Reason

		// Store alternatives
		if len(selectedMatches) > 1 {
			result.AlternativePlugins = selectedMatches[1:]
		}
	}

	// Detect dependencies and conflicts
	result.RequiredDependencies = detectRequiredDependencies(step, state)
	result.ConflictingChoices = detectConflictingChoices(step, state)

	return result
}

// Enhanced plugin analysis with comprehensive scoring
func analyzeEnhancedPlugin(plugin utils.Plugin, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO, state *FomodState, fomodDir string) (float64, []string, string) {
	var details []string
	totalFiles := 0
	matchedFiles := 0
	perfectMatches := 0

	log.Printf("       🔍 Enhanced plugin analysis: %s", plugin.Name)

	// Get plugin type with current state
	pluginType := plugin.GetPluginType(state.Flags, state.InstalledFiles)
	log.Printf("         Type: %s", pluginType)

	if pluginType == "NotUsable" {
		return 0.0, []string{"Plugin marked as not usable"}, pluginType
	}

	// Analyze individual files with enhanced scoring
	fileList := plugin.GetFileList()
	for _, file := range fileList {
		totalFiles++
		isMatch, isPerfect := enhancedMatchFile(file.Source, file.Destination, archiveFiles, modFileMap, &details)
		if isMatch {
			matchedFiles++
			if isPerfect {
				perfectMatches++
			}
		}
	}

	// Analyze folders with enhanced scoring
	folderList := plugin.GetFolderList()
	for _, folder := range folderList {
		folderMatches, folderTotal, folderPerfect := analyzeEnhancedFolderContents(folder.Source, folder.Destination, archiveFiles, modFileMap, fomodDir)
		matchedFiles += folderMatches
		totalFiles += folderTotal
		perfectMatches += folderPerfect

		if folderMatches > 0 {
			details = append(details, fmt.Sprintf("Folder '%s': %d/%d files (%d perfect)", folder.Source, folderMatches, folderTotal, folderPerfect))
		}
	}

	if totalFiles == 0 {
		// Plugin with no files - check for condition flags only
		if plugin.ConditionFlags != nil && len(plugin.ConditionFlags.Flags) > 0 {
			return 0.5, []string{"Plugin sets flags but has no files"}, pluginType
		}
		return 0.0, []string{"No files to analyze"}, pluginType
	}

	// Calculate enhanced confidence score
	baseConfidence := float64(matchedFiles) / float64(totalFiles)
	perfectRatio := float64(perfectMatches) / float64(totalFiles)

	// Enhanced scoring factors
	confidence := baseConfidence

	// Boost for perfect matches
	if perfectRatio > 0 {
		confidence += perfectRatio * 0.2 // Up to 20% boost for perfect matches
	}

	// Plugin type modifiers
	switch pluginType {
	case "Required":
		if confidence > 0 {
			confidence = confidence * 1.3 // 30% boost for required plugins
		}
	case "Recommended":
		if confidence > 0 {
			confidence = confidence * 1.1 // 10% boost for recommended plugins
		}
	case "CouldBeUsable":
		confidence = confidence * 0.8 // 20% penalty for conditional plugins
	}

	// File priority consideration
	if hasHighPriorityFiles(plugin) {
		confidence = confidence * 1.05 // 5% boost for high priority files
	}

	// Ensure confidence doesn't exceed 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	log.Printf("         🎯 Enhanced confidence: %.1f%% (base: %.1f%%, perfect: %.1f%%, type: %s)",
		confidence*100, baseConfidence*100, perfectRatio*100, pluginType)

	return confidence, details, pluginType
}

// Enhanced file matching with perfect match detection
func enhancedMatchFile(sourcePath, destPath string, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO, details *[]string) (bool, bool) {
	normalizedSource := strings.ReplaceAll(sourcePath, "\\", "/")
	normalizedDest := strings.ReplaceAll(destPath, "\\", "/")

	// Find source file with enhanced matching
	archiveFile, sourceFound, exactPath := findArchiveFileEnhanced(normalizedSource, archiveFiles)
	if !sourceFound {
		return false, false
	}

	// Find destination file with enhanced matching
	modFile, destFound, exactDestPath := findModFileEnhanced(normalizedDest, modFileMap)
	if !destFound {
		return false, false
	}

	// Check hash match
	if archiveFile.Hash == modFile.Hash {
		isPerfect := exactPath && exactDestPath
		matchType := "fuzzy"
		if isPerfect {
			matchType = "perfect"
		}

		*details = append(*details, fmt.Sprintf("File match (%s): %s -> %s", matchType, archiveFile.RelativePath, normalizedDest))
		return true, isPerfect
	}

	return false, false
}

// Enhanced folder analysis with perfect match tracking
func analyzeEnhancedFolderContents(sourcePath, destPath string, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO, fomodDir string) (int, int, int) {
	matches := 0
	total := 0
	perfectMatches := 0

	normalizedSource := strings.ReplaceAll(sourcePath, "\\", "/")
	normalizedDest := strings.ReplaceAll(destPath, "\\", "/")

	for archivePath, archiveFile := range archiveFiles {
		if !isFileInFolder(archivePath, normalizedSource) {
			continue
		}

		total++
		destFilePath := calculateDestinationPath(archivePath, normalizedSource, normalizedDest)

		modFile, found, exactMatch := findModFileEnhanced(destFilePath, modFileMap)
		if found && archiveFile.Hash == modFile.Hash {
			matches++
			if exactMatch && strings.HasPrefix(archivePath, normalizedSource+"/") {
				perfectMatches++
			}
		}
	}

	return matches, total, perfectMatches
}

// Enhanced archive file finding with exact path tracking
func findArchiveFileEnhanced(sourcePath string, archiveFiles map[string]ArchiveFile) (ArchiveFile, bool, bool) {
	// Try exact match first
	if af, exists := archiveFiles[sourcePath]; exists {
		return af, true, true
	}

	// Try suffix matching
	for archivePath, af := range archiveFiles {
		if strings.HasSuffix(archivePath, sourcePath) || strings.HasSuffix(archivePath, "/"+sourcePath) {
			return af, true, false
		}
	}

	return ArchiveFile{}, false, false
}

// Enhanced mod file finding with exact path tracking
func findModFileEnhanced(destPath string, modFileMap map[string]dtos.ModFileDTO) (dtos.ModFileDTO, bool, bool) {
	// Try exact match
	if mf, exists := modFileMap[destPath]; exists {
		return mf, true, true
	}

	// Try case-insensitive match
	if mf, exists := modFileMap[strings.ToLower(destPath)]; exists {
		return mf, true, false
	}

	return dtos.ModFileDTO{}, false, false
}

// Apply group selection logic based on group type
func applyGroupSelectionLogic(matches []PluginMatch, groupType string) []PluginMatch {
	if len(matches) == 0 {
		return matches
	}

	switch groupType {
	case "SelectExactlyOne":
		// Return only the best match
		if matches[0].Confidence > 0 {
			return matches[:1]
		}
		return []PluginMatch{}

	case "SelectAll":
		// Return all matches above threshold
		var selected []PluginMatch
		for _, match := range matches {
			if match.Confidence > 0.3 { // Threshold for "all" selection
				selected = append(selected, match)
			}
		}
		return selected

	case "SelectAtMostOne":
		// Return best match if confidence is high enough
		if matches[0].Confidence > 0.5 {
			return matches[:1]
		}
		return []PluginMatch{}

	case "SelectAtLeastOne":
		// Return best match, or multiple if they're close
		var selected []PluginMatch
		bestConfidence := matches[0].Confidence
		for _, match := range matches {
			if match.Confidence >= bestConfidence*0.9 { // Within 90% of best
				selected = append(selected, match)
			}
		}
		return selected

	case "SelectAny":
		// Return all matches above threshold
		var selected []PluginMatch
		for _, match := range matches {
			if match.Confidence > 0.4 {
				selected = append(selected, match)
			}
		}
		return selected

	default:
		return matches[:1] // Default to single selection
	}
}

// Additional helper functions for complex case handling

func analyzeRequiredFiles(requiredFiles *utils.RequiredInstallFiles, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO) int {
	matches := 0

	for _, file := range requiredFiles.Files {
		if isMatch, _ := enhancedMatchFile(file.Source, file.Destination, archiveFiles, modFileMap, &[]string{}); isMatch {
			matches++
		}
	}

	for _, folder := range requiredFiles.Folders {
		folderMatches, _, _ := analyzeEnhancedFolderContents(folder.Source, folder.Destination, archiveFiles, modFileMap, "")
		matches += folderMatches
	}

	return matches
}

func updateStateFlags(stepIdx int, step utils.InstallStep, selectedPlugin string, state *FomodState) {
	// Find the selected plugin and update flags
	for _, group := range step.Groups {
		for _, plugin := range group.Plugins {
			if plugin.Name == selectedPlugin && plugin.ConditionFlags != nil {
				for _, flag := range plugin.ConditionFlags.Flags {
					state.Flags[flag.Name] = flag.Value
					log.Printf("🏁 Set flag: %s = %s", flag.Name, flag.Value)
				}
			}
		}
	}
}

func analyzeConditionalFileInstalls(conditionalInstalls []utils.ConditionalFileInstall, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO, state *FomodState) int {
	matches := 0

	for _, install := range conditionalInstalls {
		if install.Dependencies != nil {
			// Check if conditions are met
			if evaluateCompositeDependency(install.Dependencies.Operator, install.Dependencies.Dependencies, state.Flags, state.InstalledFiles) {
				// Count matching files
				if install.Files != nil {
					for _, file := range install.Files.Files {
						if isMatch, _ := enhancedMatchFile(file.Source, file.Destination, archiveFiles, modFileMap, &[]string{}); isMatch {
							matches++
						}
					}
				}
			}
		}
	}

	return matches
}

func assessDetectionQuality(analysis *EnhancedFomodAnalysis) string {
	if analysis.TotalSteps == 0 {
		return "Low"
	}

	detectionRate := float64(analysis.SelectedCount) / float64(analysis.TotalSteps)

	// Calculate average confidence
	totalConfidence := 0.0
	for _, result := range analysis.StepResults {
		totalConfidence += result.Confidence
	}
	avgConfidence := totalConfidence / float64(len(analysis.StepResults))

	if detectionRate >= 0.8 && avgConfidence >= 0.7 {
		return "High"
	} else if detectionRate >= 0.5 && avgConfidence >= 0.5 {
		return "Medium"
	}
	return "Low"
}

func generateRecommendations(analysis *EnhancedFomodAnalysis, config *utils.ModuleConfig, state *FomodState) []string {
	var recommendations []string

	for _, result := range analysis.StepResults {
		if result.BestPlugin != "" && result.Confidence >= 0.7 {
			recommendations = append(recommendations, fmt.Sprintf("%s: %s (%.0f%% confidence)",
				result.StepName, result.BestPlugin, result.Confidence*100))
		}
	}

	return recommendations
}

func detectPotentialConflicts(analysis *EnhancedFomodAnalysis, config *utils.ModuleConfig, state *FomodState) []string {
	var conflicts []string

	// Check for conflicting selections in SelectExactlyOne groups
	for _, result := range analysis.StepResults {
		if result.GroupType == "SelectExactlyOne" && len(result.AlternativePlugins) > 0 {
			for _, alt := range result.AlternativePlugins {
				if alt.Confidence >= result.Confidence*0.9 {
					conflicts = append(conflicts, fmt.Sprintf("%s: %s vs %s (similar confidence)",
						result.StepName, result.BestPlugin, alt.Name))
				}
			}
		}
	}

	return conflicts
}

func detectMissingDependencies(config *utils.ModuleConfig, state *FomodState) []string {
	var missing []string

	if config.ModuleDependencies != nil {
		// Check module-level dependencies
		if !evaluateCompositeDependency(config.ModuleDependencies.Operator,
			config.ModuleDependencies.Dependencies, state.Flags, state.InstalledFiles) {
			missing = append(missing, "Module dependencies not satisfied")
		}
	}

	return missing
}

func detectRequiredDependencies(step utils.InstallStep, state *FomodState) []string {
	var required []string

	for _, group := range step.Groups {
		for _, plugin := range group.Plugins {
			pluginType := plugin.GetPluginType(state.Flags, state.InstalledFiles)
			if pluginType == "Required" {
				required = append(required, plugin.Name)
			}
		}
	}

	return required
}

func detectConflictingChoices(step utils.InstallStep, state *FomodState) []string {
	var conflicts []string

	// This would need more sophisticated logic based on your specific conflict detection needs
	// For now, just return empty slice

	return conflicts
}

func hasHighPriorityFiles(plugin utils.Plugin) bool {
	if plugin.Files == nil {
		return false
	}

	for _, file := range plugin.Files.Files {
		if file.GetPriority() > 5 {
			return true
		}
	}

	for _, folder := range plugin.Files.Folders {
		if folder.GetPriority() > 5 {
			return true
		}
	}

	return false
}

// Helper function to evaluate composite dependencies (you already have this, but including for completeness)
func evaluateCompositeDependency(operator string, dependencies []utils.Dependency, flags map[string]string, installedFiles map[string]bool) bool {
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

		case "dependencies":
			results[i] = evaluateCompositeDependency(dep.Operator, dep.Dependencies, flags, installedFiles)

		default:
			results[i] = false
		}
	}

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
		for _, result := range results {
			if !result {
				return false
			}
		}
		return true
	}
}

// Complex scenario handlers for real-world FOМODs

// Handle cascading dependencies (when step 2 depends on step 1 selection)
func handleCascadingDependencies(config *utils.ModuleConfig, state *FomodState) {
	log.Printf("🔗 Analyzing cascading dependencies...")

	for stepIdx, step := range config.InstallSteps {
		// Re-evaluate visibility after each step selection
		if selectedPlugin, exists := state.SelectedPlugins[stepIdx]; exists {
			log.Printf("   Step %d selected: %s", stepIdx+1, selectedPlugin)

			// Update flags for next steps
			updateStateFlags(stepIdx, step, selectedPlugin, state)

			// Re-evaluate subsequent step visibility
			for nextStepIdx := stepIdx + 1; nextStepIdx < len(config.InstallSteps); nextStepIdx++ {
				nextStep := config.InstallSteps[nextStepIdx]
				wasVisible := nextStep.IsVisible(state.Flags, state.InstalledFiles)
				log.Printf("   Step %d (%s) visibility: %t", nextStepIdx+1, nextStep.Name, wasVisible)
			}
		}
	}
}

// Handle version-specific plugin selection
func handleVersionSpecificPlugins(plugins []utils.Plugin, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO) []PluginMatch {
	var versionMatches []PluginMatch

	log.Printf("🔖 Analyzing version-specific plugins...")

	for _, plugin := range plugins {
		// Check for version indicators in plugin names
		versionScore := calculateVersionScore(plugin.Name, archiveFiles, modFileMap)

		if versionScore > 0 {
			match := PluginMatch{
				Name:       plugin.Name,
				Confidence: versionScore,
				Reason:     "Version-specific match detected",
				Type:       plugin.GetPluginType(make(map[string]string), make(map[string]bool)),
			}
			versionMatches = append(versionMatches, match)
			log.Printf("   Version match: %s (%.1f%%)", plugin.Name, versionScore*100)
		}
	}

	return versionMatches
}

// Calculate version-specific confidence based on file patterns
func calculateVersionScore(pluginName string, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO) float64 {
	// Look for version indicators in plugin name
	versionKeywords := []string{"v1", "v2", "v3", "1.0", "2.0", "3.0", "old", "new", "legacy", "updated", "SE", "LE", "AE"}

	pluginLower := strings.ToLower(pluginName)
	versionFound := false

	for _, keyword := range versionKeywords {
		if strings.Contains(pluginLower, keyword) {
			versionFound = true
			log.Printf("     Found version keyword: %s in %s", keyword, pluginName)
			break
		}
	}

	if !versionFound {
		return 0.0
	}

	// Additional scoring based on file analysis would go here
	return 0.6 // Base score for version-specific plugins
}

// Handle texture/mesh variant detection
func handleVariantDetection(plugin utils.Plugin, archiveFiles map[string]ArchiveFile, modFileMap map[string]dtos.ModFileDTO) (float64, []string) {
	var details []string
	variantScore := 0.0

	log.Printf("🎨 Analyzing variants for plugin: %s", plugin.Name)

	// Common variant patterns
	variantPatterns := map[string][]string{
		"textures": {"textures/", ".dds", ".tga", ".png"},
		"meshes":   {"meshes/", ".nif", ".obj"},
		"sounds":   {"sounds/", ".wav", ".mp3", ".ogg"},
		"scripts":  {"scripts/", ".psc", ".pex"},
	}

	fileList := plugin.GetFileList()
	folderList := plugin.GetFolderList()

	// Analyze file variants
	for _, file := range fileList {
		for variantType, patterns := range variantPatterns {
			for _, pattern := range patterns {
				if strings.Contains(strings.ToLower(file.Source), pattern) {
					variantScore += 0.1
					details = append(details, fmt.Sprintf("%s variant detected in %s", variantType, file.Source))
					log.Printf("     Variant found: %s in %s", variantType, file.Source)
					break
				}
			}
		}
	}

	// Analyze folder variants
	for _, folder := range folderList {
		for variantType, patterns := range variantPatterns {
			for _, pattern := range patterns {
				if strings.Contains(strings.ToLower(folder.Source), pattern) {
					variantScore += 0.2
					details = append(details, fmt.Sprintf("%s variant folder: %s", variantType, folder.Source))
					log.Printf("     Variant folder: %s in %s", variantType, folder.Source)
					break
				}
			}
		}
	}

	// Cap variant score
	if variantScore > 1.0 {
		variantScore = 1.0
	}

	return variantScore, details
}

// Handle DLC/expansion-specific detection
func handleDLCSpecificDetection(plugin utils.Plugin, state *FomodState) float64 {
	dlcScore := 0.0
	pluginLower := strings.ToLower(plugin.Name)

	log.Printf("🎮 Checking DLC compatibility for: %s", plugin.Name)

	// Common DLC/expansion keywords
	dlcKeywords := map[string]float64{
		"dawnguard": 0.8, "dg": 0.8,
		"dragonborn": 0.8, "db": 0.8,
		"hearthfire": 0.8, "hf": 0.8,
		"anniversary": 0.9, "ae": 0.9,
		"special": 0.7, "se": 0.7,
		"legendary": 0.7, "le": 0.7,
		"expansion": 0.6, "dlc": 0.6,
	}

	for keyword, score := range dlcKeywords {
		if strings.Contains(pluginLower, keyword) {
			dlcScore = score
			log.Printf("     DLC keyword found: %s (score: %.1f)", keyword, score)
			break
		}
	}

	// Check for DLC-specific file patterns in plugin files
	if dlcScore > 0 {
		fileList := plugin.GetFileList()
		for _, file := range fileList {
			fileLower := strings.ToLower(file.Source)
			if strings.Contains(fileLower, "dawnguard") ||
				strings.Contains(fileLower, "dragonborn") ||
				strings.Contains(fileLower, "hearthfire") {
				dlcScore += 0.1
				log.Printf("     DLC-specific file found: %s", file.Source)
			}
		}
	}

	if dlcScore > 1.0 {
		dlcScore = 1.0
	}

	return dlcScore
}

// Handle compatibility patch detection
func handleCompatibilityPatchDetection(step utils.InstallStep, state *FomodState) []string {
	var patches []string

	log.Printf("🔧 Analyzing compatibility patches for step: %s", step.Name)

	stepLower := strings.ToLower(step.Name)
	patchKeywords := []string{"patch", "compatibility", "compat", "fix", "hotfix", "update"}

	isPatchStep := false
	for _, keyword := range patchKeywords {
		if strings.Contains(stepLower, keyword) {
			isPatchStep = true
			log.Printf("     Patch step detected: %s", keyword)
			break
		}
	}

	if isPatchStep {
		for _, group := range step.Groups {
			for _, plugin := range group.Plugins {
				pluginLower := strings.ToLower(plugin.Name)

				// Check for mod-specific patches
				modKeywords := []string{"enb", "skse", "ussep", "unofficial", "weather", "lighting"}
				for _, modKeyword := range modKeywords {
					if strings.Contains(pluginLower, modKeyword) {
						patches = append(patches, fmt.Sprintf("%s patch: %s", modKeyword, plugin.Name))
						log.Printf("     Compatibility patch found: %s for %s", plugin.Name, modKeyword)
					}
				}
			}
		}
	}

	return patches
}

// Handle performance variant detection (High/Medium/Low quality options)
func handlePerformanceVariants(plugins []utils.Plugin) map[string]PluginMatch {
	performanceMap := make(map[string]PluginMatch)

	log.Printf("⚡ Analyzing performance variants...")

	for _, plugin := range plugins {
		pluginLower := strings.ToLower(plugin.Name)

		// Performance level detection
		if strings.Contains(pluginLower, "high") || strings.Contains(pluginLower, "ultra") || strings.Contains(pluginLower, "max") {
			performanceMap["high"] = PluginMatch{
				Name:       plugin.Name,
				Confidence: 0.7,
				Reason:     "High performance variant",
				Type:       plugin.GetPluginType(make(map[string]string), make(map[string]bool)),
			}
			log.Printf("     High performance variant: %s", plugin.Name)
		} else if strings.Contains(pluginLower, "medium") || strings.Contains(pluginLower, "mid") || strings.Contains(pluginLower, "standard") {
			performanceMap["medium"] = PluginMatch{
				Name:       plugin.Name,
				Confidence: 0.8, // Often the most compatible
				Reason:     "Medium performance variant",
				Type:       plugin.GetPluginType(make(map[string]string), make(map[string]bool)),
			}
			log.Printf("     Medium performance variant: %s", plugin.Name)
		} else if strings.Contains(pluginLower, "low") || strings.Contains(pluginLower, "lite") || strings.Contains(pluginLower, "performance") {
			performanceMap["low"] = PluginMatch{
				Name:       plugin.Name,
				Confidence: 0.6,
				Reason:     "Low performance variant",
				Type:       plugin.GetPluginType(make(map[string]string), make(map[string]bool)),
			}
			log.Printf("     Low performance variant: %s", plugin.Name)
		}
	}

	return performanceMap
}

// Enhanced result formatting with detailed analysis
func formatEnhancedDetectionResults(analysis *EnhancedFomodAnalysis) string {
	var results []string

	// Add quality indicator with details
	qualityDetails := fmt.Sprintf("Quality: %s (%d/%d steps, %.0f%% avg confidence)",
		analysis.DetectionQuality,
		analysis.SelectedCount,
		analysis.TotalSteps,
		calculateAverageConfidence(analysis))
	results = append(results, qualityDetails)

	// Add step results with enhanced information
	for _, stepResult := range analysis.StepResults {
		stepName := stepResult.StepName
		if stepName == "" {
			stepName = fmt.Sprintf("Step %d", stepResult.StepIndex+1)
		}

		if !stepResult.IsVisible {
			results = append(results, fmt.Sprintf("%s: [Hidden by conditions]", stepName))
		} else if stepResult.BestPlugin != "" {
			confidence := fmt.Sprintf("%.0f%%", stepResult.Confidence*100)

			// Add type information
			typeInfo := ""
			if stepResult.GroupType != "" {
				typeInfo = fmt.Sprintf(" [%s]", stepResult.GroupType)
			}

			// Add alternative information
			altInfo := ""
			if len(stepResult.AlternativePlugins) > 0 {
				highConfAlts := 0
				for _, alt := range stepResult.AlternativePlugins {
					if alt.Confidence >= stepResult.Confidence*0.9 {
						highConfAlts++
					}
				}
				if highConfAlts > 0 {
					altInfo = fmt.Sprintf(" (+%d close alternatives)", highConfAlts)
				} else {
					altInfo = fmt.Sprintf(" (+%d alternatives)", len(stepResult.AlternativePlugins))
				}
			}

			// Add dependency/conflict warnings
			warningInfo := ""
			if len(stepResult.RequiredDependencies) > 0 {
				warningInfo += " [REQUIRED]"
			}
			if len(stepResult.ConflictingChoices) > 0 {
				warningInfo += " [CONFLICTS]"
			}

			results = append(results, fmt.Sprintf("%s: %s (%s)%s%s%s",
				stepName, stepResult.BestPlugin, confidence, typeInfo, altInfo, warningInfo))
		} else {
			results = append(results, fmt.Sprintf("%s: [No suitable match found]", stepName))
		}
	}

	// Add detailed summary information
	summaryParts := []string{}

	if analysis.RequiredFileMatches > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("Required files: %d matches", analysis.RequiredFileMatches))
	}

	if analysis.ConditionalFileMatches > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("Conditional files: %d matches", analysis.ConditionalFileMatches))
	}

	if len(analysis.RecommendedChoices) > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("High confidence: %d choices", len(analysis.RecommendedChoices)))
	}

	if len(analysis.PotentialConflicts) > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("⚠️ Conflicts: %d detected", len(analysis.PotentialConflicts)))
	}

	if len(analysis.MissingDependencies) > 0 {
		summaryParts = append(summaryParts, fmt.Sprintf("❌ Missing deps: %d", len(analysis.MissingDependencies)))
	}

	if len(summaryParts) > 0 {
		results = append(results, strings.Join(summaryParts, " | "))
	}

	summary := strings.Join(results, " | ")

	// Enhanced logging
	log.Printf("🎯 Enhanced FOMOD Detection Complete:")
	log.Printf("   Quality: %s", analysis.DetectionQuality)
	log.Printf("   Detection Rate: %d/%d steps (%.1f%%)",
		analysis.SelectedCount, analysis.TotalSteps,
		float64(analysis.SelectedCount)/float64(analysis.TotalSteps)*100)
	log.Printf("   Average Confidence: %.1f%%", calculateAverageConfidence(analysis))

	if len(analysis.PotentialConflicts) > 0 {
		log.Printf("   ⚠️ Potential Conflicts:")
		for _, conflict := range analysis.PotentialConflicts {
			log.Printf("      %s", conflict)
		}
	}

	if len(analysis.MissingDependencies) > 0 {
		log.Printf("   ❌ Missing Dependencies:")
		for _, dep := range analysis.MissingDependencies {
			log.Printf("      %s", dep)
		}
	}

	return summary
}

// Missing helper functions from original code

// buildArchiveFileMap creates a comprehensive map of all files in the archive
func buildArchiveFileMap(tempDir string) (map[string]ArchiveFile, error) {
	archiveFiles := make(map[string]ArchiveFile)
	fileCount := 0

	log.Printf("🔍 Building archive file map from: %s", tempDir)

	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("❌ Walk error for %s: %v", path, err)
			return nil
		}
		if info.IsDir() {
			log.Printf("📁 Directory: %s", path)
			return nil
		}

		fileCount++
		relativePath, err := filepath.Rel(tempDir, path)
		if err != nil {
			log.Printf("❌ Failed to get relative path for %s: %v", path, err)
			return nil
		}

		// Normalize path separators
		normalizedPath := strings.ReplaceAll(relativePath, "\\", "/")

		// Compute hash
		hash, err := utils.HashFile(path)
		if err != nil {
			log.Printf("❌ Failed to hash file %s: %v", path, err)
			return nil
		}

		archiveFiles[normalizedPath] = ArchiveFile{
			RelativePath: normalizedPath,
			FullPath:     path,
			Hash:         hash,
			Size:         info.Size(),
		}

		log.Printf("📄 [%d] %s (hash: %s, size: %d)", fileCount, normalizedPath, hash, info.Size())

		return nil
	})

	log.Printf("✅ Archive indexing complete: %d files processed", fileCount)
	return archiveFiles, err
}

// buildModFileMap creates lookup maps for installed mod files
func buildModFileMap(modFiles []dtos.ModFileDTO) map[string]dtos.ModFileDTO {
	modFileMap := make(map[string]dtos.ModFileDTO)

	log.Printf("🗂️ Building mod file map from %d installed files:", len(modFiles))

	for i, modFile := range modFiles {
		normalizedPath := strings.ReplaceAll(modFile.Path, "\\", "/")

		// Store both original case and lowercase for flexible matching
		modFileMap[normalizedPath] = modFile
		modFileMap[strings.ToLower(normalizedPath)] = modFile

		log.Printf("📋 [%d] %s (hash: %s, size: %d)", i+1, normalizedPath, modFile.Hash, modFile.Size)
	}

	log.Printf("✅ Mod file mapping complete: %d files mapped (%d with case variants)", len(modFiles), len(modFileMap))
	return modFileMap
}

// buildInstalledFileMap creates a boolean map for quick file existence checks
func buildInstalledFileMap(modFileMap map[string]dtos.ModFileDTO) map[string]bool {
	installedFiles := make(map[string]bool)

	for path := range modFileMap {
		installedFiles[path] = true
	}

	return installedFiles
}

// isFileInFolder checks if a file path is within a folder path
func isFileInFolder(filePath, folderPath string) bool {
	log.Printf("                 🔍 Checking if '%s' is in folder '%s'", filePath, folderPath)

	// Handle various folder path patterns
	if folderPath == "" {
		log.Printf("                 ✅ Empty folder path - matches all files")
		return true // Empty folder path means all files
	}

	// Direct prefix match
	if strings.HasPrefix(filePath, folderPath+"/") {
		log.Printf("                 ✅ Direct prefix match")
		return true
	}

	// Exact match (file is the folder itself)
	if filePath == folderPath {
		log.Printf("                 ✅ Exact match")
		return true
	}

	// Suffix matching for files with prefixes
	pathParts := strings.Split(filePath, "/")
	log.Printf("                 🔍 Path parts: %v", pathParts)

	for i, part := range pathParts {
		if part == folderPath || strings.Contains(part, folderPath) {
			log.Printf("                 ✅ Found folder part at index %d: %s", i, part)
			return true
		}
	}

	// Check if any part of the path starts with the folder name
	for i, part := range pathParts {
		if strings.HasPrefix(part, folderPath) {
			log.Printf("                 ✅ Found folder prefix at index %d: %s", i, part)
			return true
		}
	}

	log.Printf("                 ❌ File not in folder")
	return false
}

// calculateDestinationPath computes the final destination path for a file
func calculateDestinationPath(archivePath, sourcePath, destPath string) string {
	log.Printf("                 🎯 Calculating destination:")
	log.Printf("                     Archive: %s", archivePath)
	log.Printf("                     Source:  %s", sourcePath)
	log.Printf("                     Dest:    %s", destPath)

	if destPath == "" {
		log.Printf("                 📁 Empty destination - stripping source prefix")

		// Empty destination means strip the source prefix
		if strings.HasPrefix(archivePath, sourcePath+"/") {
			result := strings.TrimPrefix(archivePath, sourcePath+"/")
			log.Printf("                 ✅ Direct prefix strip: %s", result)
			return result
		}

		// Handle complex prefix scenarios
		parts := strings.Split(archivePath, "/")
		log.Printf("                 🔍 Archive parts: %v", parts)

		for i, part := range parts {
			if strings.Contains(part, sourcePath) || part == sourcePath {
				if i+1 < len(parts) {
					result := strings.Join(parts[i+1:], "/")
					log.Printf("                 ✅ Complex prefix strip at index %d: %s", i, result)
					return result
				}
			}
		}

		log.Printf("                 ⚠️ Fallback to original path: %s", archivePath)
		return archivePath
	}

	log.Printf("                 📂 Non-empty destination - calculating relative path")

	// Non-empty destination
	relativePath := strings.TrimPrefix(archivePath, sourcePath+"/")
	if relativePath == archivePath {
		log.Printf("                 🔍 No direct prefix - trying complex matching")

		// Handle complex scenarios
		parts := strings.Split(archivePath, "/")
		for i, part := range parts {
			if strings.Contains(part, sourcePath) {
				if i+1 < len(parts) {
					relativePath = strings.Join(parts[i+1:], "/")
					log.Printf("                 ✅ Found relative path at index %d: %s", i, relativePath)
					break
				}
			}
		}
	} else {
		log.Printf("                 ✅ Direct relative path: %s", relativePath)
	}

	result := destPath + "/" + relativePath
	log.Printf("                 🎯 Final destination: %s", result)
	return result
}

// cleanupTempDir removes the temporary directory
func cleanupTempDir(tempDir string) {
	go func() {
		if err := os.RemoveAll(tempDir); err != nil {
			log.Printf("Failed to cleanup temp directory: %v", err)
		}
	}()
}

// Helper function to calculate average confidence
func calculateAverageConfidence(analysis *EnhancedFomodAnalysis) float64 {
	if len(analysis.StepResults) == 0 {
		return 0.0
	}

	total := 0.0
	count := 0

	for _, result := range analysis.StepResults {
		if result.IsVisible {
			total += result.Confidence
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return (total / float64(count)) * 100
}
