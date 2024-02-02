# naive-nx

This is a wrapper around `nx` command helps with running pre-determined set of targets which are:

- type-check
- lint
- test

The exposed binary is `naive-nx`. This wrapper is a stopgap for `nx affected`, it will look at diff-ed filenames and assumes that those are the only affected projects - hence the *naive* name.

## Installation

Download a release, run `chmod +x naive-nx` and move to appropriate location such as `sudo cp naive-nx /usr/local/bin` so you can launch it anywhere.

## Usage

`naive-nx --help` shows API.

The command can be executed even if you're in a sub-directory of a repository. It will diff against `master` branch by default. You can override this with flag `naive-nx --base-ref origin/me/some-feature-branch` for example.

You can pass flag `naive-nx --stubborn`, this will force to run `test`, `type-check` and `lint` for detected projects without checking whether those targets are available. This command is faster but it can fail, so it's a tradeoff.
