# bumper (go-bumper)

[![Go](https://github.com/hedzr/go-bumper/actions/workflows/go.yml/badge.svg)](https://github.com/hedzr/go-bumper/actions/workflows/go.yml)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/go-bumper.svg?label=release)](https://github.com/hedzr/go-bumper/releases)
[![](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/go-bumper)

go-bumper is a CLI app to provide the formula versions bumping 
operation for Homebrew.

We made Bumper because `brew bump-formula-pr` cannot update for 
multiple sha256 values in a formula. I reviewed its source but
the relevant codes might be [hard to modify](https://github.com/Homebrew/brew/blob/1ca3ed87e28c450a24ee144a23fe2ba8b2a73640/Library/Homebrew/dev-cmd/bump-formula-pr.rb#L145-L341).

So here is go-bumper. For its usage you may check the go.yml in 
[go-faker](https://github.com/hedzr/go-faker) and [go-faker releases
page](https://github.com/hedzr/go-faker/releases).

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

[go-bumper](https://github.com/hedzr/go-bumper) is written for bumping [go-faker](https://github.com/hedzr/go-faker) automatically perfectly. my 
old github actions cannot work properly with `brew bump-formula-pr`, 
because there are three or more urls for darwin/linux and multiple
cpu arch. As we knew brew can bump-formula-pr fine on a single url
formula, and [go-bumper](https://github.com/hedzr/go-bumper) is yet another one but for multiple urls.

[go-bumper](https://github.com/hedzr/go-bumper) is still WIP.



## Conventions

version pattern: `x.y.z[-xxx]`
- eg: 1.2.3, 1.2.3-b1, 1.2.3-beta29, 1.2.3-alpha, 1.2.3-7, ...

## GitHub Personal Token

To use bumper, a personal access token with read-write permissions (such as `GH_RW_TOKEN`) should be specified into your repo Secrets section. See go-faker go.yml for instance, the command-line sample:

```bash
echo .pr/bin/bumper brew -act $ACTOR -am $ACTOR_EMAIL -f $APP_NAME \
  -ref $GITHUB_REF -ver $RELEASE_VERSION -t $HOMEBREW_TAP \
  --sha $BINARIES_ASC \
  --push
.pr/bin/bumper brew -act $ACTOR -am $ACTOR_EMAIL -f $APP_NAME \
  -ref $GITHUB_REF -ver $RELEASE_VERSION -t $HOMEBREW_TAP \
  --sha $BINARIES_ASC \
  --token ${{ secrets.GH_RW_TOKEN }} \
  --push
```


## LICENSE

Apache-2.0

