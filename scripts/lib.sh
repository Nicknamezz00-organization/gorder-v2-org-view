#!/usr/bin/env bash

set -euo pipefail

ROOT_MODULE="$PWD"

function set_root_dir {
  ROOT_DIR=$ROOT_MODULE
#  echo set to $ROOT_DIR
}
set_root_dir

COLOR_RED='\033[0;31m'
COLOR_ORANGE='\033[0;33m'
COLOR_GREEN='\033[0;32m'
COLOR_LIGHTCYAN='\033[0;36m'
COLOR_BLUE='\033[0;94m'
COLOR_MAGENTA='\033[95m'
COLOR_BOLD='\033[1m'
COLOR_NONE='\033[0m' # No Color


function log_error {
  >&2 echo -n -e "${COLOR_BOLD}${COLOR_RED}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_warning {
  >&2 echo -n -e "${COLOR_ORANGE}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_callout {
  >&2 echo -n -e "${COLOR_LIGHTCYAN}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_cmd {
  >&2 echo -n -e "${COLOR_BLUE}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_success {
  >&2 echo -n -e "${COLOR_GREEN}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function log_info {
  >&2 echo -n -e "${COLOR_NONE}"
  >&2 echo "$@"
  >&2 echo -n -e "${COLOR_NONE}"
}

function modules() {
  modules=$(ls internal)
  echo "${modules[@]}"
}

function modules_exp() {
  for m in $(modules); do
    echo -n "${m}/... "
  done
}

# From http://stackoverflow.com/a/12498485
function relativePath {
  # both $1 and $2 are absolute paths beginning with /
  # returns relative path to $2 from $1
  local source=$1
  local target=$2

  local commonPart=$source
  local result=""

  while [[ "${target#"$commonPart"}" == "${target}" ]]; do
    # no match, means that candidate common part is not correct
    # go up one level (reduce common part)
    commonPart="$(dirname "$commonPart")"
    # and record that we went back, with correct / handling
    if [[ -z $result ]]; then
      result=".."
    else
      result="../$result"
    fi
  done

  if [[ $commonPart == "/" ]]; then
    # special case for root (no common path)
    result="$result/"
  fi

  # since we now have identified the common part,
  # compute the non-common part
  local forwardPart="${target#"$commonPart"}"

  # and now stick all parts together
  if [[ -n $result ]] && [[ -n $forwardPart ]]; then
    result="$result$forwardPart"
  elif [[ -n $forwardPart ]]; then
    # extra slash removal
    result="${forwardPart:1}"
  fi

  echo "$result"
}

function module_dirs() {
  echo "internal/common internal/order internal/stock internal/payment"
}

function module_subdir {
  relativePath "${ROOT_DIR}" "${PWD}"
}

####    Running actions against multiple modules ####

# run [command...] - runs given command, printing it first and
# again if it failed (in RED). Use to wrap important test commands
# that user might want to re-execute to shorten the feedback loop when fixing
# the test.
function run {
  local rpath
  local command
  rpath=$(module_subdir)
  # Quoting all components as the commands are fully copy-parsable:
  command=("${@}")
  command=("${command[@]@Q}")
  if [[ "${rpath}" != "." && "${rpath}" != "" ]]; then
    repro="(cd ${rpath} && ${command[*]})"
  else
    repro="${command[*]}"
  fi

  log_cmd "% ${repro}"
  "${@}" 2> >(while read -r line; do echo -e "${COLOR_NONE}stderr: ${COLOR_MAGENTA}${line}${COLOR_NONE}">&2; done)
  local error_code=$?
  if [ ${error_code} -ne 0 ]; then
    log_error -e "FAIL: (code:${error_code}):\\n  % ${repro}"
    return ${error_code}
  fi
}

# receives a directory, relative to ROOT_DIR. If it exists, delete all files under it,
# otherwise create that directory
function prepare_dir() {
  local dir="$1"
  if [ -d "$dir" ]; then
    log_warning "Directory $dir exists. Deleting all files under $dir."
    run find "$dir" -mindepth 1 -delete
  else
    log_callout "Directory $dir does not exist. Creating directory $dir."
    # Create the directory
    run mkdir -p "$dir"
  fi
}

# run_for_module [module] [cmd]
# executes given command in the given module for given pkgs.
#   module_name - "." (in future: tests, client, server)
#   cmd         - cmd to be executed - that takes package as last argument
function run_for_module {
  local module=${1:-"."}
  shift 1
  (
    cd "${ROOT_DIR}/${module}" && "$@"
  )
}

#  run_for_modules [cmd]
#  run given command across all modules and packages
#  (unless the set is limited using ${PKG} or / ${USERMOD})
function run_for_modules {
  KEEP_GOING_MODULE=${KEEP_GOING_MODULE:-false}
  local pkg="${PKG:-./...}"
  log_info "pkg = $pkg"
  local fail_mod=false
  if [ -z "${USERMOD:-}" ]; then
    for m in $(module_dirs); do
      if run_for_module "${m}" "$@" "${pkg}"; then
        continue
      else
        if [ "$KEEP_GOING_MODULE" = false ]; then
          log_error "There was a Failure in module ${m}, aborting..."
          return 1
        fi
        log_error "There was a Failure in module ${m}, keep going..."
        fail_mod=true
      fi
    done
    if [ "$fail_mod" = true ]; then
      return 1
    fi
  else
    run_for_module "${USERMOD}" "$@" "${pkg}" || return "$?"
  fi
}

# generic_checker [cmd...]
# executes given command in the current module, and clearly fails if it
# failed or returned output.
function generic_checker {
  local cmd=("$@")
  if ! output=$("${cmd[@]}"); then
    echo "${output}"
    log_error -e "FAIL: '${cmd[*]}' checking failed (!=0 return code)"
    return 255
  fi
  if [ -n "${output}" ]; then
    echo "${output}"
    log_error -e "FAIL: '${cmd[*]}' checking failed (printed output)"
    return 255
  fi
}

echo "go modules: $(modules)"