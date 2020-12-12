version=$(git describe --tags --dirty)
gitCommit=$(git rev-parse HEAD)
buildTime=$(date +%Y-%m-%d_%H:%M:%S%z)
echo "-X github.com/mt-inside/envbin/pkg/data.Version=${version} -X github.com/mt-inside/envbin/pkg/data.GitCommit=${gitCommit} -X github.com/mt-inside/envbin/pkg/data.BuildTime=${buildTime}"
