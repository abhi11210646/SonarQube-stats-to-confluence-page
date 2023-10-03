#!/bin/bash

# Generate kubernetes yaml output from template files and environment files.
#
# Designed to be called from the project root.
#
# $ ./kubernetes/generate-yaml.sh
#
# Arguments are passed through environment variables. This is convenient when
# the script is used on GitLab CI. See the example below:
#
# ```yaml
# job:
#   variables:
#     RELEASE_ENVIRONMENT: test
#   script:
#     - ./kubernetes/generate-yaml.sh
# ```


# Bash strict mode (http://redsymbol.net/articles/unofficial-bash-strict-mode/)
set -euo pipefail
IFS=$'\n\t'

##############################################################################
# Helper functions

function exit_with_message() {
  MESSAGE=$1
  CODE=$2
  (>&2 echo $MESSAGE; exit $CODE)
}

###############################################################################
# VARIABLES
#
# Variables intended for use in the template files should be exported. Ones
# that are only used within this script need not be exported.

export APP_NAME=devin-dashboard
GITLAB_ORG_NAME=wp-in

# Used to locate the right environment file.
export RELEASE_ENVIRONMENT=${RELEASE_ENVIRONMENT:?Required ENV var}

# DEPLOY_SHA falls back to the value of CI_COMMIT_SHORT_SHA from GitLab CI if
# not passed. Either value must be available. This is the tag/sha of the docker
# image which will be deployed.
if [ -n "${DEPLOY_SHA:=}" ]; then
  DEPLOY_SHA=$DEPLOY_SHA
elif [ -n "${CI_COMMIT_SHORT_SHA:=}" ]; then
  DEPLOY_SHA=$CI_COMMIT_SHORT_SHA
else
  exit_with_message "Either CI_COMMIT_SHORT_SHA or DEPLOY_SHA must be set" 1
fi

# The registry prefix, used in the docker image name
# HARBOR_REGISTRY_PREFIX is set by systems on most projects in GitLab.
HARBOR_REGISTRY_PREFIX=${HARBOR_REGISTRY_PREFIX:-harbor.one.com/$GITLAB_ORG_NAME}
REGISTRY_PREFIX=${REGISTRY_PREFIX:-$HARBOR_REGISTRY_PREFIX}

# The name of the docker image
IMAGE_NAME=$REGISTRY_PREFIX/$APP_NAME

# The image to be deployed
export IMAGE=$IMAGE_NAME:${DEPLOY_SHA}

# The version number of the release, or just a normal git ref
export RELEASE_VERSION=${CI_COMMIT_REF_NAME:-n/a}

# Assigning default values for optional vars in environment files
export INGRESS_CLASS=nginx

# If we need to know what environment we are deploying to, it can be done by
# reading the GitLab CI `CI_ENVIRONMENT_NAME` variable.
# ENVIRONMENT=${CI_ENVIRONMENT_NAME:?Missing CI_ENVIRONMENT_NAME}


##############################################################################
# Load variables from environment file.

ENVFILE=./kubernetes/environments/${RELEASE_ENVIRONMENT}.env

if [[ -e $ENVFILE ]] ; then
  set -o allexport
  source $ENVFILE
  set +o allexport
else
  exit_with_message "Could not find env file: ${ENVFILE}" 1
fi

##############################################################################
# Invariant checks

# Guard against releasing to unexpected environments on CI by accident
if [ "${CI_ENVIRONMENT_NAME:=}" != "" ] ; then
  if [ "${K8S_ENV}" != "" ] ; then
    if [ "${K8S_ENV}" != "${CI_ENVIRONMENT_NAME}" ] ; then
        exit_with_message "Attempting to release to '${CI_ENVIRONMENT_NAME}' but expects '${K8S_ENV}'" 1
    fi
  fi
fi

##############################################################################
# Load the templates and substitute the values.

for file in $(find kubernetes/template -type f -name '*.yaml') ; do
  # Write the yaml separator if it's not already in the file.
  [[ $(head -c 3 $file) != "---" ]] && echo '---'  ;
  envsubst < $file ;
done
