#!/bin/sh


COMPONENT=$1
TAG=$2
COMMIT=$3

#echo "Checking docker/${COMPONENT}@${TAG}"
(cd docker-ce && git checkout ${TAG} &> /dev/null) 
(cd "${COMPONENT}-extract" && git checkout ${COMMIT} &> /dev/null)
diff -r --exclude=".git" "docker-ce/components/${COMPONENT}" "${COMPONENT}-extract"