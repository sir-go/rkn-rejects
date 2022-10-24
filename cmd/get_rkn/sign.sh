#!/bin/bash
/opt/cprocsp/bin/amd64/csptest -sfsign -sign -in /tmp/req.xml \
-out /tmp/req.xml.signed -my "${COMPANY_NAME}" -detached -add
