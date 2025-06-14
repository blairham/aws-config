package flags

import (
	"flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsage(t *testing.T) {
	// Test with no flags
	usageText := "This is sample usage text"
	result := Usage(usageText, nil)
	assert.Equal(t, usageText, result)

	// Test with flags
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Bool("test-flag", false, "Test flag")
	fs.String("config", "", "Config file path")

	usageWithFlags := Usage(usageText, fs)
	assert.Equal(t, usageText, usageWithFlags, "Expected Usage to return the text as is even with flags")
}

func TestUsager(t *testing.T) {
	// Test with no flags
	usageText := "This is sample usage text"
	usager := Usager{Usage: usageText}
	result := usager.String()
	assert.Equal(t, usageText, result)

	// Test with space trimming
	spacedUsage := "  Text with spaces  \n\n"
	usager = Usager{Usage: spacedUsage}
	result = usager.String()
	assert.Equal(t, "Text with spaces", result)
}

func TestUsagerWithFlags(t *testing.T) {
	usageText := "Command usage"
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Bool("test-flag", false, "Test flag")

	usager := Usager{Usage: usageText, Flags: fs}
	result := usager.String()

	// Should include the usage text
	assert.True(t, strings.Contains(result, usageText), "Result should contain the usage text")
}
