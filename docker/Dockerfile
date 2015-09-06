FROM golang:1.5

RUN git clone git://git.postgresql.org/git/postgresql $HOME/.pgbrew/repository

RUN apt-get update && \
    apt-get install -y \
        libreadline6 \
        libreadline6-dev \
        zlib1g-dev \
        bison \
        flex

RUN cd $HOME/.pgbrew/repository && \
    git checkout REL9_4_4 && \
    ./configure --prefix=$HOME/.pgbrew/versions/9.4.4 && \
    make install && \
    echo -n "7c055f3ec3bd338a1ebb8c73cff3d01df626471e\tREL9_4_4" > $HOME/.pgbrew/versions/9.4.4/.pgbrew_info && \
    make clean

RUN cd $HOME/.pgbrew/repository && \
    git checkout REL9_3_9 && \
    ./configure --prefix=$HOME/.pgbrew/versions/9.3.9-debug --enable-debug --enable-cassert && \
    make install && \
    echo -n "553e576e05b50f9faffbd3dd721e44fc3746898d\tREL9_3_9" > $HOME/.pgbrew/versions/9.3.9-debug/.pgbrew_info && \
    make clean

RUN echo -n "{\"base-path\":\"$HOME/.pgbrew\"}" > $HOME/.pgbrew/config.json && \
    ln -s $HOME/.pgbrew/versions/9.4.4 $HOME/.pgbrew/current
