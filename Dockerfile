# From https://gist.github.com/AverageMarcus/78fbcf45e72e09d9d5e75924f0db4573
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21 as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app/
ADD . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o kustomize-plugin-merger main.go

FROM --platform=${TARGETPLATFORM:-linux/amd64} scratch
WORKDIR /app/
COPY --from=builder /app/kustomize-plugin-merger /app/kustomize-plugin-merger
ENTRYPOINT ["/app/kustomize-plugin-merger"]