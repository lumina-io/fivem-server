FROM debian:bookworm

RUN apt-get update && \
    apt-get install -yqq \
    wget curl xz-utils

WORKDIR /app/fivem
RUN wget https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/7669-5e71f91d7bd846d8eb059c6a7a0db79f16e9d601/fx.tar.xz && \
    tar xvf ./fx.tar.xz

COPY ./template/fivem-server/start.sh /app/fivem/start.sh

CMD ["bash", "/app/fivem/start.sh"]
