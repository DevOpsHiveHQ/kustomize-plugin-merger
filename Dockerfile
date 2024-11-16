FROM ubuntu:latest@sha256:278628f08d4979fb9af9ead44277dbc9c92c2465922310916ad0c46ec9999295 as base
RUN useradd -u 1001 merger

FROM scratch
COPY --from=base /etc/passwd /etc/passwd
COPY kustomize-plugin-merger /
USER 1001
ENTRYPOINT ["/kustomize-plugin-merger"]
