FROM debian:bookworm

RUN apt-get update && \
    apt-get install -yqq \
    wget curl xz-utils

WORKDIR /app/fivem
RUN wget https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/7878-37639cddcb2e5b0ef05f96d5a482b236b7349a1e/fx.tar.xz && \
    tar xvf ./fx.tar.xz

COPY ./template/fivem-server/start.sh /app/fivem/start.sh

CMD ["bash", "/app/fivem/start.sh"]
