version=$(git describe --tags --dirty)
buildTime=$(date +"%F %T%Z")
echo -X "'"github.com/mt-inside/envbin/pkg/data.Version=${version}"'" -X "'"github.com/mt-inside/envbin/pkg/data.BuildTime=${buildTime}"'"
