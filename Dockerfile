FROM debian:bookworm

ARG FIVEM_ARTIFACT_URL="https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/"
ARG USER_ID='0'
ARG GROUP_ID='0'

RUN apt-get update && \
    apt-get install -yqq \
    wget curl xz-utils make git lua5.4

RUN curl -sSL https://bit.ly/install-xq | bash && \
    mv /usr/local/bin/xq /usr/bin/xq

# Ensure user
RUN useradd app --uid ${USER_ID} -U -s /bin/bash
USER app

WORKDIR /app/fivem

RUN _ARTIFACT=$(curl ${FIVEM_ARTIFACT_URL} | xq -q 'body > section > div > nav > a:nth-child(4)' -a 'href') && \
    wget "${FIVEM_ARTIFACT_URL}/${_ARTIFACT}" && \
    tar xvf ./fx.tar.xz

COPY ./template/fivem-server/start.sh /app/fivem/start.sh

CMD ["bash", "/app/fivem/start.sh"]
