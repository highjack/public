#!/bin/bash

if [ "$#" -ne 2 ]; then
        echo "./scan <package-name> <version>"
        exit
fi

docker run -ti -v $(pwd)/results:/results -e PACKAGE=$1 -e VERSION=$2  ps-scan
