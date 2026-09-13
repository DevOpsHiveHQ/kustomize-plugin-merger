FROM ubuntu:latest@sha256:513c074113a871b51a8d16ab445c88779d6452d937a164fb5cc479f32668a41d as base
RUN useradd -u 1001 merger

FROM scratch
COPY --from=base /etc/passwd /etc/passwd
COPY kustomize-plugin-merger /
USER 1001
ENTRYPOINT ["/kustomize-plugin-merger"]
