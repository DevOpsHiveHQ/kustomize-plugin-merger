FROM ubuntu:latest@sha256:f144425ff09be612d6d9ad965196e9cdc23dae1f42110a8a11a3e9a8198759f7 as base
RUN useradd -u 1001 merger

FROM scratch
COPY --from=base /etc/passwd /etc/passwd
COPY kustomize-plugin-merger /
USER 1001
ENTRYPOINT ["/kustomize-plugin-merger"]
