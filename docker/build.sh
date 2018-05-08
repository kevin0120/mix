
set -e
BUILD_PATH=`pwd`
docker build --build-arg CODING_PWD=$CODING_PWD CODING_USER=$CODING_USER GITHUB_PWD=$GITHUB_PWD GITHUB_USER=$GITHUB_USER -f ./Dockerfile -t gubinempower/sa-backend:20180508 . --no-cache
exit 0
