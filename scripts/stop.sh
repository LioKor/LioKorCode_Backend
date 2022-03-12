# shellcheck disable=SC2006

id1=`pgrep main_service`

kill $id1

rm -rf build/