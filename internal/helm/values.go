package helm

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Values represents Helm chart values
type Values map[string]interface{}

// MergeValues merges multiple values maps with priority to later values
// Later values override earlier values for the same keys
func MergeValues(valueMaps ...Values) Values {
	result := make(Values)
	
	for _, values := range valueMaps {
		result = mergeRecursive(result, values)
	}
	
	return result
}

// mergeRecursive recursively merges two values maps
func mergeRecursive(dst, src Values) Values {
	for key, srcVal := range src {
		if dstVal, exists := dst[key]; exists {
			// Both values exist, check if they are maps
			srcMap, srcIsMap := srcVal.(map[string]interface{})
			dstMap, dstIsMap := dstVal.(map[string]interface{})
			
			if srcIsMap && dstIsMap {
				// Both are maps, merge recursively
				dst[key] = mergeRecursive(dstMap, srcMap)
			} else {
				// Not both maps, override with source value
				dst[key] = srcVal
			}
		} else {
			// Key doesn't exist in destination, add it
			dst[key] = srcVal
		}
	}
	
	return dst
}

// LoadValuesFile loads values from a YAML file
func LoadValuesFile(path string) (Values, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read values file %s: %w", path, err)
	}
	
	var values Values
	if err := yaml.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("failed to parse values file %s: %w", path, err)
	}
	
	return values, nil
}

// LoadValuesFiles loads and merges multiple values files
// Later files override earlier files for the same keys
func LoadValuesFiles(paths ...string) (Values, error) {
	var valueMaps []Values
	
	for _, path := range paths {
		values, err := LoadValuesFile(path)
		if err != nil {
			return nil, err
		}
		valueMaps = append(valueMaps, values)
	}
	
	return MergeValues(valueMaps...), nil
}

// ToYAML converts values to YAML string
func (v Values) ToYAML() (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal values to YAML: %w", err)
	}
	return string(data), nil
}

// WriteToFile writes values to a YAML file
func (v Values) WriteToFile(path string) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal values to YAML: %w", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write values file %s: %w", path, err)
	}
	
	return nil
}

// Get retrieves a value by key path (e.g., "app.name" or "database.host")
func (v Values) Get(keyPath string) (interface{}, bool) {
	return getNestedValue(v, keyPath)
}

// Set sets a value by key path (e.g., "app.name" or "database.host")
func (v Values) Set(keyPath string, value interface{}) {
	setNestedValue(v, keyPath, value)
}

// getNestedValue retrieves a nested value using dot notation
func getNestedValue(values Values, keyPath string) (interface{}, bool) {
	keys := splitKeyPath(keyPath)
	current := interface{}(values)
	
	for _, key := range keys {
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[key]; exists {
				current = val
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	
	return current, true
}

// setNestedValue sets a nested value using dot notation
func setNestedValue(values Values, keyPath string, value interface{}) {
	keys := splitKeyPath(keyPath)
	current := values
	
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		if _, exists := current[key]; !exists {
			current[key] = make(map[string]interface{})
		}
		
		if m, ok := current[key].(map[string]interface{}); ok {
			current = m
		} else {
			// Overwrite non-map value with a map
			newMap := make(map[string]interface{})
			current[key] = newMap
			current = newMap
		}
	}
	
	current[keys[len(keys)-1]] = value
}

// splitKeyPath splits a key path by dots
func splitKeyPath(keyPath string) []string {
	var keys []string
	current := ""
	
	for _, char := range keyPath {
		if char == '.' {
			if current != "" {
				keys = append(keys, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		keys = append(keys, current)
	}
	
	return keys
}
