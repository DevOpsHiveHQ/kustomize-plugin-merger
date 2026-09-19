FROM ubuntu:latest@sha256:da6fc2be547864451aa253836dd926da33623312df4a9a243e35dc877c378a78 as base
RUN useradd -u 1001 merger

FROM scratch
COPY --from=base /etc/passwd /etc/passwd
COPY kustomize-plugin-merger /
USER 1001
ENTRYPOINT ["/kustomize-plugin-merger"]
