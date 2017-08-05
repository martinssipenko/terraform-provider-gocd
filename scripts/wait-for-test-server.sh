#! /bin/bash -x

function get_status {
    curl -H 'Accept: application/vnd.go.cd.v3+json' \
        --write-out %{http_code} \
        --silent \
        --output /dev/null \
        http://localhost:8153/go/api/admin/templates
}

counter=0
wait_length=5

while [ $counter -lt 30 ]; do

    code=$(get_status)
    if [ "200" == "$code" ]; then
        echo "Got status ${code}. Exiting..."
        exit 0
    fi

    echo "Got status ${code}. Waiting for '${wait_length}' seconds..."

    sleep "${wait_length}"

done

exit 1