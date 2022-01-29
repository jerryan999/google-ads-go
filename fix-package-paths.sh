PACKAGES=('common' 'enums' 'errors' 'resources' 'services')
GOOGLE_ADS_VERSION=v9

function fix_package_path() {
    FILE=$1
    PACKAGE=$2
    MATCH="google.golang.org\/genproto\/googleapis\/ads\/googleads\/$GOOGLE_ADS_VERSION\/"
    REPLACE="github.com\/butlerhq\/google-ads-go\/"
    sed -i "" "s/$MATCH$PACKAGE/$REPLACE$PACKAGE/g" $FILE
}

function fix_package_name() {
    FILE=$1
    PACKAGE=$2
    sed -i "" "s/google_ads_googleads_$GOOGLE_ADS_VERSION\_$PACKAGE/$PACKAGE/g" $FILE
}

echo "fixing packages"
for file in ./google/ads/googleads/$GOOGLE_ADS_VERSION/**/*.pb.go; do
    for p in "${PACKAGES[@]}"; do
        fix_package_path $file $p
        fix_package_name $file $p
    done
done
mv ./google/ads/googleads/$GOOGLE_ADS_VERSION/* ./
echo "finished fixing packages"