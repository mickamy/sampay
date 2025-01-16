ARG GO_VERSION=1.23.4

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

ARG GH_PAT
RUN git config --global url."https://${GH_PAT}@github.com/".insteadOf "https://github.com/"

ADD Makefile ./
RUN make install-build

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH
ARG MAIN_GO_FILE

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /build/app $MAIN_GO_FILE

FROM alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=build /build/app /build/

EXPOSE 8080

ENTRYPOINT [ "/build/app" ]
