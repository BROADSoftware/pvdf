#! /bin/bash

VERSION=v0.2.0

UNCOMMITTED="no"

if [ "x$GITHUB_TOKEN" == "x" ]; then
  echo "Please export GITHUB_TOKEN=...."
  exit 1
fi

MYDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BASE=$MYDIR
cd $BASE || exit

cd $BASE/pvscanner/kubernetes || exit
cat <<EOF >./kustomization.yaml
# release.sh generated file. Do not edit
resources:
- deploy.yaml
EOF
kustomize edit set image registry.gitlab.com/pvdf/pvscanner=registry.gitlab.com/pvdf/pvscanner:$VERSION
mkdir -p kustomized
kustomize build . >./kustomized/deploy.yaml

cd $BASE

sh -c "cd pvscanner; make push VERSION=${VERSION}"

TEST=$(git status --porcelain|wc -l)
if [ 0 -ne $TEST -a $UNCOMMITTED != "yes" ]; then
   echo "Please, commit before releasing"
   exit 1
fi

echo "Let's go"
git tag $VERSION

goreleaser --rm-dist release --skip-validate



