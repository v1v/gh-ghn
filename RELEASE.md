# Release

The release process is based on the GitHub CLI release process and utilizes the
[gh-extension-precompile](https://github.com/cli/gh-extension-precompile) action.

## Release Process

The release process is as follows:
- Create a tag in the format `vX.Y.Z` on the main branch and push it to the repository.
  ```shell
  git checkout main
  git tag v1.0.0
  git push origin v1.0.0
  ```
- The release workflow will be triggered and creates a release with the precompiled binaries and assets attached.