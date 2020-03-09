#!/bin/sh

function check() {
    docker run --rm slubecki/check-extract cli $1
    docker run --rm slubecki/check-extract engine $1
    docker run --rm slubecki/check-extract packaging $1
}

check "17.06"
check "17.07"
check "17.09"
check "17.10"
check "17.11"
check "17.12"
check "18.01"
check "18.02"
check "18.03"
check "18.04"
check "18.05"
