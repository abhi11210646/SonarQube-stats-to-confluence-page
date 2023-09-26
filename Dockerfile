FROM harbor.one.com/standard-images/ubuntu:focal

# Install necessary packages
RUN apt-get update && \
apt-get install -y --no-install-recommends \
python3 build-essential

# Install Nodejs
ENV NODE_VERSION 18.18.0
RUN set -eux && \
  cd / && \
  curl -O https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.gz && \
  tar -zxf node-v$NODE_VERSION-linux-x64.tar.gz -C /usr/lib && \
  ln -s /usr/lib/node-v$NODE_VERSION-linux-x64/bin/node /usr/bin/node && \
  ln -s /usr/lib/node-v$NODE_VERSION-linux-x64/bin/npm /usr/bin/npm && \
  ln -s /usr/lib/node-v$NODE_VERSION-linux-x64/bin/npx /usr/bin/npx && \
  rm node-v$NODE_VERSION-linux-x64.tar.gz


# TODO: Change user from root to node(any name)

RUN mkdir /app

WORKDIR /app

COPY package.json .

# Install packages
RUN npm install

# Copy rest of your code
COPY . .

# CMD ["node","script.js"]
