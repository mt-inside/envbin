source common.sh

set -x
curl ${base_url}/api/allocate?value=1073741824
