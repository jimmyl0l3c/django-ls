#!/bin/bash

if [[ "$#" -ne 1 ]]; then
    echo "Usage: $0 <version>"
    exit 1
fi

version="$1"

if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Invalid format: $version"
    echo "Expected semver without pre-release tags (e.g., 1.2.3)"
    exit 1
fi

current_branch=$(git branch --show-current)
if [[ "$current_branch" != "master" ]]; then
    echo "Error: Not on default branch (currently on: $current_branch)"
    exit 1
fi

echo "Releasing: $version"

echo "Updating meta.go"
sed -i "s/[0-9]\+\.[0-9]\+\.[0-9]\+/$version/" meta.go || exit 1
git add meta.go || exit 1
git commit -S -m "deploy: 🚀 v$version" || exit 1

echo "Tagging version"
git tag "v$version" || exit 1

echo "Pushing changes"
git push --follow-tags || exit 1
