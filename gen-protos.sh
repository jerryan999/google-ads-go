# Args from make
PROTO_ROOT_DIR=$1
PROTO_SRC_DIR=$2
PROTOC_GO_ARGS=$3

GOOGLE_ADS_VERSION=v15

# Protos in order to compile
PROTOS_ENUMS=googleapis/google/ads/googleads/$GOOGLE_ADS_VERSION/enums/*.proto
PROTOS_ERRORS=googleapis/google/ads/googleads/$GOOGLE_ADS_VERSION/errors/*.proto
PROTOS_COMMON=googleapis/google/ads/googleads/$GOOGLE_ADS_VERSION/common/*.proto
PROTOS_RESOURCES=googleapis/google/ads/googleads/$GOOGLE_ADS_VERSION/resources/*.proto
PROTOS_SERVICES=googleapis/google/ads/googleads/$GOOGLE_ADS_VERSION/services/*.proto

function compile_protos() {
    PROTOS=$*
    for file in $PROTOS; do
        echo "converting proto $(basename $file)"
        FOLDER_PATH=$(dirname $file)
        PACKAGE=$(basename $FOLDER_PATH)
        protoc -I=$PROTO_ROOT_DIR $PROTOC_GO_ARGS $file
    done
}

compile_protos $PROTOS_ENUMS
compile_protos $PROTOS_ERRORS
compile_protos $PROTOS_COMMON
compile_protos $PROTOS_RESOURCES
compile_protos $PROTOS_SERVICES

cd googleapis && git checkout . && cd ../