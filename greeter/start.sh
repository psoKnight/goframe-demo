#!/bin/sh
set -eu

basepath=$(pwd)
confpath="${basepath}/conf/config.yml"

exec ./greeter -c="${confpath}"
