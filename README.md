# bumper (go-bumper)

Bumper is a golang CLI app to provide the formula versions bumping 
operation for Homebrew.

We made bumper because `brew bump-formula-pr` cannot update for 
multiple sha256 values in a formula. I reviewed its source but
the relevant codes might be [hard to modify](https://github.com/Homebrew/brew/blob/1ca3ed87e28c450a24ee144a23fe2ba8b2a73640/Library/Homebrew/dev-cmd/bump-formula-pr.rb#L145-L341).

So here is go-bumper. For its usage yoo may check the go.yml in 
[go-faker](https://github.com/hedzr/go-faker) and go-faker releases
page.

The flow of go-bumper needs:

- you're writing building flow with bash/shell in github actions 
- all binaries were built into bin/ and tar gzipped (*.tgz/*.tar.gz/...)
- `binaries.asc` was generated for all binaries with sha256sum
- you're updating the homebrew tap repo

go-bumper will:

- pull the tap repo and parse the formula.rb you specified
- bump the versions if matched `\d+\.\d+\.\d+` (can be customized)
- locate and update bottle url and its sha256sum (depend on `binaries.asc`)

## WIP

go-bumper is written for bumping go-faker automatically perfectly. my 
old github actions cannot work properly with `brew bump-formula-pr`, 
because there are three or more urls for darwin/linux and multiple
cpu arch. As we knew brew can bump-formula-pr fine on a single url
formula, and go-bumper is yet another one but for multiple urls.

go-bumper is WIP.

## LICENSE

Apache-2.0

