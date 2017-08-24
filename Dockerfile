ARG GOCD_VERSION=17.8.0
FROM gocd/gocd-server:v${GOCD_VERSION}

ARG UID

COPY ./scripts/shutdown-clean-goserver.sh /shutdown.sh

RUN apk --no-cache add shadow && \
    usermod -u ${UID} go && \
    chmod +x /shutdown.sh