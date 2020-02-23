source common.sh

set -x
curl -X POST ${base_url}/api/bandwidth?value=1
time curl -N ${base_url}
curl -X POST ${base_url}/api/bandwidth?value=100
time curl -N ${base_url}
curl -X POST ${base_url}/api/bandwidth?value=0
