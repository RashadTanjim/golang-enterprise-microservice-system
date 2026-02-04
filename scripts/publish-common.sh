#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Function to validate version format
validate_version() {
    if [[ $1 =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
        return 0
    else
        return 1
    fi
}

# Function to increment version
increment_version() {
    local version=$1
    local increment_type=${2:-minor} # patch, minor, major

    # Remove 'v' prefix and extract version parts
    version_clean=$(echo $version | sed 's/^v//')

    # Split version and suffix (e.g., "1.5.0-alpha" -> "1.5.0" and "alpha")
    if [[ $version_clean == *"-"* ]]; then
        version_number=$(echo $version_clean | cut -d'-' -f1)
        version_suffix="-$(echo $version_clean | cut -d'-' -f2-)"
    else
        version_number=$version_clean
        version_suffix=""
    fi

    # Split into major.minor.patch
    IFS='.' read -r major minor patch <<< "$version_number"

    case $increment_type in
        "patch")
            patch=$((patch + 1))
            ;;
        "minor")
            minor=$((minor + 1))
            patch=0
            ;;
        "major")
            major=$((major + 1))
            minor=0
            patch=0
            ;;
    esac

    echo "v${major}.${minor}.${patch}${version_suffix}"
}

MODULE="common"
SPECIFIED_VERSION="$1"

print_info "Common Module Tag Release Script"
echo "=================================="
print_info "Module: $MODULE"

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository!"
    exit 1
fi

# Ensure you are in repo root
print_info "Navigating to repository root..."
cd "$(git rev-parse --show-toplevel)"

# Check if module directory exists
if [ ! -d "$MODULE" ]; then
    print_error "Module directory '$MODULE' not found!"
    exit 1
fi

print_info "Going to module directory to verify go.mod..."

# Go to module directory
cd "$MODULE"

# Verify go.mod exists
if [ ! -f "go.mod" ]; then
    print_error "go.mod not found in $MODULE/ directory!"
    print_info "Make sure this is a valid Go module."
    exit 1
fi

print_success "go.mod verified in $MODULE/ directory."

# Go back to root before continuing
print_info "Returning to repository root for tagging..."
cd ".."

print_info "Fetching latest tags from remote..."

# Fetch all tags from remote to ensure we have the latest
if ! git fetch --tags origin; then
    print_warning "Failed to fetch tags from remote. Proceeding with local tags only..."
    print_warning "This might result in version conflicts if other developers created tags."
else
    print_success "Successfully fetched latest tags from remote."
fi

print_info "Checking for existing $MODULE tags..."

# Get all module tags and find the latest
latest_tag=$(git tag -l "$MODULE/v*" | sort -V | tail -n1)

if [ -z "$latest_tag" ]; then
    print_warning "No existing $MODULE tags found."
    latest_version="v0.0.0"
    echo ""
    print_info "Current latest version: ${latest_version}"

    if [ -z "$SPECIFIED_VERSION" ]; then
        echo ""
        echo "What should be the first version?"
        echo "  1) v1.0.0 - First stable release"
        echo "  2) v0.1.0 - Initial development version"
        echo "  3) Custom version"
        echo ""
        echo -n "Select version (1-3) or press Enter for v1.0.0: "
        read first_version_choice

        case $first_version_choice in
            "2")
                new_version="v0.1.0"
                ;;
            "3")
                echo -n "Enter custom version (e.g., v1.2.3): "
                read custom_version
                if [[ $custom_version != v* ]]; then
                    custom_version="v$custom_version"
                fi
                if ! validate_version "$custom_version"; then
                    print_error "Invalid version format! Use format: v1.2.3 or v1.2.3-alpha"
                    exit 1
                fi
                new_version="$custom_version"
                ;;
            *)
                new_version="v1.0.0"
                ;;
        esac
    else
        new_version="$SPECIFIED_VERSION"
        if [[ $new_version != v* ]]; then
            new_version="v$new_version"
        fi
    fi
