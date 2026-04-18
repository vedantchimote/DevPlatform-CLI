package logger

import (
	"testing"

	"github.com/devplatform/devplatform-cli/test/testutil"
)

// TestColorize tests the Colorize function
func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		color    Color
		contains []string
	}{
		{
			name:     "red color",
			text:     "error",
			color:    ColorRed,
			contains: []string{"error", "\033[31m", "\033[0m"},
		},
		{
			name:     "green color",
			text:     "success",
			color:    ColorGreen,
			contains: []string{"success", "\033[32m", "\033[0m"},
		},
		{
			name:     "yellow color",
			text:     "warning",
			color:    ColorYellow,
			contains: []string{"warning", "\033[33m", "\033[0m"},
		},
		{
			name:     "blue color",
			text:     "info",
			color:    ColorBlue,
			contains: []string{"info", "\033[34m", "\033[0m"},
		},
		{
			name:     "cyan color",
			text:     "debug",
			color:    ColorCyan,
			contains: []string{"debug", "\033[36m", "\033[0m"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Colorize(tt.text, tt.color)
			for _, substr := range tt.contains {
				testutil.AssertContains(t, result, substr)
			}
		})
	}
}

// TestColorHelpers tests the color helper functions
func TestColorHelpers(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string) string
		text     string
		contains string
	}{
		{
			name:     "Red",
			fn:       Red,
			text:     "error",
			contains: "\033[31m",
		},
		{
			name:     "Green",
			fn:       Green,
			text:     "success",
			contains: "\033[32m",
		},
		{
			name:     "Yellow",
			fn:       Yellow,
			text:     "warning",
			contains: "\033[33m",
		},
		{
			name:     "Blue",
			fn:       Blue,
			text:     "info",
			contains: "\033[34m",
		},
		{
			name:     "Cyan",
			fn:       Cyan,
			text:     "debug",
			contains: "\033[36m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.text)
			testutil.AssertContains(t, result, tt.text)
			testutil.AssertContains(t, result, tt.contains)
			testutil.AssertContains(t, result, string(ColorReset))
		})
	}
}

// TestColorConstants tests that color constants are defined
func TestColorConstants(t *testing.T) {
	testutil.AssertEqual(t, Color("\033[0m"), ColorReset)
	testutil.AssertEqual(t, Color("\033[31m"), ColorRed)
	testutil.AssertEqual(t, Color("\033[32m"), ColorGreen)
	testutil.AssertEqual(t, Color("\033[33m"), ColorYellow)
	testutil.AssertEqual(t, Color("\033[34m"), ColorBlue)
	testutil.AssertEqual(t, Color("\033[36m"), ColorCyan)
}

// TestColorizeEmptyString tests colorizing an empty string
func TestColorizeEmptyString(t *testing.T) {
	result := Colorize("", ColorRed)
	testutil.AssertContains(t, result, string(ColorRed))
	testutil.AssertContains(t, result, string(ColorReset))
}

// TestColorizeMultiline tests colorizing multiline text
func TestColorizeMultiline(t *testing.T) {
	text := "line1\nline2\nline3"
	result := Colorize(text, ColorGreen)
	testutil.AssertContains(t, result, "line1")
	testutil.AssertContains(t, result, "line2")
	testutil.AssertContains(t, result, "line3")
	testutil.AssertContains(t, result, string(ColorGreen))
	testutil.AssertContains(t, result, string(ColorReset))
}
