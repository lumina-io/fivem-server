FROM alpine:latest AS download_assets
ARG FIVEM_ARTIFACT_URL="https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/"

WORKDIR /fivem

# Download fxserver
RUN apk add --no-cache curl xq \
    && _ARTIFACT=$(curl ${FIVEM_ARTIFACT_URL} | xq -q 'body > section > div > nav > a:nth-child(4)' -a 'href') \
    && curl -#OL "${FIVEM_ARTIFACT_URL}/${_ARTIFACT}" \
    && tar xvf ./fx.tar.xz \
    && rm -rf fx.tar.xz

# For debug: direct download
#ARG DIRECT_URL="https://runtime.fivem.net/artifacts/fivem/build_proot_linux/master/16237-081a76900402908b96facbdbc3d2299ebf1c4714/fx.tar.xz"
#RUN apk add --no-cache curl \
#    && curl -#OL "${DIRECT_URL}" \
#    && tar xvf ./fx.tar.xz \
#    && rm -rf fx.tar.xz

FROM golang:latest AS build_utils
WORKDIR /app
COPY ./kontra/go.mod ./kontra/go.sum /app/
RUN go mod download

COPY ./kontra .
RUN go build -a -v -o /kontra ./main.go

# Prepare image
FROM alpine:latest
ARG USER_ID='0'
ARG GROUP_ID='0'

# useradd app --uid ${USER_ID} -U -s /bin/bash
RUN apk add --no-cache bash libstdc++ libgcc make lua5.4 tzdata \
    && ln -sf /usr/bin/lua5.4 /usr/bin/lua \
    && if [ ${USER_ID} != "0" ]; then \
    adduser app -u ${USER_ID} -h /app -s /bin/bash -D; \
    fi

USER app
WORKDIR /app/fivem

COPY --from=build_utils /kontra /bin/kontra
COPY --from=download_assets /fivem/ /app/fivem/
COPY ./template/fivem-server/start.sh /app/fivem/start.sh

CMD ["bash", "/app/fivem/start.sh"]
