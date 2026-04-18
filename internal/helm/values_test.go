package helm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestMergeValues tests merging multiple values maps
func TestMergeValues(t *testing.T) {
	tests := []struct {
		name     string
		values   []Values
		expected Values
	}{
		{
			name: "merge two simple maps",
			values: []Values{
				{"key1": "value1", "key2": "value2"},
				{"key2": "override", "key3": "value3"},
			},
			expected: Values{
				"key1": "value1",
				"key2": "override",
				"key3": "value3",
			},
		},
		{
			name: "merge nested maps",
			values: []Values{
				{"app": map[string]interface{}{"name": "myapp", "version": "1.0"}},
				{"app": map[string]interface{}{"version": "2.0", "env": "prod"}},
			},
			expected: Values{
				"app": map[string]interface{}{
					"name":    "myapp",
					"version": "2.0",
					"env":     "prod",
				},
			},
		},
		{
			name: "merge three maps",
			values: []Values{
				{"a": 1},
				{"b": 2},
				{"c": 3},
			},
			expected: Values{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
		{
			name:     "merge empty maps",
			values:   []Values{{}, {}},
			expected: Values{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeValues(tt.values...)
			
			for key, expectedVal := range tt.expected {
				actualVal, exists := result[key]
				testutil.AssertTrue(t, exists, "Key should exist: "+key)
				
				// For nested maps, compare recursively
				if expectedMap, ok := expectedVal.(map[string]interface{}); ok {
					// Try both map[string]interface{} and Values type
					var actualMap map[string]interface{}
					if m, ok := actualVal.(map[string]interface{}); ok {
						actualMap = m
					} else if m, ok := actualVal.(Values); ok {
						actualMap = map[string]interface{}(m)
					} else {
						testutil.AssertTrue(t, false, "Value should be a map")
						continue
					}
					
					for nestedKey, nestedExpected := range expectedMap {
						nestedActual := actualMap[nestedKey]
						testutil.AssertEqual(t, nestedExpected, nestedActual)
					}
				} else {
					testutil.AssertEqual(t, expectedVal, actualVal)
				}
			}
		})
	}
}

// TestLoadValuesFile tests loading values from a YAML file
func TestLoadValuesFile(t *testing.T) {
	// Create a temporary values file
	tmpDir := t.TempDir()
	valuesPath := filepath.Join(tmpDir, "values.yaml")
	
	yamlContent := `
app:
  name: myapp
  replicas: 3
  image:
    repository: nginx
    tag: latest
database:
  host: localhost
  port: 5432
`
	
	err := os.WriteFile(valuesPath, []byte(yamlContent), 0644)
	testutil.AssertNoError(t, err)

	// Load the values
	values, err := LoadValuesFile(valuesPath)
	testutil.AssertNoError(t, err)
	testutil.AssertTrue(t, values != nil, "Values should not be nil")

	// Verify values
	app, exists := values["app"]
	testutil.AssertTrue(t, exists, "app key should exist")
	
	// Handle both map[string]interface{} and Values type
	var appMap map[string]interface{}
	if m, ok := app.(map[string]interface{}); ok {
		appMap = m
	} else if m, ok := app.(Values); ok {
		appMap = map[string]interface{}(m)
	} else {
		t.Fatalf("app should be a map, got %T", app)
	}
	
	testutil.AssertEqual(t, "myapp", appMap["name"])
	testutil.AssertEqual(t, 3, appMap["replicas"])
}

// TestLoadValuesFileNotFound tests loading a non-existent file
func TestLoadValuesFileNotFound(t *testing.T) {
	_, err := LoadValuesFile("/nonexistent/values.yaml")
	testutil.AssertError(t, err)
	testutil.AssertContains(t, err.Error(), "failed to read values file")
}

// TestLoadValuesFileInvalidYAML tests loading invalid YAML
func TestLoadValuesFileInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	valuesPath := filepath.Join(tmpDir, "invalid.yaml")
	
	invalidYAML := `
app:
  name: myapp
  invalid yaml here: [
`
	
	err := os.WriteFile(valuesPath, []byte(invalidYAML), 0644)
	testutil.AssertNoError(t, err)

	_, err = LoadValuesFile(valuesPath)
	testutil.AssertError(t, err)
	testutil.AssertContains(t, err.Error(), "failed to parse values file")
}

// TestLoadValuesFiles tests loading and merging multiple values files
func TestLoadValuesFiles(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create base values file
	baseValues := filepath.Join(tmpDir, "base.yaml")
	baseYAML := `
app:
  name: myapp
  replicas: 2
`
	err := os.WriteFile(baseValues, []byte(baseYAML), 0644)
	testutil.AssertNoError(t, err)

	// Create override values file
	overrideValues := filepath.Join(tmpDir, "override.yaml")
	overrideYAML := `
app:
  replicas: 5
  env: production
`
	err = os.WriteFile(overrideValues, []byte(overrideYAML), 0644)
	testutil.AssertNoError(t, err)

	// Load and merge
	values, err := LoadValuesFiles(baseValues, overrideValues)
	testutil.AssertNoError(t, err)

	appVal := values["app"]
	var app map[string]interface{}
	if m, ok := appVal.(map[string]interface{}); ok {
		app = m
	} else if m, ok := appVal.(Values); ok {
		app = map[string]interface{}(m)
	} else {
		t.Fatalf("app should be a map, got %T", appVal)
	}
	
	testutil.AssertEqual(t, "myapp", app["name"])
	testutil.AssertEqual(t, 5, app["replicas"]) // Should be overridden
	testutil.AssertEqual(t, "production", app["env"])
}

// TestValuesToYAML tests converting values to YAML string
func TestValuesToYAML(t *testing.T) {
	values := Values{
		"app": map[string]interface{}{
			"name":     "myapp",
			"replicas": 3,
		},
	}

	yaml, err := values.ToYAML()
	testutil.AssertNoError(t, err)
	testutil.AssertContains(t, yaml, "app:")
	testutil.AssertContains(t, yaml, "name: myapp")
	testutil.AssertContains(t, yaml, "replicas: 3")
}

// TestValuesWriteToFile tests writing values to a file
func TestValuesWriteToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.yaml")

	values := Values{
		"app": map[string]interface{}{
			"name": "myapp",
		},
	}

	err := values.WriteToFile(outputPath)
	testutil.AssertNoError(t, err)

	// Verify file was created
	_, err = os.Stat(outputPath)
	testutil.AssertNoError(t, err)

	// Load and verify content
	loaded, err := LoadValuesFile(outputPath)
	testutil.AssertNoError(t, err)
	
	appVal := loaded["app"]
	var app map[string]interface{}
	if m, ok := appVal.(map[string]interface{}); ok {
		app = m
	} else if m, ok := appVal.(Values); ok {
		app = map[string]interface{}(m)
	} else {
		t.Fatalf("app should be a map, got %T", appVal)
	}
	
	testutil.AssertEqual(t, "myapp", app["name"])
}

