# syntax=docker/dockerfile:1.7

FROM alpine:latest AS download_assets
ARG FIVEM_ARTIFACT_URL

WORKDIR /fivem

# Download fxserver
RUN FIVEM_ARTIFACT_BASE="https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/" \
    && if [ "${FIVEM_ARTIFACT_URL}" = "" ]; then \
        apk add --no-cache xq; \
        _ARTIFACT=$(wget -O- ${FIVEM_ARTIFACT_BASE} | xq -q 'body > section > div > nav > a:nth-child(4)' -a 'href'); \
        FIVEM_ARTIFACT_URL="${FIVEM_ARTIFACT_BASE}/${_ARTIFACT}"; \
    fi \
    && wget -O- "${FIVEM_ARTIFACT_URL}" \
    | tar xvJ \
        --exclude alpine/bin \
        --exclude alpine/dev \
        --exclude alpine/etc \
        --exclude alpine/home \
        # --exclude alpine/lib \
        --exclude alpine/lib64 \
        --exclude alpine/media \
        --exclude alpine/mnt \
        --exclude alpine/proc \
        --exclude alpine/root \
        --exclude alpine/run \
        --exclude alpine/sbin \
        --exclude alpine/srv \
        --exclude alpine/sys \
        --exclude alpine/tmp \
        # --exclude alpine/usr \
        --exclude alpine/usr/bin \
        # --exclude alpine/usr/glibc-compat \
        --exclude alpine/usr/sbin \
        --exclude alpine/usr/share \
        --exclude alpine/usr/libexec \
        --exclude alpine/utils \
        --exclude alpine/var

FROM golang:latest AS build_utils
WORKDIR /app
COPY --link ./kontra/go.mod ./kontra/go.sum /app/
RUN go mod download

COPY --link ./kontra .
RUN go build -a -v -o /kontra ./cmd/kontra

# Prepare image
FROM alpine:latest
ARG USER_ID='0'
ARG GROUP_ID='0'

ENV TXHOST_TXA_URL="http://localhost:40120"
ENV TXHOST_TXA_PORT="40120"
ENV TXHOST_TMP_HIDE_ADS="true"

# libstdc++ -> for fxserver
# gcompat -> glibc-compat (fxserver compat / kontra compat)
RUN apk add --no-cache lua5.4 tzdata libstdc++ gcompat \
    && ln -sf /usr/bin/lua5.4 /usr/bin/lua \
    && if [ "${USER_ID}" != "0" ]; then \
    adduser app -u ${USER_ID} -h /app -s /bin/sh -D; \
    fi

COPY --link --from=build_utils /kontra /bin/kontra
COPY --link ./template/fivem-server/start.sh /app/fivem/start.sh

USER app
WORKDIR /app/fivem

COPY --link --from=download_assets /fivem/ /app/fivem/

CMD ["kontra", "sh", "/app/fivem/start.sh"]
