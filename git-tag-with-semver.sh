#!/bin/bash

this_script=${0##*/}
    
usage="
Usage: $this_script [-f|--force]

Apply SEMVER tag locally and to the remote. Requires
a clean git status. Use --force to reapply the tag.

"


force=
case "$1" in
    -h|--h*)
        echo >&2 "$usage"
        exit 1
        ;;
    -f|--force)
        force="--force"
        shift
        ;;
esac

# Common functions, etc:
case "$0" in
    /*|*/*) this_dir=$(dirname $(realpath $0));;
    *) this_dir=$(dirname $(realpath $(which $0)));;
esac
project_root_dir=$this_dir

set -e
set -x; cd $project_root_dir; set +x
export PATH="$(realpath $this_dir)${PATH+:}${PATH}"

# Must have semver:
semver=$(cat semver.txt)
if [[ "$semver" != v* ]]; then
    semver="v$semver"
fi

# Must be in in proper git state:
if ! check-git-state.sh; then
    echo >&2 "$this_script: cannot continue"
    exit 1
fi

# The git tag must include the sub-directory path to the tag:
git tag $force $semver testutils/$semver 
git push $force origin tag $semver testutils/$semver


 

