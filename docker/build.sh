
set -e
BUILD_PATH=`pwd`
docker build -f ./Dockerfile -t gubinempower/sa-backend:20180508 . --no-cache
exit 0
