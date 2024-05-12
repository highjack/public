#!/bin/bash
docker run -v "$(pwd)"/complete:/complete -e url=$1 -e pkg_ver=$2 -ti goscan
