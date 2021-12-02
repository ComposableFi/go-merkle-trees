# Update the VARIANT arg in devcontainer.json to pick an Go version
ARG VARIANT=latest
FROM golang:${VARIANT}

# This Dockerfile adds a non-root user with sudo access. Update the “remoteUser” property in
# devcontainer.json to use it. More info: https://aka.ms/vscode-remote/containers/non-root-user.
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Options for common setup script - SHA generated on release
ARG INSTALL_ZSH="true"
ARG UPGRADE_PACKAGES="false"
ARG COMMON_SCRIPT_SOURCE="https://raw.githubusercontent.com/microsoft/vscode-dev-containers/master/script-library/common-debian.sh"
ARG COMMON_SCRIPT_SHA="dev-mode"

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Install needed packages and setup non-root user. Use a separate RUN statement to add your own dependencies.
RUN apt-get update \
    && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends python3 python3-pip default-jre curl ca-certificates 2>&1 \
    && curl -sSL  ${COMMON_SCRIPT_SOURCE} -o /tmp/common-setup.sh \
    && ([ "${COMMON_SCRIPT_SHA}" = "dev-mode" ] || (echo "${COMMON_SCRIPT_SHA} /tmp/common-setup.sh" | sha256sum -c -)) \
    && /bin/bash /tmp/common-setup.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" \
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

# Install
RUN apt-get update \
    &&  apt-get -y install --no-install-recommends \
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*


WORKDIR /tmp
RUN curl -s https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-4.6.2.2472-linux.zip -o sonar-scanner-cli-4.6.2.2472-linux.zip

# Install pre-commit
RUN pip install --no-cache-dir pre-commit==2.15.0

# Install required packages fot stack installation
RUN apt-get update -y && apt-get install -y --no-install-recommends libffi-dev libgmp-dev xz-utils zlib1g-dev \
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

# Install stack
RUN curl -sSL https://get.haskellstack.org/ | sh

# Visual Studio Code extensions path
USER vscode
RUN mkdir -p /home/$USERNAME/.vscode-server/extensions \
    /home/$USERNAME/.vscode-server-insiders/extensions \
    && chown -R $USERNAME \
    /home/$USERNAME/.vscode-server \
    /home/$USERNAME/.vscode-server-insiders


# Install Go tools
ARG GO_TOOLS_WITH_MODULES="\
    golang.org/x/tools/gopls \
    honnef.co/go/tools/... \
    golang.org/x/tools/cmd/gorename \
    golang.org/x/tools/cmd/goimports \
    golang.org/x/tools/cmd/guru \
    golang.org/x/lint/golint \
    github.com/mdempsky/gocode \
    github.com/cweill/gotests/... \
    github.com/haya14busa/goplay/cmd/goplay \
    github.com/sqs/goreturns \
    github.com/josharian/impl \
    github.com/davidrjenni/reftools/cmd/fillstruct \
    github.com/uudashr/gopkgs/v2/cmd/gopkgs \
    github.com/ramya-rao-a/go-outline \
    github.com/acroca/go-symbols \
    github.com/godoctor/godoctor \
    github.com/rogpeppe/godef \
    github.com/zmb3/gogetdoc \
    github.com/fatih/gomodifytags \
    github.com/mgechev/revive \
    github.com/go-delve/delve/cmd/dlv \
    github.com/fzipp/gocyclo/cmd/gocyclo \
    github.com/securego/gosec/v2/cmd/gosec"

RUN mkdir -p /tmp/gotools
WORKDIR /tmp/gotools
USER root
RUN export GOPATH=/tmp/gotools \
    # Go tools w/module support
    && export GO111MODULE=on \
    && (echo "${GO_TOOLS_WITH_MODULES}" | xargs -n 1 go get -x )2>&1 \
    # gocode-gomod
    && export GO111MODULE=auto \
    && go get -x -d github.com/stamblerre/gocode 2>&1 \
    && go build -o gocode-gomod github.com/stamblerre/gocode \
    # golangci-lint
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin 2>&1 \
    # Move Go tools into path and clean up
    && mv /tmp/gotools/bin/* /usr/local/bin/ \
    && mv gocode-gomod /usr/local/bin/ \
    && rm -rf /tmp/gotools
USER vscode

ENV GO111MODULE=auto

ENV GOPATH=/home/${USERNAME}/go
ENV PATH="${PATH}:$GOPATH/bin"

# COPY install_dependencies.sh /tmp
# RUN cd /tmp && ./install_dependencies.sh

RUN mkdir -p $GOPATH/src/github.com/ComposableFi/merkle-go
RUN chown -R vscode $GOPATH
