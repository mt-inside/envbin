source common.sh

set -x
curl --header "Accept: application/json" ${base_url} | jq
