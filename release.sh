#! /bin/bash

VERSION=v0.1.0

UNCOMMITTED="yes"


if [ "x$GITHUB_TOKEN" == "x" ]; then
  echo "Please export GITHUB_TOKEN=...."
  exit 1
fi

MYDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BASE=$MYDIR
cd $BASE || exit



cd $BASE/pvscanner/kubernetes || exit
cat <<EOF >./kustomization.yaml
resources:
- deploy.yaml
EOF
kustomize edit set image pvdf/pvscanner=pvdf/pvscanner:$VERSION
mkdir -p kustomized
kustomize build . >./kustomized/deploy.yaml


cd $BASE

sh -c "cd vgsd; make build VERSION=${VERSION}"
sh -c "cd pvscanner; make push VERSION=${VERSION}"



TEST=$(git status --porcelain|wc -l)
if [ 0 -ne $TEST -a $UNCOMMITTED != "yes" ]; then
   echo "Please, commit before releasing"
   exit 1
fi

echo "Let's go"
git tag $VERSION

goreleaser --rm-dist release --skip-validate



