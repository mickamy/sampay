FROM alpine AS build

ARG GRPC_TAP_VERSION=0.0.3
ARG TARGETARCH

ADD https://github.com/mickamy/grpc-tap/releases/download/v${GRPC_TAP_VERSION}/grpc-tap_${GRPC_TAP_VERSION}_linux_${TARGETARCH}.tar.gz /tmp/grpc-tap.tar.gz
RUN tar -xzf /tmp/grpc-tap.tar.gz -C /usr/local/bin && rm /tmp/grpc-tap.tar.gz

FROM alpine

RUN apk add --no-cache curl

COPY --from=build /usr/local/bin/grpc-tapd /usr/local/bin/grpc-tapd
