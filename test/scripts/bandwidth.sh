source common.sh

set -x
curl ${base_url}/handlers/bandwidth?value=1
time curl -N ${base_url}
curl ${base_url}/handlers/bandwidth?value=100
time curl -N ${base_url}
curl ${base_url}/handlers/bandwidth?value=0
