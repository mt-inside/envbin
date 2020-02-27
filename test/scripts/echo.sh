source common.sh

set -x
curl -d @- ${base_url}/echo <<< "foo bar"
