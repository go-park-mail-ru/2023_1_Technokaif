#!/bin/bash

BASEDIR=$(dirname $(realpath "$0"))

$BASEDIR/build_api.sh
$BASEDIR/build_auth.sh
$BASEDIR/build_search.sh
$BASEDIR/build_user.sh
