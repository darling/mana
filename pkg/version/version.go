package version

import "fmt"

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

func (b BuildInfo) String() string {
	return fmt.Sprintf("mana version %s\ncommit: %s\nbuilt: %s", b.Version, b.Commit, b.Date)
}

func (b BuildInfo) GetVersion() string {
	return b.Version
}

func (b BuildInfo) GetCommit() string {
	return b.Commit
}

func (b BuildInfo) GetDate() string {
	return b.Date
}
