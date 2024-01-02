FROM --platform=$BUILDPLATFORM golang:1.21 AS builder

WORKDIR /workspace

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /workspace
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o epgstation-exporter


FROM gcr.io/distroless/static:nonroot
LABEL maintainer="Miscord Developers <info@miscord.win>"
WORKDIR /
COPY --from=builder /workspace/epgstation-exporter .
USER 65532:65532

ENTRYPOINT [ "/epgstation-exporter" ]
EXPOSE 2121

