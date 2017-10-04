FROM buildpack-deps:xenial

RUN apt-get update \
    && apt-get -y upgrade \
    && apt-get -y install python-pip make build-essential curl openssl vim jq gettext \
    && rm -rf /var/lib/apt/lists/*

ENV GO_VERSION 1.8
RUN mkdir /tmp/godeb \
    && curl "https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz" > /tmp/godeb/godeb.tar.gz \
    && (cd /tmp/godeb; tar zvxf godeb.tar.gz; ./godeb install "${GO_VERSION}") \
    && rm -rf /tmp/godeb

# See: https://github.com/Azure/azure-cli/blob/master/packaged_releases/bundled/README.md#using-the-bundled-installer
ENV AZURE_CLI_BUNDLE_VERSION 0.2.10-1
RUN mkdir /tmp/azurecli \
    && curl "https://azurecliprod.blob.core.windows.net/bundled/azure-cli_bundle_${AZURE_CLI_BUNDLE_VERSION}.tar.gz" > /tmp/azurecli/azure-cli_bundle.tar.gz \
    && (cd /tmp/azurecli \
      && tar -xvzf azure-cli_bundle.tar.gz \
      && azure-cli_bundle_*/installer --bin-dir /usr/local/bin) \
    && rm -rf /tmp/azurecli

RUN mkdir /gopath \
    && mkdir /gopath/src \
    && mkdir /gopath/src/acs-engine-template-generator \
    && mkdir /gopath/bin
ENV GOPATH /gopath
ENV GOROOT /usr/local/go
ENV PATH "${PATH}:${GOPATH}/bin:$GOROOT/bin"

WORKDIR /tmp
RUN wget https://github.com/Azure/acs-engine/releases/download/v0.6.0/acs-engine-v0.6.0-linux-amd64.tar.gz \
    && tar -xf acs-engine-v0.6.0-linux-amd64.tar.gz \
    && mv /tmp/acs-engine-v0.6.0-linux-amd64/acs-engine ${GOPATH}/bin/acs-engine

RUN go get -d github.com/Azure/acs-engine
COPY src/acs-engine-template-generator.go ${GOPATH}/src/acs-engine-template-generator
RUN go install acs-engine-template-generator

COPY src/deploy.sh ${GOPATH}/bin
RUN sed -i 's/\r$//' ${GOPATH}/bin/deploy.sh
WORKDIR ${GOPATH}/bin
ENTRYPOINT ${GOPATH}/bin/deploy.sh
