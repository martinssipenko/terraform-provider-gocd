FROM gocd/gocd-server:v17.8.0

ARG UID

RUN apk --no-cache add shadow && \
    usermod -u ${UID} go