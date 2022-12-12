#!/bin/bash

# this is a requests signing command line (using a CryptoPRO CLI tool)

/opt/cprocsp/bin/amd64/csptest -sfsign -sign -in /tmp/req.xml \
-out /tmp/req.xml.signed -my "${COMPANY_NAME}" -detached -add
