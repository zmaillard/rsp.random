# syntax=docker/dockerfile:1

##
## Build
##
ARG XCPUTRANSLATE_VERSION=v0.8.0
ARG BUILDPLATFORM=linux/amd64
FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPUTRANSLATE_VERSION} AS xcputranslate

FROM --platform=${BUILDPLATFORM} golang:1.26-bookworm AS build

WORKDIR /app
COPY --from=xcputranslate /xcputranslate xcputranslate
ARG TARGETPLATFORM

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

ENV CGO_ENABLED=0
ARG VERSION_APP=0.0.0
ENV VERSION ${VERSION_APP}

RUN GOARCH="$(./xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(./xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
    go build \
    --ldflags "-X 'main.Version=${VERSION_APP}'" \
    -o /main main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
ARG TARGETPLATFORM

WORKDIR /
RUN mkdir /data

COPY --from=build /main /main
EXPOSE 1333

USER nonroot:nonroot

ENTRYPOINT ["/main"]
