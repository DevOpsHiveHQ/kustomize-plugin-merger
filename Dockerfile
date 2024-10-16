FROM ubuntu:latest@sha256:d4f6f70979d0758d7a6f81e34a61195677f4f4fa576eaf808b79f17499fd93d1 as base
RUN useradd -u 1001 merger

FROM scratch
COPY --from=base /etc/passwd /etc/passwd
COPY kustomize-plugin-merger /
USER 1001
ENTRYPOINT ["/kustomize-plugin-merger"]
