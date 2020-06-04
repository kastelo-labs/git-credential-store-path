# git-credential-store-path

This is a Git credential helper that can serve different credentials based
on the repo path. It's useful when you have multiple GitHub users with
access to different organizations.

## Installation

Make sure that the `git-credential-store-path` binary is in your `$PATH`.

Add the following to your `~/.gitconfig`:

```
[credential]
	helper = store-path
[url "https://github.com/"]
    insteadOf = git@github.com:
```

The URL replacement is to force usage of HTTPS instead of SSH, as the
credentials helper can't assist with SSH connections.

## Usage

The helper reads `~/.git-credentials` in the same way as the regular
`git-credentials-store` helper and could be used as a drop in replacement.
In addition, it also does matching on path prefixes and uses the first
matching entry. For example, the following `.git-credentials` file:

```
https://org1-user:someapikey@github.com/Org1/
https://org2-user:someapikey@github.com/Org2/
https://myuser:someapikey@github.com
```

Operations on a repository `Org1/foo` will use the `org1-user` credentials,
operations on `Org2/bar` will use the `org2-user` credentials, and all other
repos will use `myuser`.
