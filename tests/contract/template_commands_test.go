package contract

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	templateTestBinaryPath string
)

// buildTemplateTestBinary builds the my-context binary for testing
func buildTemplateTestBinary(t *testing.T) string {
	if templateTestBinaryPath != "" {
		return templateTestBinaryPath
	}

	// Find project root
	dir, err := os.Getwd()
	require.NoError(t, err)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}

	// Build binary
	tmpDir, err := os.MkdirTemp("", "template-test-*")
	require.NoError(t, err)

	binaryName := "my-context-template-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	templateTestBinaryPath = filepath.Join(tmpDir, binaryName)

	cmd := exec.Command("go", "build", "-o", templateTestBinaryPath, "./cmd/my-context/")
	cmd.Dir = dir
	err = cmd.Run()
	require.NoError(t, err, "Failed to build test binary")

	return templateTestBinaryPath
}

// runTemplateCommand executes a my-context command with isolated MY_CONTEXT_HOME
func runTemplateCommand(binary, testDir string, args ...string) (stdout string, exitCode int) {
	cmd := exec.Command(binary, args...)

	// Set up isolated environment
	env := make([]string, 0, len(os.Environ())+1)
	for _, e := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(e), "MY_CONTEXT_HOME=") {
			env = append(env, e)
		}
	}
	env = append(env, "MY_CONTEXT_HOME="+testDir)
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	stdout = strings.TrimSpace(string(output))

	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			exitCode = exit.ExitCode()
		} else {
			exitCode = 1
		}
	}

	return
}

// TestTemplateListEmpty tests listing templates when none exist
func TestTemplateListEmpty(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runTemplateCommand(binary, testDir, "template", "list")

	assert.Equal(t, 0, exitCode, "template list should succeed: %s", stdout)
	assert.Contains(t, stdout, "No templates found")
}

// TestTemplateCreateAndList tests creating and listing templates
func TestTemplateCreateAndList(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	// Create template file manually (since create is interactive)
	templatesDir := filepath.Join(testDir, "templates")
	err := os.MkdirAll(templatesDir, 0o755)
	require.NoError(t, err)

	templateContent := `name: test-template
description: Test template for contract tests
labels:
  - test
  - contract
initial_notes:
  - "Test note 1"
  - "Test note 2"
`
	templatePath := filepath.Join(templatesDir, "test-template.yaml")
	err = os.WriteFile(templatePath, []byte(templateContent), 0o600)
	require.NoError(t, err)

	// List templates
	stdout, exitCode := runTemplateCommand(binary, testDir, "template", "list")

	assert.Equal(t, 0, exitCode, "template list should succeed: %s", stdout)
	assert.Contains(t, stdout, "test-template")
	assert.Contains(t, stdout, "Test template for contract tests")
}

// TestTemplateShow tests showing template details
func TestTemplateShow(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	// Create template file
	templatesDir := filepath.Join(testDir, "templates")
	err := os.MkdirAll(templatesDir, 0o755)
	require.NoError(t, err)

	templateContent := `name: show-test
description: Template for show command test
labels:
  - label1
  - label2
initial_notes:
  - "Note 1"
  - "Note 2"
  - "Note 3"
`
	templatePath := filepath.Join(templatesDir, "show-test.yaml")
	err = os.WriteFile(templatePath, []byte(templateContent), 0o600)
	require.NoError(t, err)

	// Show template
	stdout, exitCode := runTemplateCommand(binary, testDir, "template", "show", "show-test")

	assert.Equal(t, 0, exitCode, "template show should succeed: %s", stdout)
	assert.Contains(t, stdout, "show-test")
	assert.Contains(t, stdout, "Template for show command test")
	assert.Contains(t, stdout, "label1")
	assert.Contains(t, stdout, "label2")
	assert.Contains(t, stdout, "Note 1")
	assert.Contains(t, stdout, "Note 2")
	assert.Contains(t, stdout, "Note 3")
}

// TestTemplateShowNotFound tests showing non-existent template
func TestTemplateShowNotFound(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runTemplateCommand(binary, testDir, "template", "show", "nonexistent")

	assert.NotEqual(t, 0, exitCode, "template show should fail for non-existent template")
	assert.Contains(t, stdout, "not found")
}

