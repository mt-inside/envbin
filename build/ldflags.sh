if [ -n "$1" ]
then
    version="$1"
else
    version=$(git describe --tags --dirty)
fi

buildTime=$(date +%s)

echo -X "'"github.com/mt-inside/envbin/pkg/data.Version=${version}"'" -X "'"github.com/mt-inside/envbin/pkg/data.TimeUnix=${buildTime}"'"
