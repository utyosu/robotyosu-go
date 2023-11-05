REQUIRED_GO_VERSION=1.21.3

LOCAL_GO_VERSION=`go version | cut -c 14- | cut -d' ' -f1`
if [ "${LOCAL_GO_VERSION}" != "${REQUIRED_GO_VERSION}" ]; then
  echo "Error: Please update golang (local=${LOCAL_GO_VERSION}, required=${REQUIRED_GO_VERSION})"
  exit 1
fi
