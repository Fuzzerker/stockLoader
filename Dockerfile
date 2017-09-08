FROM scratch
COPY stockLoader /
ENTRYPOINT ["./stockLoader"]
