#!/bin/bash

ORG_NAME="go-zen-chu"
TARGET_DIR="${HOME}/repos"
LIMIT=200

repos=$(gh repo list --no-archived --limit ${LIMIT} --json nameWithOwner -q '.[].nameWithOwner')
mkdir -p "${TARGET_DIR}/github.com/${ORG_NAME}"

function clone_repos() {
    local repos=$1
    local dir=$2
    local host=$3

    echo "$repo $dir $host"

    pushd "$dir" || return

    for org_repo in $repos; do
        repo_dir=$(basename "$org_repo")

        if [ -d "$repo_dir" ]; then
            echo "Updating existing repository: $repo_dir"
            pushd "$repo_dir" || return
            git fetch --depth 1 origin
            git reset --hard origin/$(git rev-parse --abbrev-ref HEAD)
            popd || return
        else
            echo "Cloning new repository: $repo_dir"
            git clone --depth 1 "git@${host}:${org_repo}.git" "$repo_dir"
        fi
    done

    popd || return
}

clone_repos "$repos" "${TARGET_DIR}/github.com/${ORG_NAME}" "github.com"
