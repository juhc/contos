FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TERM=xterm \
    SYSLINUX_SITE=https://mirrors.edge.kernel.org/ubuntu/pool/main/s/syslinux \
    SYSLINUX_VERSION=4.05+dfsg-6+deb8u1

RUN apt-get -q update && \
    apt-get -q -y install --no-install-recommends \
      ca-certificates \
      bc build-essential cpio file git python3 unzip rsync wget curl \
      syslinux syslinux-common isolinux xorriso dosfstools mtools \
      xz-utils patch && \
    wget -q "${SYSLINUX_SITE}/syslinux-common_${SYSLINUX_VERSION}_all.deb" && \
    wget -q "${SYSLINUX_SITE}/syslinux_${SYSLINUX_VERSION}_amd64.deb" && \
    dpkg -i "syslinux-common_${SYSLINUX_VERSION}_all.deb" && \
    dpkg -i "syslinux_${SYSLINUX_VERSION}_amd64.deb" && \
    rm -f "syslinux-common_${SYSLINUX_VERSION}_all.deb" "syslinux_${SYSLINUX_VERSION}_amd64.deb" && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

ENV SRC_DIR=/build \
    OVERLAY=/overlay \
    BR_ROOT=/build/buildroot

RUN mkdir -p ${SRC_DIR} ${OVERLAY}

ENV BR_VERSION 2024.11
RUN wget -qO- https://buildroot.org/downloads/buildroot-${BR_VERSION}.tar.xz | tar xJ && \
    mv buildroot-${BR_VERSION} ${BR_ROOT}

COPY overlay ${OVERLAY}
WORKDIR ${OVERLAY}

RUN mkdir -p etc/ssl/certs && \
    cp /etc/ssl/certs/ca-certificates.crt etc/ssl/certs/

RUN mkdir -p usr/share/bash-completion/completions && \
    wget -qO usr/share/bash-completion/bash_completion https://raw.githubusercontent.com/scop/bash-completion/master/bash_completion && \
    chmod +x usr/share/bash-completion/bash_completion

ENV DOCKER_VERSION 26.1.0
ENV DOCKER_REVISION barge.2
RUN mkdir -p usr/bin && \
    wget -qO- https://download.docker.com/linux/static/stable/x86_64/docker-${DOCKER_VERSION}.tgz | tar xz -C usr/bin --strip-components=1 docker/docker

RUN wget -qO usr/share/bash-completion/completions/docker https://raw.githubusercontent.com/docker/cli/v${DOCKER_VERSION}/contrib/completion/bash/docker

ENV DINIT_VERSION 1.2.5
RUN mkdir -p usr/bin && \
    wget -qO usr/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v${DINIT_VERSION}/dumb-init_${DINIT_VERSION}_x86_64 && \
    chmod +x usr/bin/dumb-init

ENV VERSION 1.0
RUN mkdir -p etc && \
    echo -e "█▀▀ █▀█ █▄░█ ▀█▀ █▀█ █▀\n█▄▄ █▄█ █░▀█ ░█░ █▄█ ▄█\nWelcome to Contos ${VERSION}, Docker version ${DOCKER_VERSION}" > etc/motd && \
    echo "NAME=\"Contos\"" > etc/os-release && \
    echo "VERSION=${VERSION}" >> etc/os-release && \
    echo "ID=contos" >> etc/os-release && \
    echo "ID_LIKE=busybox" >> etc/os-release && \
    echo "VERSION_ID=${VERSION}" >> etc/os-release && \
    echo "PRETTY_NAME=\"Contos ${VERSION}\"" >> etc/os-release

RUN wget -qO usr/bin/pkg https://raw.githubusercontent.com/bargees/barge-pkg/master/pkg && \
    chmod +x usr/bin/pkg

COPY configs ${SRC_DIR}/configs
RUN cp ${SRC_DIR}/configs/buildroot.config ${BR_ROOT}/.config && \
    cp ${SRC_DIR}/configs/busybox.config ${BR_ROOT}/package/busybox/busybox.config

COPY scripts ${SRC_DIR}/scripts

VOLUME ${BR_ROOT}/dl ${BR_ROOT}/ccache

WORKDIR ${BR_ROOT}

ENV FORCE_UNSAFE_CONFIGURE=1

CMD ["../scripts/build.sh"]