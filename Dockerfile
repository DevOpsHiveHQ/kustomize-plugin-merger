FROM scratch
ENTRYPOINT ["/kustomize-plugin-merger"]
COPY kustomize-plugin-merger /
