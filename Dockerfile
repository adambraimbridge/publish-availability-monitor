FROM alpine:3.4

COPY . /app
ADD config.json.template /config.json
ADD startup.sh /

RUN apk update \
  && apk add bash git bzr go ca-certificates \
  && export GOPATH=/gopath \
  && REPO_PATH="github.com/Financial-Times/publish-availability-monitor" \
  && mkdir -p $GOPATH/src/${REPO_PATH} \
  && mv /app/* $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get -t -d -v ./... \
  && go build \
  && go test ./... \
  && mv publish-availability-monitor / \
  && apk del go git bzr \
  && rm -rf /app $GOPATH /var/cache/apk/*

CMD [ "/startup.sh" ]
