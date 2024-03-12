#!/bin/bash
set -eou pipefail

# Note: using '-i.bak' to support different versions of sed when using in-place editing.

# Swap tag in for main for URLs if the version is vx.x.x*
if [[ "${APP_VERSION}" =~ ^v[0-9]+.[0-9]+.[0-9]+(.*)? ]]; then
  sed -i.bak -e "s%spinkube/spin-operator/main%spinkube/spin-operator/${APP_VERSION}%g" "${STAGING_DIR}/${CHART_NAME}-${CHART_VERSION}/README.md"
  sed -i.bak -e "s%spinkube/spin-operator/main%spinkube/spin-operator/${APP_VERSION}%g" "${STAGING_DIR}/${CHART_NAME}-${CHART_VERSION}/templates/NOTES.txt"
fi

## Update Chart.yaml with CHART_VERSION and APP_VERSION
sed -r -i.bak -e "s%^version: .*%version: ${CHART_VERSION}%g" "${STAGING_DIR}/${CHART_NAME}-${CHART_VERSION}/Chart.yaml"
sed -r -i.bak -e "s%^appVersion: .*%appVersion: ${APP_VERSION}%g" "${STAGING_DIR}/${CHART_NAME}-${CHART_VERSION}/Chart.yaml"

## Update README.md with CHART_VERSION
sed -i.bak -e "s%{{ CHART_VERSION }}%${CHART_VERSION}%g" "${STAGING_DIR}/${CHART_NAME}-${CHART_VERSION}/README.md"

# Cleanup
find "${STAGING_DIR}/${CHART_NAME}-${CHART_VERSION}" -type f -name '*.bak' -print0 | xargs -0 rm --