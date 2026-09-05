FROM ubuntu:latest@sha256:2260313b31c8c011cd2eebe728008efac1b3982be73eb71348ea2648d2c0e09b as base
RUN useradd -u 1001 merger

FROM scratch
COPY --from=base /etc/passwd /etc/passwd
COPY kustomize-plugin-merger /
USER 1001
ENTRYPOINT ["/kustomize-plugin-merger"]