// TestValuesGet tests getting nested values
func TestValuesGet(t *testing.T) {
	values := Values{
		"app": map[string]interface{}{
			"name": "myapp",
			"database": map[string]interface{}{
				"host": "localhost",
				"port": 5432,
			},
		},
	}

	tests := []struct {
		name     string
		keyPath  string
		expected interface{}
		exists   bool
	}{
		{
			name:     "get top level key",
			keyPath:  "app",
			expected: values["app"],
			exists:   true,
		},
		{
			name:     "get nested key",
			keyPath:  "app.name",
			expected: "myapp",
			exists:   true,
		},
		{
			name:     "get deeply nested key",
			keyPath:  "app.database.host",
			expected: "localhost",
			exists:   true,
		},
		{
			name:     "get non-existent key",
			keyPath:  "app.nonexistent",
			expected: nil,
			exists:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, exists := values.Get(tt.keyPath)
			testutil.AssertEqual(t, tt.exists, exists)
			if tt.exists {
				testutil.AssertEqual(t, tt.expected, val)
			}
		})
	}
}

// TestValuesSet tests setting nested values
func TestValuesSet(t *testing.T) {
	values := Values{}

	// Set top level value
	values.Set("app", "myapp")
	val, exists := values.Get("app")
	testutil.AssertTrue(t, exists, "Key should exist")
	testutil.AssertEqual(t, "myapp", val)

	// Set nested value
	values.Set("database.host", "localhost")
	val, exists = values.Get("database.host")
	testutil.AssertTrue(t, exists, "Nested key should exist")
	testutil.AssertEqual(t, "localhost", val)

	// Set deeply nested value
	values.Set("app.config.timeout", 30)
	val, exists = values.Get("app.config.timeout")
	testutil.AssertTrue(t, exists, "Deeply nested key should exist")
	testutil.AssertEqual(t, 30, val)
}

