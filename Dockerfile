# syntax=docker/dockerfile:1.7-labs

# (i.e. 'cgr.dev', 'docker.io')
ARG REGISTRY=docker.io

FROM --platform=${BUILDPLATFORM} ${REGISTRY}/chainguard/go AS go
WORKDIR /src

ARG TARGETOS TARGETARCH
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

COPY go.* .
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go mod download

FROM go AS build
# only really needed for CGO deps (like SQLite, etc.)
ENV CGO_ENABLED=0

COPY . .
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build \
    -tags=osusergo,netgo,timetzdata \
    -ldflags="-s -w -extldflags=-static" \
    -o=/usr/local/bin/a.out ./cmd/quickstart/

FROM ${REGISTRY}/chainguard/cc-dynamic

COPY --from=build /usr/local/bin /usr/local/bin

ENTRYPOINT [ "a.out" ]
CMD [ "serve" ]