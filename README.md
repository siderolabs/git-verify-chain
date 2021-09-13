# git-verify-chain

A tool to verify git commit signatures since the given tag/commit.

```
$ git-verify-chain -from v0.12.0
2021/09/13 14:41:58 OK

$ git-verify-chain -from v0.10.0
2021/09/13 14:42:04 failed to verify commit "faecae44fde60fc626ccb01da3b221519a9d41d7":
git verify-commit --verbose faecae44fde60fc626ccb01da3b221519a9d41d7: exit status 1
```

It is a wrapper around `git` and `gpg` tools.
They are invoked directly instead of handling git/GnuPG internals with native Go code as it already
[done](https://github.com/talos-systems/conform/blob/master/internal/git/git.go)
by [conform](https://github.com/talos-systems/conform).
