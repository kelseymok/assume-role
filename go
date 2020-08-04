#!/usr/bin/env bash

set -e
set -o nounset
set -o pipefail

script_dir=$(cd "$(dirname "$0")" ; pwd -P)
project_root="${script_dir}/."
app_name=${app_name}

goal_build() {
  pushd "${script_dir}" > /dev/null
      docker build -t ${app_name}  .
  popd > /dev/null
}

goal_assume-role() {
  pushd "${script_dir}" > /dev/null
  
    base_profile="${1:-}"
    if [ -z "${base_profile}" ]; then
      echo "Base Profile not supplied. Usage: <BASE PROFILE> <ROLE>"
      exit 1
    fi

    role="${2:-}"
    if [ -z "${role}" ]; then
      echo "Role name not supplied. Usage: <BASE PROFILE> <ROLE>"
      exit 1
    fi

    mounted_dir=$(cd ${PROJECT_ROOT}/..; pwd)
    home_dir=$(cd ~; pwd)
    echo "Mounting ${mounted_dir}"

    if [ ! -z "$(docker ps -a | grep ${app_name})" ]; then
      docker rm -f ${app_name}
    fi

    docker run -it \
      --name ${app_name} \
      -v "${mounted_dir}:/app" \
      -v "${home_dir}/.aws:/root/.aws" \
      --entrypoint="" \
      ${app_name} assume-role ${base_profile} ${role}
  popd > /dev/null
}

goal_run() {
  pushd "${SCRIPT_DIR}" > /dev/null
    mounted_dir=$(cd ${PROJECT_ROOT}/..; pwd)
    home_dir=$(cd ~; pwd)
    echo "Mounting ${mounted_dir}"

    if [ ! -z "$(docker ps -a | grep ${app_name})" ]; then
      docker rm -f ${app_name}
    fi


    docker run -it \
      --name ${app_name} \
      -v "${mounted_dir}:/app" \
      -v "${home_dir}/.aws:/root/.aws" \
      ${app_name}
  popd > /dev/null
}

cleanup-containers() {
  container=$1
  if [ -z "${container}" ]; then
    echo "You provided an empty container. Usage <func> <container-name>"
    exit 1
  fi

  if [ ! -z "$(docker ps -a | grep $container)" ]; then
    docker rm -f $container
  fi
}

TARGET=${1:-}
if type -t "goal_${TARGET}" &>/dev/null; then
  "goal_${TARGET}" ${@:2}
else
  echo "Usage: $0 <goal>

goal:
    build                   - Builds container
    run                     - Runs container
    assume-role             - Assumes Role from argument and updates Shared Credentials with valid session
"
  exit 1
fi
