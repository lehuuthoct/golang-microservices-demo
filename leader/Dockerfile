FROM alpine:latest

EXPOSE 9001

# add leader service to container
ADD leader-service-linux-amd64 /

# add system checking service to container
ADD system-checker-linux-amd64 /

# check health automatically of leader service every 5s
HEALTHCHECK --interval=5s --timeout=5s CMD ["./system-checker-linux-amd64", "-port=9001"] || exit 1

ENTRYPOINT ["./leader-service-linux-amd64"]


