FROM debian:bookworm

RUN apt-get update && \
    apt-get install -yqq \
    wget curl xz-utils

WORKDIR /app/fivem
RUN wget https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/7654-8b93ef263f644196fd83621df65c3c0b687da124/fx.tar.xz && \
    tar xvf ./fx.tar.xz

COPY ./template/fivem-server/start.sh /app/fivem/start.sh

CMD ["bash", "/app/fivem/start.sh"]
