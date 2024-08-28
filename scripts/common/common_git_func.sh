#!/usr/bin/env bash -x

##############################
# Internal Utility Functions #
##############################

function setup_prompt_git()
{
   status "Setting up Git prompt..."

   local        BLUE="\[\033[0;34m\]"
   local        CYAN="\[\033[0;36m\]"
   local         RED="\[\033[0;31m\]"
   local   LIGHT_RED="\[\033[1;31m\]"
   local       GREEN="\[\033[0;32m\]"
   local LIGHT_GREEN="\[\033[1;32m\]"
   local       WHITE="\[\033[1;37m\]"
   local  LIGHT_GRAY="\[\033[0;37m\]"
   local     DEFAULT="\[\033[0m\]"
   PS1="\u@\h:[\W]$CYAN [\$(parse_git_branch_prompt)] $DEFAULT\$ "

   local prompt_file="$1"
   [ "$prompt_file" = "" ] && prompt_file=~/.system_jenkins_custom_prompt

   if [ -f ${prompt_file} ]
   then
      export PS1=`grep -v "#" $prompt_file`
   else
      # Default JENKINS project prompt:
      export PS1="\u@\h:[\W]$CYAN [\$(parse_git_branch_prompt)] $DEFAULT\$ "
  fi
  status "Done. Setup Git prompt"
}

function setup_git()
{
   status "Setting up git..."
   local repository_path=$1

   # Bash completion for git
   . ${repository_path}/env/git/completion.bash
   . ${repository_path}/env/git/config

   ####################
   # installing hooks #
   ####################

   # git for appending Change-Id to commit logs (required for Gerrit)
   mkdir -p ${repository_path}/.git/hooks
   if [ ! -f "${repository_path}/.git/hooks/commit-msg" ]; then
      if [ ! -f ${repository_path}/env/git/hooks/commit-msg ]; then
         scp -q -p -P 29418 ${ERICSSON_SIGNUM}@gerrit.ericsson.se:hooks/commit-msg ${repository_path}/env/git/hooks/
      fi
   fi
   # UDM JENKINS HOOKS
   cp -f ${repository_path}/env/git/hooks/* ${repository_path}/.git/hooks
   status "Done. Setup git"
}

function setup_bash_git()
{
   status "Setting up Git bash..."
   local repository_path=$1
   alias git-push="$(echo ${repository_path})/env/git/git-push.bash"
   alias git-push-current-brach="$(echo ${repository_path})/env/git/git-push.bash $(parse_git_branch_prompt)"
   alias git-push-draft="$(echo ${repository_path})/env/git/git-push-draft.bash"
   alias git-push-topic="$(echo ${repository_path})/env/git/git-push-topic.bash"
   status "Done Setup Git bash..."
}