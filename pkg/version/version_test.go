package version

import (
	"testing"
)

func TestBuildInfo_String(t *testing.T) {
	info := BuildInfo{
		Version: "1.0.0",
		Commit:  "abc1234",
		Date:    "2025-01-01",
	}

	expected := "mana version 1.0.0\ncommit: abc1234\nbuilt: 2025-01-01"
	if got := info.String(); got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}
}

func TestBuildInfo_Getters(t *testing.T) {
	info := BuildInfo{
		Version: "1.0.0",
		Commit:  "abc1234",
		Date:    "2025-01-01",
	}

	t.Run("GetVersion", func(t *testing.T) {
		if got := info.GetVersion(); got != info.Version {
			t.Errorf("GetVersion() = %q, want %q", got, info.Version)
		}
	})

	t.Run("GetCommit", func(t *testing.T) {
		if got := info.GetCommit(); got != info.Commit {
			t.Errorf("GetCommit() = %q, want %q", got, info.Commit)
		}
	})

	t.Run("GetDate", func(t *testing.T) {
		if got := info.GetDate(); got != info.Date {
			t.Errorf("GetDate() = %q, want %q", got, info.Date)
		}
	})
}
