FROM debian:bullseye-slim

# Install curl
RUN apt-get update && apt-get install -y curl

# Install git
RUN apt-get install -y git

# Install Spin
RUN curl -fsSL https://developer.fermyon.com/downloads/install.sh | bash
RUN mv spin /usr/local/bin

ENTRYPOINT [ "spin" ]
