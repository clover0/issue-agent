# Release

Releases are automated using tagpr and GoReleaser through GitHub Actions.

1. Merge PR created by tagpr.
1. tapr create versioned git tag.
1. GoReleaser builds and publishes the release.
