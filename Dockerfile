FROM alpine:3.4

COPY . /source
ADD config.json.template /config.json
ADD startup.sh /

RUN apk update \
  && apk add bash git bzr go ca-certificates \
  && cd /source/ \
  && BUILDINFO_PACKAGE="github.com/Financial-Times/service-status-go/buildinfo." \
  && VERSION="version=$(git describe --tag --always 2> /dev/null)" \
  && DATETIME="dateTime=$(date -u +%Y%m%d%H%M%S)" \
  && REPOSITORY="repository=$(git config --get remote.origin.url)" \
  && REVISION="revision=$(git rev-parse HEAD)" \
  && BUILDER="builder=$(go version)" \
  && LDFLAGS="-X '"${BUILDINFO_PACKAGE}$VERSION"' -X '"${BUILDINFO_PACKAGE}$DATETIME"' -X '"${BUILDINFO_PACKAGE}$REPOSITORY"' -X '"${BUILDINFO_PACKAGE}$REVISION"' -X '"${BUILDINFO_PACKAGE}$BUILDER"'" \
  && cd - \
  && export GOPATH=/gopath \
  && REPO_PATH="github.com/Financial-Times/publish-availability-monitor" \
  && mkdir -p $GOPATH/src/${REPO_PATH} \
  && mv /source/* $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get -t -d -v ./... \
  && go build -ldflags="${LDFLAGS}" \
  && go test ./... \
  && mv publish-availability-monitor / \
  && apk del go git bzr \
  && rm -rf /source $GOPATH /var/cache/apk/*

CMD [ "/startup.sh" ]
