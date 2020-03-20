FROM alpine/git:v2.24.1 AS branches

WORKDIR /tmp/extract

RUN apk update && apk add --no-cache diffutils
RUN git clone https://github.com/docker/docker-ce.git
RUN git clone https://github.com/silvin-lubecki/cli-extract.git
RUN git clone https://github.com/silvin-lubecki/engine-extract.git
RUN git clone https://github.com/silvin-lubecki/packaging-extract.git
COPY check/check.sh check.sh
ENTRYPOINT [ "/tmp/extract/check.sh" ]
CMD []

FROM branches as tags
COPY check/check-tags.sh check-tags.sh
ENTRYPOINT [ "/tmp/extract/check-tags.sh" ]
CMD []

FROM golang AS detect
COPY --from=branches /tmp/extract /tmp/extract
RUN cd /tmp/extract/docker-ce &&\
    git remote add silvin https://github.com/silvin-lubecki/docker-ce &&\
    git fetch silvin
RUN cd /tmp/extract/packaging-extract &&\
    git remote add docker https://github.com/docker/docker-ce-packaging &&\
    git fetch docker &&\
    git checkout -b 17.06-extract-packaging origin/17.06-extract-packaging &&\
    git checkout -b 17.07-extract-packaging origin/17.07-extract-packaging &&\
    git checkout -b 17.09-extract-packaging origin/17.09-extract-packaging &&\
    git checkout -b 17.10-extract-packaging origin/17.10-extract-packaging &&\
    git checkout -b 17.11-extract-packaging origin/17.11-extract-packaging &&\
    git checkout -b 17.12-extract-packaging origin/17.12-extract-packaging &&\
    git checkout -b 18.01-extract-packaging origin/18.01-extract-packaging &&\
    git checkout -b 18.02-extract-packaging origin/18.02-extract-packaging &&\
    git checkout -b 18.03-extract-packaging origin/18.03-extract-packaging &&\
    git checkout -b 18.04-extract-packaging origin/18.04-extract-packaging &&\
    git checkout -b 18.05-extract-packaging origin/18.05-extract-packaging
RUN cd /tmp/extract/cli-extract &&\
    git remote add docker https://github.com/docker/cli &&\
    git fetch docker &&\
    git checkout -b 17.06-extract-cli origin/17.06-extract-cli &&\
    git checkout -b 17.07-extract-cli origin/17.07-extract-cli &&\
    git checkout -b 17.09-extract-cli origin/17.09-extract-cli &&\
    git checkout -b 17.10-extract-cli origin/17.10-extract-cli &&\
    git checkout -b 17.11-extract-cli origin/17.11-extract-cli &&\
    git checkout -b 17.12-extract-cli origin/17.12-extract-cli &&\
    git checkout -b 18.01-extract-cli origin/18.01-extract-cli &&\
    git checkout -b 18.02-extract-cli origin/18.02-extract-cli &&\
    git checkout -b 18.03-extract-cli origin/18.03-extract-cli &&\
    git checkout -b 18.04-extract-cli origin/18.04-extract-cli &&\
    git checkout -b 18.05-extract-cli origin/18.05-extract-cli
RUN cd /tmp/extract/engine-extract &&\
    git remote add docker https://github.com/docker/engine &&\
    git fetch docker &&\
    git checkout -b 17.06-extract-engine origin/17.06-extract-engine &&\
    git checkout -b 17.07-extract-engine origin/17.07-extract-engine &&\
    git checkout -b 17.09-extract-engine origin/17.09-extract-engine &&\
    git checkout -b 17.10-extract-engine origin/17.10-extract-engine &&\
    git checkout -b 17.11-extract-engine origin/17.11-extract-engine &&\
    git checkout -b 17.12-extract-engine origin/17.12-extract-engine &&\
    git checkout -b 18.01-extract-engine origin/18.01-extract-engine &&\
    git checkout -b 18.02-extract-engine origin/18.02-extract-engine &&\
    git checkout -b 18.03-extract-engine origin/18.03-extract-engine &&\
    git checkout -b 18.04-extract-engine origin/18.04-extract-engine &&\
    git checkout -b 18.05-extract-engine origin/18.05-extract-engine
WORKDIR /go/src/github.com/silvin-lubecki/docker-ce-tags
COPY . .

RUN go get gopkg.in/src-d/go-git.v4
RUN go get gopkg.in/yaml.v2
RUN make build
ENTRYPOINT [ "./docker-ce-tags" ]
CMD []
