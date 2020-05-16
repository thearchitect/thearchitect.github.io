# cat Dockerfile | docker build --build-arg USERNAME=neo -t home - && docker run --rm -it -h zion --network none home
# cat Dockerfile | docker build --build-arg USERNAME=neo -t home - && docker run --rm -it -h zion --network none home

################################################################
#### Base configuration
####
FROM ubuntu:latest as base

ARG DEBIAN_FRONTEND=noninteractive

RUN \
    apt-get update &&\
    apt-get install -y sudo locales software-properties-common \
        zile mc zsh git htop wget curl \
        cmatrix toilet \
        iproute2 nmap iputils-tracepath dnsutils traceroute inetutils-ping inetutils-traceroute

# LOCALE
RUN printf 'en_US.UTF-8 UTF-8' > /etc/locale.gen && locale-gen
ENV LANG        en_US.UTF-8
ENV LANGUAGE    en_US:en
ENV LC_ALL      en_US.UTF-8

# TOILET
RUN \
    printf "toilet --gay --width 170 --font mono12 \"Hello, \${LOGNAME}!\"" > /bin/hello &&\
    chmod +x /bin/hello

RUN \
    wget --output-document=/bin/powerline-go https://github.com/justjanne/powerline-go/releases/latest/download/powerline-go-linux-amd64 &&\
    chmod +x /bin/powerline-go

# CLEANUP
RUN \
    apt-get autoremove -y &&\
    apt-get clean -y &&\
    apt-get purge -y  --allow-remove-essential apt &&\
    rm -rf /var/lib/apt/lists/

# FIXME
# ENV HOSTNAME="zion"
# RUN printf "$HOSTNAME" > /etc/hostname
# RUN printf "127.0.0.1\tlocalhost\t$HOSTNAME" > /etc/hosts

# COPY /thearchitect /thearchitect

################################################################
#### User setup
####
FROM base as user

ARG USERNAME=johnny
ARG USER_UID=1000
ARG USER_GID=$USER_UID
RUN \
    groupadd --gid $USER_GID $USERNAME &&\
    useradd -s $(which zsh) --uid $USER_UID --gid $USER_GID --groups sudo -m $USERNAME &&\
    printf "\n\n$USERNAME ALL=(ALL) NOPASSWD:ALL\n\n" >> /etc/sudoers

# ZSH CONFIG
COPY zshrc /root/.zshrc
COPY zshrc /home/${USERNAME}/.zshrc
RUN \
    chsh -s $(which zsh) root &&\
    chsh -s $(which zsh) ${USERNAME}

RUN chown -R ${USER_UID}:${USER_GID} /home/${USERNAME}

USER $USERNAME

################################################################
#### Final image
####
FROM user

ENTRYPOINT ["zsh", "-l", "-c", "mc / /home/johnny; zsh -l"]
