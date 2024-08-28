#!/bin/bash -x

#########################################################
# Utility Functions to export outside the setup.sh scope
#########################################################

function setup_root_project_dir_git()
{
   export GIT_ROOT="$(git rev-parse --show-toplevel 2> /dev/null)"
}

# main setup function
function setup_project()
{
   status "Setting Git Environment..."
   set_signum_ericsson
   setup_git ${GIT_ROOT}
   setup_prompt_git
   setup_bash_git ${GIT_ROOT}
   status "Done. Set Git Environment."
}

function init_setup_project_git()
{
   local repository_path=$1
   . ${repository_path}/scripts/common/common_shell_func.sh
   . ${repository_path}/scripts/common/common_git_func.sh
}

#############
# EXECUTION #
#############

# Check own shell execution (preffixed with dot or source)
[[ x"${BASH_SOURCE[0]}" = x"$0" ]] && { echo "This script must be executed 'sourced' !"; exit 1; }

setup_root_project_dir_git
init_setup_project_git ${GIT_ROOT}
setup_root_project_dir
setup_project

info ""
status "Done. Git Enviroment finished"
info "\tGIT aliases               : `git config --get-regexp '^alias\.' | sed 's/.*\.\([^ ]*\) .*/\1/' | tr '\n' ' '`"
info ""

unset GIT_ROOT