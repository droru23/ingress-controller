#!/bin/bash
if [ -z "$GOPATH" ]; then
	export GOPATH=$HOME/go
fi

# ROOT_PACKAGE :: the package (relative to $GOPATH/src) that is the target for code generation
ROOT_PACKAGE="github.wdf.sap.corp/i531295/Ingress-Controller"
# CUSTOM_RESOURCE_NAME :: the name of the custom resource that we're generating client code for
CUSTOM_RESOURCE_NAME="SimpleIngress"
# CUSTOM_RESOURCE_VERSION :: the version of the resource
CUSTOM_RESOURCE_VERSION="v1"

# retrieve the code-generator scripts and bins
CODE_GEN_VER=v0.26.8
go get -u k8s.io/code-generator/...@$CODE_GEN_VER

# run the code-generator entrypoint script
$GOPATH/pkg/mod/k8s.io/code-generator@$CODE_GEN_VER/generate-groups.sh all "$ROOT_PACKAGE/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION" \
 --go-header-file "$GOPATH/pkg/mod/k8s.io/code-generator@$CODE_GEN_VER"/hack/boilerplate.go.txt