// TestStartWithTemplate tests starting a context with a template
func TestStartWithTemplate(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	// Create template file
	templatesDir := filepath.Join(testDir, "templates")
	err := os.MkdirAll(templatesDir, 0o755)
	require.NoError(t, err)

	templateContent := `name: debug
description: Debug template
labels:
  - debug
  - investigation
initial_notes:
  - "Investigation started"
  - "Hypothesis:"
  - "Findings:"
`
	templatePath := filepath.Join(templatesDir, "debug.yaml")
	err = os.WriteFile(templatePath, []byte(templateContent), 0o600)
	require.NoError(t, err)

	// Start context with template
	stdout, exitCode := runTemplateCommand(binary, testDir, "start", "Bug investigation", "--template", "debug")

	assert.Equal(t, 0, exitCode, "start with template should succeed: %s", stdout)
	assert.Contains(t, stdout, "Started: Bug investigation")

	// Show context to verify labels and notes
	stdout, exitCode = runTemplateCommand(binary, testDir, "show", "--json")
	assert.Equal(t, 0, exitCode, "show should succeed: %s", stdout)

	// Parse JSON output
	var result map[string]interface{}
	err = json.Unmarshal([]byte(stdout), &result)
	require.NoError(t, err, "Failed to parse JSON output: %s", stdout)

	dataWrapper := result["data"].(map[string]interface{})
	data := dataWrapper["data"].(map[string]interface{})

	// Check labels - context contains metadata
	contextData := data["context"].(map[string]interface{})
	metadata := contextData["metadata"].(map[string]interface{})
	labels := metadata["labels"].([]interface{})
	assert.Contains(t, labels, "debug")
	assert.Contains(t, labels, "investigation")

	// Check notes
	notes := data["notes"].([]interface{})
	assert.GreaterOrEqual(t, len(notes), 3, "Should have at least 3 notes from template")

	noteTexts := make([]string, 0)
	for _, n := range notes {
		noteMap := n.(map[string]interface{})
		noteTexts = append(noteTexts, noteMap["text_content"].(string))
	}
	assert.Contains(t, noteTexts, "Investigation started")
	assert.Contains(t, noteTexts, "Hypothesis:")
	assert.Contains(t, noteTexts, "Findings:")
}

// TestStartWithTemplateShortFlag tests starting a context with -T flag
func TestStartWithTemplateShortFlag(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	// Create template file
	templatesDir := filepath.Join(testDir, "templates")
	err := os.MkdirAll(templatesDir, 0o755)
	require.NoError(t, err)

	templateContent := `name: feature
description: Feature template
labels:
  - feature
initial_notes:
  - "Requirements:"
`
	templatePath := filepath.Join(templatesDir, "feature.yaml")
	err = os.WriteFile(templatePath, []byte(templateContent), 0o600)
	require.NoError(t, err)

	// Start context with -T flag
	stdout, exitCode := runTemplateCommand(binary, testDir, "start", "New feature", "-T", "feature")

	assert.Equal(t, 0, exitCode, "start with -T flag should succeed: %s", stdout)
	assert.Contains(t, stdout, "Started: New feature")
}

// TestStartWithTemplateNotFound tests starting with non-existent template
func TestStartWithTemplateNotFound(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	stdout, exitCode := runTemplateCommand(binary, testDir, "start", "Test", "--template", "nonexistent")

	assert.NotEqual(t, 0, exitCode, "start with non-existent template should fail")
	assert.Contains(t, stdout, "template error")
	assert.Contains(t, stdout, "not found")
}

