package main

import (
	"github.com/hedzr/go-bumper/bumper/cmdrrel"
	"os"
)

// main for testing.
// You should use bumper/main.go as the real main entry
func main() {
	// $GITHUB_REF, $GITHUB_ACTION, $GITHUB_TOKEN
	_ = os.Setenv("GITHUB_REF", "refs/tags/v0.2.4")
	_ = os.Setenv("GITHUB_ACTOR", "hedzr")
	_ = os.Setenv("GITHUB_ACTOR_MAIL", "hedzrz@gmail.com")
	_ = os.Setenv("TAP", "hedzr/homebrew-brew")
	_ = os.Setenv("FORMULA", "faker")

	//print("hello, bumper [--tap=homebrew-brew] []")
	cmdrrel.Entry()
}
