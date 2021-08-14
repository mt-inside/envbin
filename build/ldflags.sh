version=$(git describe --tags --dirty)
buildTime=$(date +%s)
echo -X "'"github.com/mt-inside/envbin/pkg/data.Version=${version}"'" -X "'"github.com/mt-inside/envbin/pkg/data.TimeUnix=${buildTime}"'"