// TestSplitKeyPath tests splitting key paths
func TestSplitKeyPath(t *testing.T) {
	tests := []struct {
		name     string
		keyPath  string
		expected []string
	}{
		{
			name:     "single key",
			keyPath:  "app",
			expected: []string{"app"},
		},
		{
			name:     "nested key",
			keyPath:  "app.name",
			expected: []string{"app", "name"},
		},
		{
			name:     "deeply nested key",
			keyPath:  "app.database.connection.host",
			expected: []string{"app", "database", "connection", "host"},
		},
		{
			name:     "empty key",
			keyPath:  "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitKeyPath(tt.keyPath)
			testutil.AssertEqual(t, len(tt.expected), len(result))
			for i, expected := range tt.expected {
				testutil.AssertEqual(t, expected, result[i])
			}
		})
	}
}

// TestMergeValuesOverride tests that later values override earlier ones
func TestMergeValuesOverride(t *testing.T) {
	base := Values{
		"replicas": 2,
		"image":    "nginx:1.19",
	}

	override := Values{
		"replicas": 5,
	}

	result := MergeValues(base, override)

	testutil.AssertEqual(t, 5, result["replicas"])
	testutil.AssertEqual(t, "nginx:1.19", result["image"])
}

// TestMergeValuesDeepNesting tests merging deeply nested structures
func TestMergeValuesDeepNesting(t *testing.T) {
	base := Values{
		"app": map[string]interface{}{
			"config": map[string]interface{}{
				"database": map[string]interface{}{
					"host": "localhost",
					"port": 5432,
				},
			},
		},
	}

	override := Values{
		"app": map[string]interface{}{
			"config": map[string]interface{}{
				"database": map[string]interface{}{
					"port": 3306,
					"user": "admin",
				},
			},
		},
	}

	result := MergeValues(base, override)

	// Navigate to the database config - handle Values type
	appVal := result["app"]
	var app map[string]interface{}
	if m, ok := appVal.(map[string]interface{}); ok {
		app = m
	} else if m, ok := appVal.(Values); ok {
		app = map[string]interface{}(m)
	} else {
		t.Fatalf("app should be a map, got %T", appVal)
	}
	
	configVal := app["config"]
	var config map[string]interface{}
	if m, ok := configVal.(map[string]interface{}); ok {
		config = m
	} else if m, ok := configVal.(Values); ok {
		config = map[string]interface{}(m)
	} else {
		t.Fatalf("config should be a map, got %T", configVal)
	}
	
	databaseVal := config["database"]
	var database map[string]interface{}
	if m, ok := databaseVal.(map[string]interface{}); ok {
		database = m
	} else if m, ok := databaseVal.(Values); ok {
		database = map[string]interface{}(m)
	} else {
		t.Fatalf("database should be a map, got %T", databaseVal)
	}

	testutil.AssertEqual(t, "localhost", database["host"]) // Preserved from base
	testutil.AssertEqual(t, 3306, database["port"])        // Overridden
	testutil.AssertEqual(t, "admin", database["user"])     // Added from override
}

// TestLoadValuesFilesWithFixtures tests loading values using test fixtures
func TestLoadValuesFilesWithFixtures(t *testing.T) {
	// Use test fixtures
	values, err := LoadValuesFile("../../test/fixtures/helm/values-dev.yaml")
	testutil.AssertNoError(t, err)
	testutil.AssertTrue(t, values != nil, "Values should not be nil")
}

// TestValuesSetOverwriteNonMap tests that Set overwrites non-map values with maps
func TestValuesSetOverwriteNonMap(t *testing.T) {
	values := Values{
		"app": "simple-string",
	}

	// Set a nested value under app, which should convert app to a map
	values.Set("app.name", "myapp")

	val, exists := values.Get("app.name")
	testutil.AssertTrue(t, exists, "Nested key should exist")
	testutil.AssertEqual(t, "myapp", val)

	// app should now be a map
	app, exists := values.Get("app")
	testutil.AssertTrue(t, exists, "app should exist")
	
	// Check if it's a map (either type)
	_, isMap := app.(map[string]interface{})
	if !isMap {
		_, isMap = app.(Values)
	}
	testutil.AssertTrue(t, isMap, "app should be a map")
}

// TestEmptyValues tests operations on empty values
func TestEmptyValues(t *testing.T) {
	values := Values{}

	// Get from empty values
	_, exists := values.Get("nonexistent")
	testutil.AssertFalse(t, exists, "Key should not exist in empty values")

	// Set in empty values
	values.Set("key", "value")
	val, exists := values.Get("key")
	testutil.AssertTrue(t, exists, "Key should exist after set")
	testutil.AssertEqual(t, "value", val)

	// ToYAML on empty values
	yaml, err := Values{}.ToYAML()
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, "{}\n", yaml)
}
