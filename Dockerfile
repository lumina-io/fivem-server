FROM debian:bookworm

ARG FIVEM_ARTIFACT_URL="https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/9635-0776d79840adc87a90786a4ad9ad9b12dacb8886/fx.tar.xz"

RUN apt-get update && \
    apt-get install -yqq \
    wget curl xz-utils

WORKDIR /app/fivem
RUN wget $FIVEM_ARTIFACT_URL && \
    tar xvf ./fx.tar.xz

COPY ./template/fivem-server/start.sh /app/fivem/start.sh

CMD ["bash", "/app/fivem/start.sh"]
