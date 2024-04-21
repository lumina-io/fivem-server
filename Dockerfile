FROM debian:bookworm

RUN apt-get update && \
    apt-get install -yqq \
    wget curl xz-utils

WORKDIR /app/fivem
RUN wget https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/7969-55dab4e102a780a94c0f3cfa54fd2e6a0c069f89/fx.tar.xz && \
    tar xvf ./fx.tar.xz

COPY ./template/fivem-server/start.sh /app/fivem/start.sh

CMD ["bash", "/app/fivem/start.sh"]