// TestStartWithTemplateCombinesLabels tests that template labels combine with command-line labels
func TestStartWithTemplateCombinesLabels(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	// Create template file
	templatesDir := filepath.Join(testDir, "templates")
	err := os.MkdirAll(templatesDir, 0o755)
	require.NoError(t, err)

	templateContent := `name: combo
description: Combination test
labels:
  - template-label
initial_notes: []
`
	templatePath := filepath.Join(templatesDir, "combo.yaml")
	err = os.WriteFile(templatePath, []byte(templateContent), 0o600)
	require.NoError(t, err)

	// Start context with template and additional labels
	stdout, exitCode := runTemplateCommand(binary, testDir, "start", "Combined test", "--template", "combo", "--labels", "manual-label,another-label")

	assert.Equal(t, 0, exitCode, "start with combined labels should succeed: %s", stdout)

	// Show context to verify all labels are present
	stdout, exitCode = runTemplateCommand(binary, testDir, "show", "--json")
	assert.Equal(t, 0, exitCode, "show should succeed: %s", stdout)

	// Parse JSON output
	var result map[string]interface{}
	err = json.Unmarshal([]byte(stdout), &result)
	require.NoError(t, err, "Failed to parse JSON output: %s", stdout)

	dataWrapper := result["data"].(map[string]interface{})
	data := dataWrapper["data"].(map[string]interface{})
	contextData := data["context"].(map[string]interface{})
	metadata := contextData["metadata"].(map[string]interface{})
	labels := metadata["labels"].([]interface{})

	// Should have labels from both template and command line
	assert.Contains(t, labels, "template-label")
	assert.Contains(t, labels, "manual-label")
	assert.Contains(t, labels, "another-label")
}

// TestTemplateListJSON tests JSON output for template list
func TestTemplateListJSON(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	// Create template file
	templatesDir := filepath.Join(testDir, "templates")
	err := os.MkdirAll(templatesDir, 0o755)
	require.NoError(t, err)

	templateContent := `name: json-test
description: JSON test template
labels:
  - test
initial_notes:
  - "Note"
`
	templatePath := filepath.Join(templatesDir, "json-test.yaml")
	err = os.WriteFile(templatePath, []byte(templateContent), 0o600)
	require.NoError(t, err)

	// List templates with JSON output
	stdout, exitCode := runTemplateCommand(binary, testDir, "template", "list", "--json")

	assert.Equal(t, 0, exitCode, "template list --json should succeed: %s", stdout)

	// Parse JSON output
	var result map[string]interface{}
	err = json.Unmarshal([]byte(stdout), &result)
	require.NoError(t, err, "Failed to parse JSON output: %s", stdout)

	assert.Equal(t, "template_list", result["command"])

	dataWrapper := result["data"].(map[string]interface{})
	data := dataWrapper["data"].(map[string]interface{})
	templates := data["templates"].([]interface{})
	assert.Len(t, templates, 1)

	template := templates[0].(map[string]interface{})
	assert.Equal(t, "json-test", template["name"])
	assert.Equal(t, "JSON test template", template["description"])
}

// TestTemplateShowJSON tests JSON output for template show
func TestTemplateShowJSON(t *testing.T) {
	binary := buildTemplateTestBinary(t)
	testDir := t.TempDir()

	// Create template file
	templatesDir := filepath.Join(testDir, "templates")
	err := os.MkdirAll(templatesDir, 0o755)
	require.NoError(t, err)

	templateContent := `name: show-json
description: Show JSON test
labels:
  - label1
initial_notes:
  - "Note 1"
`
	templatePath := filepath.Join(templatesDir, "show-json.yaml")
	err = os.WriteFile(templatePath, []byte(templateContent), 0o600)
	require.NoError(t, err)

	// Show template with JSON output
	stdout, exitCode := runTemplateCommand(binary, testDir, "template", "show", "show-json", "--json")

	assert.Equal(t, 0, exitCode, "template show --json should succeed: %s", stdout)

	// Parse JSON output
	var result map[string]interface{}
	err = json.Unmarshal([]byte(stdout), &result)
	require.NoError(t, err, "Failed to parse JSON output: %s", stdout)

	assert.Equal(t, "template_show", result["command"])

	dataWrapper := result["data"].(map[string]interface{})
	data := dataWrapper["data"].(map[string]interface{})
	template := data["template"].(map[string]interface{})
	assert.Equal(t, "show-json", template["name"])
	assert.Equal(t, "Show JSON test", template["description"])
}
