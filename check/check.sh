#!/bin/sh
set -e

COMPONENT=$1
BRANCH=$2

echo "Checking docker/${COMPONENT}@${BRANCH}"
(cd docker-ce && git checkout ${BRANCH} &> /dev/null) 
(cd "${COMPONENT}-extract" && git checkout "${BRANCH}-extract-${COMPONENT}" &> /dev/null)
diff -r --exclude=".git" "docker-ce/components/${COMPONENT}" "${COMPONENT}-extract"
echo "OK"