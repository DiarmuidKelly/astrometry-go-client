#!/bin/bash
# Automated release script
# Bumps version, updates changelog, and creates git tag
# Usage: ./auto-release.sh [major|minor|patch]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Get bump type from argument
BUMP_TYPE=${1:-patch}

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Auto Release - Version Bump: ${BUMP_TYPE}${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Read current version from VERSION file
if [ ! -f "VERSION" ]; then
    echo -e "${RED}Error: VERSION file not found${NC}"
    exit 1
fi

CURRENT_VERSION=$(cat VERSION)
echo -e "${YELLOW}Current version: ${CURRENT_VERSION}${NC}"

# Split version into components
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

# Bump version based on type
case "$BUMP_TYPE" in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo -e "${RED}Error: Invalid bump type. Use major, minor, or patch${NC}"
        exit 1
        ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
echo -e "${GREEN}New version: ${NEW_VERSION}${NC}"
echo ""

# Update VERSION file
echo "$NEW_VERSION" > VERSION
echo -e "${GREEN}✓ Updated VERSION file${NC}"

# Update CHANGELOG.md
if [ -f "CHANGELOG.md" ]; then
    # Get current date
    RELEASE_DATE=$(date +%Y-%m-%d)

    # Create new changelog entry
    TEMP_FILE=$(mktemp)

    # Read the changelog and insert new version after the header
    awk -v version="$NEW_VERSION" -v date="$RELEASE_DATE" '
        /^# Changelog/ {
            print $0
            print ""
            print "## [" version "] - " date
            print ""
            print "### Changes"
            print ""
            print "- Release created from PR merge"
            print ""
            next
        }
        { print }
    ' CHANGELOG.md > "$TEMP_FILE"

    mv "$TEMP_FILE" CHANGELOG.md
    echo -e "${GREEN}✓ Updated CHANGELOG.md${NC}"
fi

# Commit changes
git add VERSION CHANGELOG.md
git commit -m "chore: bump version to ${NEW_VERSION}"

echo -e "${GREEN}✓ Committed version bump${NC}"

# Create git tag
git tag -a "v${NEW_VERSION}" -m "Release v${NEW_VERSION}"
echo -e "${GREEN}✓ Created tag v${NEW_VERSION}${NC}"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Release v${NEW_VERSION} prepared!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Next steps:"
echo "  git push origin main"
echo "  git push --tags"
echo ""
