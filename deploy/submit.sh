
#!/bin/bash -e

CURDIR="$(dirname $0)"
TAG=$(git rev-parse --short HEAD)
echo "TAG: ${TAG}"
export CLOUDSDK_CORE_PROJECT=eoscanada-shared-services

APP=devproxy

pushd ${CURDIR} >/dev/null
  gcloud builds submit \
        --config cloudbuild.yaml \
        --timeout 15m \
        --substitutions=TAG_NAME=${TAG},_APP=${APP} \
        ..
popd >/dev/null
