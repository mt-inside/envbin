source common.sh

set -x
curl -X POST ${base_url}/api/allocate?value=1000000000
