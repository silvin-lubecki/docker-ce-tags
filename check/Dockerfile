FROM alpine/git:v2.24.1 AS branches

WORKDIR /tmp/extract

RUN apk update && apk add --no-cache diffutils
RUN git clone https://github.com/docker/docker-ce.git
RUN git clone https://github.com/silvin-lubecki/cli-extract.git
RUN git clone https://github.com/silvin-lubecki/engine-extract.git
RUN git clone https://github.com/silvin-lubecki/packaging-extract.git
COPY check.sh check.sh
ENTRYPOINT [ "/tmp/extract/check.sh" ]
CMD []

FROM branches as tags
COPY check-tags.sh check-tags.sh
ENTRYPOINT [ "/tmp/extract/check-tags.sh" ]
CMD []