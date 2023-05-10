#!/bin/bash

# Get current version
current_version=$(grep 'var version' version.go | cut -d '"' -f 2)
echo "Current version: $current_version"
# Prompt user for new version
read -p "Enter new version (in semver format): " new_version

# Validate new version format
if ! [[ "$new_version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version format. Must be in semver format (e.g. v1.2.3)"
  exit 1
fi

# Replace version in version.go
sed -i.bak "s/$current_version/$new_version/g" version.go

# Delete backup file
rm version.go.bak

# Create a git commit with the version bump
git add version.go
git commit -m "Bump version to $new_version"

# Create git tag
git tag -a "$new_version" -m "Version $new_version"

echo "Please Push changes and tag to remote"
echo "git push origin && git push origin --tags"