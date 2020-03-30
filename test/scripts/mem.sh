source common.sh

set -x
curl ${base_url}/handlers/allocate?value=1073741824
