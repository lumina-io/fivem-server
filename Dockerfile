FROM alpine:latest AS download_assets
ARG FIVEM_ARTIFACT_URL="https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/"

WORKDIR /fivem

# Download fxserver
RUN apk add --no-cache curl xq \
    && _ARTIFACT=$(curl ${FIVEM_ARTIFACT_URL} | xq -q 'body > section > div > nav > a:nth-child(4)' -a 'href') \
    && curl -#OL "${FIVEM_ARTIFACT_URL}/${_ARTIFACT}" \
    && tar xvf ./fx.tar.xz \
    && rm -rf fx.tar.xz

# Prepare image
FROM alpine:latest
ARG USER_ID='0'
ARG GROUP_ID='0'

# useradd app --uid ${USER_ID} -U -s /bin/bash
RUN apk add --no-cache bash libstdc++ libgcc make lua5.4 \
    && ln -sf /usr/bin/lua5.4 /usr/bin/lua \
    && if [ ${USER_ID} != "0" ]; then \
    adduser app -u ${USER_ID} -h /app -s /bin/bash -D; \
    fi

USER app
WORKDIR /app/fivem

COPY --from=download_assets /fivem/ /app/fivem/
COPY ./template/fivem-server/start.sh /app/fivem/start.sh

CMD ["bash", "/app/fivem/start.sh"]
