#!/bin/sh
set -e

COMPONENT=$1
TAG=$2
COMMIT=$3

#echo "Checking docker/${COMPONENT}@${TAG}"
(cd docker-ce && git checkout ${TAG} > /dev/null 2>&1 ) 
(cd "${COMPONENT}-extract" && git checkout ${COMMIT} > /dev/null 2>&1 )
diff -r --exclude=".git" "docker-ce/components/${COMPONENT}" "${COMPONENT}-extract"
echo "OK"