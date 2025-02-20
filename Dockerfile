# Stage 1: Download gittisane
FROM alpine:latest as downloader
RUN apk add --no-cache curl bash grep
WORKDIR /download/
COPY download.sh .
RUN /download/download.sh
RUN ls /download/downloads
RUN cp /download/downloads/gitea-Linux /download/gittisane-linux-amd64

# Stage 2: Runtime
FROM geti2p/i2p:latest

# Copy gittisane binary
COPY --from=downloader /download/gittisane-linux-amd64 /usr/local/bin/gittisane
RUN chmod +x /usr/local/bin/gittisane

# Create data directories
RUN mkdir -p /data/gitea

# Configure I2P for SAM
RUN echo "i2cp.tcp.host=127.0.0.1\n\
i2cp.tcp.port=7654\n\
sam.enabled=true\n\
sam.host=127.0.0.1\n\
sam.port=7656" >> /i2p/router.config

# Setup volumes
VOLUME ["/data/gitea", "/i2p/.i2p"]
WORKDIR /data/gitea

# Create startup script
COPY start.sh /usr/local/bin/start.sh

ENTRYPOINT ["start.sh"]