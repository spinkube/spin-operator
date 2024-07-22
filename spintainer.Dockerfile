FROM debian:bullseye-slim AS install
ARG spin_version=v2.6.0
WORKDIR /opt/install

# Install curl
RUN apt-get update && apt-get install -y curl

# Install git
RUN apt-get install -y git

# Install Static Spin binary
RUN ARCH=$(uname -m | sed s/x86_64/amd64/) && \
  curl -fsSL -o spin-${spin_version}-static-linux-${ARCH}.tar.gz https://github.com/fermyon/spin/releases/download/${spin_version}/spin-${spin_version}-static-linux-${ARCH}.tar.gz && \
  tar xvf spin-${spin_version}-static-linux-${ARCH}.tar.gz

FROM gcr.io/distroless/static-debian11
COPY --from=install /opt/install/spin /usr/local/bin/spin

ENTRYPOINT [ "/usr/local/bin/spin" ]
