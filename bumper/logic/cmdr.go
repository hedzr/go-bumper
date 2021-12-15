package logic

import (
	"github.com/hedzr/cmdr"
)

func AttachToCmdr(root *cmdr.RootCmdOpt) {
	brewTool(root)
	debTool(root)
}

func brewTool(root *cmdr.RootCmdOpt) {

	cc := root.NewSubCommand("brew", "b", "homebrew").
		Description("homebrew tap bumper").
		Group("").
		//TailPlaceholder("[text1, text2, ...]").
		//PreAction(func(cmd *cmdr.Command, remainArgs []string) (err error) {
		//	fmt.Printf("[PRE] DebugMode=%v, TraceMode=%v. InDebugging/IsDebuggerAttached=%v\n",
		//		cmdr.GetDebugMode(), logex.GetTraceMode(), cmdr.InDebugging())
		//	for ix, s := range remainArgs {
		//		fmt.Printf("[PRE] %5d. %s\n", ix, s)
		//	}
		//
		//	fmt.Printf("[PRE] Debug=%v, Trace=%v\n", cmdr.GetDebugMode(), cmdr.GetTraceMode())
		//
		//	// return nil to be continue,
		//	// return cmdr.ErrShouldBeStopException to stop the following actions without error
		//	// return other errors for application purpose
		//	return
		//}).
		//PostAction(func(cmd *cmdr.Command, remainArgs []string) {
		//	for ix, s := range remainArgs {
		//		fmt.Printf("[POST] %5d. %s\n", ix, s)
		//	}
		//}).
		Action(brewIt)

	cmdr.NewString("").
		Titles("tap", "t").
		Description("tap repository (owner/repo for github, or fqdn)").
		Group("").
		Placeholder("NAME").
		EnvKeys("TAP").
		AttachTo(cc)

	cmdr.NewString("").
		Titles("formula", "f").
		Description("formula name").
		Group("").
		Placeholder("FORMULA").
		EnvKeys("FORMULA").
		AttachTo(cc)

	cmdr.NewString("").
		Titles("release-version", "ver", "ver").
		Description("release version").
		Group("").
		Placeholder("VERSION").
		EnvKeys("RELEASE_VERSION", "VERSION").
		AttachTo(cc)
	cmdr.NewString("").
		Titles("actor", "act", "github-actor").
		Description("github actor").
		Group("").
		Placeholder("NAME").
		EnvKeys("GITHUB_ACTOR").
		AttachTo(cc)
	cmdr.NewString("").
		Titles("actor-mail", "am", "github-actor-mail").
		Description("github actor email address").
		Group("").
		Placeholder("EMAIL").
		EnvKeys("GITHUB_ACTOR_MAIL").
		AttachTo(cc)
	cmdr.NewString("").
		Titles("ref", "ref", "github-ref").
		Description("github ref").
		Group("").
		Placeholder("REF").
		EnvKeys("GITHUB_REF").
		AttachTo(cc)
	cmdr.NewString("bin/binaries.asc").
		Titles("sha256file", "sha", "sha").
		Description("sha256 filename").
		Group("SHA").
		Placeholder("FILE").
		EnvKeys("SHA_FILE").
		AttachTo(cc)

	cmdr.NewString("").
		Titles("token", "tkn", "tkn", "github-token").
		Description("github token").
		Group("For Remote").
		Placeholder("TOKEN").
		EnvKeys("GITHUB_TOKEN").
		AttachTo(cc)
	cmdr.NewBool().
		Titles("push", "p").
		Description("push the new commit back to tap repository right now").
		Group("For Remote").
		AttachTo(cc)
}

func debTool(root *cmdr.RootCmdOpt) {

	cc := root.NewSubCommand("deb", "d").
		Description("deb bumper").
		Group("").
		TailPlaceholder("[text1, text2, ...]").
		Action(debIt)

	cmdr.NewBool().
		Titles("bool", "b").
		Description("bool option").
		Group("").
		AttachTo(cc)
}
