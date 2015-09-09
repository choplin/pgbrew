FROM golang:1.5

ADD pgenv /bin/

RUN apt-get update && \
    apt-get install -y \
        libreadline6 \
        libreadline6-dev \
        zlib1g-dev \
        bison \
        flex

RUN groupadd -r postgres && useradd -m -r -g postgres postgres

USER postgres

RUN pgenv init

RUN pgenv install REL9_4_4 && \
    pgenv install -d REL9_3_9

RUN pgenv cluster create 9.4.4 && \
    pgenv cluster create 9.3.9-debug
