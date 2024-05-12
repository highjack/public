#!/bin/bash
docker build . -t ps-scan
mkdir results 2>/dev/null
chmod +x ./ps-scan.sh