else
    print_success "Latest $MODULE tag: $latest_tag"

    # Show recent tags for context
    print_info "Recent $MODULE tags:"
    git tag -l "$MODULE/v*" | sort -V | tail -n5 | while read tag; do
        echo "  - $tag"
    done

    # Extract version from tag (remove module/ prefix)
    latest_version=$(echo $latest_tag | sed "s/$MODULE\///")

    echo ""
    print_info "Current latest version: ${latest_version}"

    if [ -z "$SPECIFIED_VERSION" ]; then
        # Pre-calculate version options
        patch_version=$(increment_version $latest_version "patch")
        minor_version=$(increment_version $latest_version "minor")
        major_version=$(increment_version $latest_version "major")

        # Ask user what type of increment they want
        echo ""
        echo "What type of version increment?"
        echo "  1) Patch (${latest_version} -> ${patch_version}) - Bug fixes"
        echo "  2) Minor (${latest_version} -> ${minor_version}) - New features"
        echo "  3) Major (${latest_version} -> ${major_version}) - Breaking changes"
        echo "  4) Custom version"
        echo ""
        echo -n "Select increment type (1-4) or press Enter for patch: "
        read increment_choice

        case $increment_choice in
            "2")
                new_version="$minor_version"
                ;;
            "3")
                new_version="$major_version"
                ;;
            "4")
                echo -n "Enter custom version (e.g., v1.2.3): "
                read custom_version
                if [[ $custom_version != v* ]]; then
                    custom_version="v$custom_version"
                fi
                if ! validate_version "$custom_version"; then
                    print_error "Invalid version format! Use format: v1.2.3 or v1.2.3-alpha"
                    exit 1
                fi
                new_version="$custom_version"
                ;;
            *)
                new_version="$patch_version"
                ;;
        esac
    else
        new_version="$SPECIFIED_VERSION"
        if [[ $new_version != v* ]]; then
            new_version="v$new_version"
        fi
    fi
fi

# Validate the final version
if ! validate_version "$new_version"; then
    print_error "Invalid version format! Use format: v1.2.3 or v1.2.3-alpha"
    exit 1
fi

# Create the full tag name
new_tag="$MODULE/$new_version"

print_info "Selected version: ${new_version}"
print_info "Full tag: ${new_tag}"

# Check if tag already exists locally
if git rev-parse "$new_tag" >/dev/null 2>&1; then
    print_error "Tag $new_tag already exists locally!"
    exit 1
fi

# Check if tag exists on remote
print_info "Verifying tag doesn't exist on remote..."
if git ls-remote --tags origin | grep -q "refs/tags/$new_tag$"; then
    print_error "Tag $new_tag already exists on remote!"
    print_error "Another developer may have created this tag. Please choose a different version."
    exit 1
fi

echo ""
print_warning "About to create and push tag: $new_tag"
echo -n "Continue? (y/N): "
read confirm

if [[ $confirm != [yY] && $confirm != [yY][eE][sS] ]]; then
    print_info "Operation cancelled."
    exit 0
fi

print_info "Creating tag $new_tag..."

# Create the tag
if ! git tag "$new_tag"; then
    print_error "Failed to create tag!"
    exit 1
fi

print_success "Tag $new_tag created successfully."

print_info "Pushing tag to remote..."

# Push the tag to remote
if ! git push origin "$new_tag"; then
    print_error "Failed to push tag to remote!"
    print_warning "Tag exists locally but not on remote."
    print_info "You can manually push it later with: git push origin $new_tag"
    exit 1
else
    print_success "Tag $new_tag pushed to remote successfully."
fi

echo ""
print_success "Tag creation completed!"
echo "=================================="
print_info "Tag created: $new_tag"
echo ""
print_info "To use this version in other projects:"
echo "  go get github.com/RashadTanjim/enterprise-microservice-system/common@$new_version"
echo ""
print_info "To view tag details:"
echo "  git show $new_tag"
echo ""
print_info "To list all $MODULE tags:"
echo "  git tag -l '$MODULE/v*'"
