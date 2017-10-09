FROM resin/rpi-raspbian
MAINTAINER steve@bargelt.com

# add our user and group first to make sure their IDs get assigned consistently, regardless of whatever dependencies get added
RUN groupadd -r mongodb && useradd -r -g mongodb mongodb

RUN apt-get update \
	&& apt-get install -y --no-install-recommends \
		ca-certificates \
		jq \
	&& rm -rf /var/lib/apt/lists/*

RUN mkdir /docker-entrypoint-initdb.d

# Install required packages
RUN apt-get update && apt-get -y dist-upgrade && \
    apt-get install -y --no-install-suggests mongodb && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /data/db /data/configdb \
	&& chown -R mongodb:mongodb /data/db /data/configdb
VOLUME /data/db /data/configdb

COPY docker-entrypoint.sh /usr/local/bin/
RUN ln -s usr/local/bin/docker-entrypoint.sh /entrypoint.sh # backwards compat
ENTRYPOINT ["docker-entrypoint.sh"]

EXPOSE 27017
CMD ["mongod"]