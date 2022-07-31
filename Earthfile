VERSION 0.5
ARG IMAGE_REPO="ghcr.io/cartermckinnon"

proto-builder:
    # toolchain last updated: April 16, 2022.
    FROM ubuntu:22.04
    # Get rid of the warning: "debconf: unable to initialize frontend: Dialog"
    # https://github.com/moby/moby/issues/27988
    RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
    RUN apt-get update && apt-get install wget unzip golang git -y
    # https://github.com/protocolbuffers/protobuf/releases
    WORKDIR /tmp
    RUN wget -O protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.20.0/protoc-3.20.0-linux-x86_64.zip && \
        unzip protoc.zip -d /protoc
    ENV PATH=$PATH:/protoc/bin
    # https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go?tab=versions
    RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
    ENV PATH=$PATH:/root/go/bin
    SAVE IMAGE --push $IMAGE_REPO/filter-feed/proto-builder:earthly-cache

proto:
    FROM +proto-builder
    WORKDIR /proto
    COPY proto/*.proto .
    RUN mkdir go/ js/
    RUN /protoc/bin/protoc \
        -I=/proto \
        --go_out=go/ \
        --js_out=import_style=commonjs:js/ \
        *.proto
    # disable eslint on generated JS files (https://github.com/grpc/grpc-web/issues/447)
    RUN find js/ -type f -exec sh -c "echo '/* eslint-disable */' | cat - {} > /tmp/out && mv /tmp/out {}" \;
    SAVE ARTIFACT go/pkg/api/ /go AS LOCAL pkg/api
    SAVE ARTIFACT js/ /js AS LOCAL ui/src/api

builder:
    FROM golang
    WORKDIR /go/src/github.com/cartermckinnon/filter-feed
    COPY . .
    COPY +proto/go pkg/api
    RUN go build -o /go/bin/filter-feed
    SAVE ARTIFACT /go/bin/filter-feed AS LOCAL bin/filter-feed

filter-feed:
    FROM ubuntu:21.04
    RUN apt-get update && apt-get install -y ca-certificates
    LABEL org.opencontainers.image.source="https://github.com/cartermckinnon/filter-feed/"
    COPY +builder/filter-feed /usr/bin/filter-feed
    ENTRYPOINT ["/usr/bin/filter-feed"]
    CMD ["server"]
    ARG VERSION="0.0.0-dev"
    SAVE IMAGE --push $IMAGE_REPO/filter-feed:$VERSION

ui-builder:
    FROM node:lts
    WORKDIR /workdir
    COPY ui/package.json .
    COPY ui/package-lock.json .
    COPY ui/webpack.config.js .
    RUN npm install
    COPY ui/src src/
    COPY +proto/js src/api
    RUN npm run build && \
        mkdir -p build/css/ && \
        cp src/css/* build/css/
    SAVE ARTIFACT /workdir/build /ui
    SAVE IMAGE --push $IMAGE_REPO/filter-feed/ui-builder:earthly-cache

ui:
    FROM nginx:stable
    LABEL org.opencontainers.image.source="https://github.com/cartermckinnon/filter-feed"
    COPY +ui-builder/ui /var/www
    COPY ui/nginx.conf /etc/nginx/conf.d/default.conf
    CMD ["nginx","-g","daemon off;"]
    ARG VERSION="latest"
    SAVE IMAGE --push $IMAGE_REPO/filter-feed/ui:$VERSION