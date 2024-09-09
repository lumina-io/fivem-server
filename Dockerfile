FROM debian:bookworm

ARG FIVEM_ARTIFACT_URL="https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/9780-861f70ca1991dba2b59b544ec776e8ffbbe0263f/fx.tar.xz"

RUN apt-get update && \
    apt-get install -yqq \
    wget curl xz-utils

WORKDIR /app/fivem
RUN wget $FIVEM_ARTIFACT_URL && \
    tar xvf ./fx.tar.xz

COPY ./template/fivem-server/start.sh /app/fivem/start.sh

CMD ["bash", "/app/fivem/start.sh"]
