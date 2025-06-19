#!/bin/bash

# Version bumping script for semantic versioning
# Uses git tags to manage versions with semver compatibility

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}INFO:${NC} $1"
}

print_success() {
    echo -e "${GREEN}SUCCESS:${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}WARNING:${NC} $1"
}

print_error() {
    echo -e "${RED}ERROR:${NC} $1"
}

# Function to get the current version from git tags
get_current_version() {
    local current_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [[ -z "$current_tag" ]]; then
        echo "0.0.0"
    else
        echo "$current_tag" | sed 's/^v//'
    fi
}

# Function to parse semantic version
parse_version() {
    local version=$1
    local major=$(echo "$version" | cut -d. -f1)
    local minor=$(echo "$version" | cut -d. -f2)
    local patch=$(echo "$version" | cut -d. -f3)
    echo "$major $minor $patch"
}

# Function to calculate next version
calculate_next_version() {
    local type=$1
    local current=$2
    
    read -r major minor patch <<< $(parse_version "$current")
    
    case $type in
        "major")
            echo "$((major + 1)).0.0"
            ;;
        "minor")
            echo "$major.$((minor + 1)).0"
            ;;
        "patch")
            echo "$major.$minor.$((patch + 1))"
            ;;
        *)
            echo "Invalid version type: $type" >&2
            exit 1
            ;;
    esac
}

# Function to check if working directory is clean
check_git_status() {
    if [[ -n $(git status --porcelain) ]]; then
        print_error "Working directory is not clean. Please commit or stash your changes."
        git status --short
        exit 1
    fi
}

# Function to check version consistency
check_version_consistency() {
    local version=$1
    local warnings=0
    
    print_info "Checking version consistency..."
    
    # Check that GoReleaser will handle the version correctly
    if [[ ! -f "meta/version.go" ]]; then
        print_warning "meta/version.go not found - version may not be set correctly"
        warnings=$((warnings + 1))
    fi
    
    if [[ $warnings -gt 0 ]]; then
        print_warning "Found $warnings potential version issues"
        print_info "The version will be set by GoReleaser when the release is built"
    else
        print_success "Version consistency checks passed"
    fi
}

# Function to create and push tag
create_tag() {
    local version=$1
    local tag="v$version"
    
    # Check version consistency before creating tag
    check_version_consistency "$version"
    
    print_info "Creating tag: $tag"
    git tag -a "$tag" -m "Release $tag"
    
    print_info "Pushing tag to origin..."
    git push origin "$tag"
    
    print_success "Tag $tag created and pushed successfully!"
    print_info "GoReleaser will set the actual version in the binary when the release is built"
}

# Function to show version information
show_version() {
    local type=${1:-"patch"}
    local current=$(get_current_version)
    local next=$(calculate_next_version "$type" "$current")
    echo "$next"
}

# Function to bump version
bump_version() {
    local type=${1:-"patch"}
    
    # Check git status
    check_git_status
    
    # Get current version
    local current=$(get_current_version)
    print_info "Current version: $current"
    
    # Calculate next version
    local next=$(calculate_next_version "$type" "$current")
    print_info "Next version: $next"
    
    # Confirm with user
    echo -n "Are you sure you want to bump from $current to $next? [y/N]: "
    read -r confirm
    
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        print_warning "Version bump cancelled."
        exit 0
    fi
    
    # Create and push tag
    create_tag "$next"
    
    print_success "Version bumped from $current to $next"
    print_info "To create a release, run: make release"
}

# Function to show help
show_help() {
    cat << EOF
Usage: $0 [COMMAND] [TYPE]

COMMANDS:
    patch       Bump patch version (default)
    minor       Bump minor version  
    major       Bump major version
    show        Show next version without creating tag
    help        Show this help message

EXAMPLES:
    $0 patch              # Bump patch version (1.2.3 → 1.2.4)
    $0 minor              # Bump minor version (1.2.3 → 1.3.0)  
    $0 major              # Bump major version (1.2.3 → 2.0.0)
    $0 show patch         # Show next patch version
    $0 show minor         # Show next minor version
    $0 show major         # Show next major version

This script manages semantic versioning using git tags.

EOF
}

# Main script logic
main() {
    local command=${1:-"patch"}
    local type=${2:-"patch"}
    
    case $command in
        "patch"|"minor"|"major")
            bump_version "$command"
            ;;
        "show")
            show_version "$type"
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Check if git is available
if ! command -v git &> /dev/null; then
    print_error "git is required but not installed."
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository."
    exit 1
fi

# Run main function
main "$@"