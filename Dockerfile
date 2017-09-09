FROM gocd/gocd-server:v17.9.0

ARG UID

RUN apk --no-cache add shadow && \
    usermod -u ${UID} go