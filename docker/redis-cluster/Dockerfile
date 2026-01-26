FROM redis:7.2.5

RUN apt-get update -qq && \
    apt-get install --no-install-recommends -yqq \
        supervisor gettext-base && \
    apt-get clean -yqq

RUN mkdir /redis-conf && mkdir /redis-data

COPY redis-cluster.tmpl /redis-conf/redis-cluster.tmpl
COPY redis.tmpl         /redis-conf/redis.tmpl

# Add startup script
COPY docker-entrypoint.sh /docker-entrypoint.sh

# Add script that generates supervisor conf file based on environment variables
COPY generate-supervisor-conf.sh /generate-supervisor-conf.sh

RUN chmod 755 /docker-entrypoint.sh

EXPOSE 7000 7001 7002 7003 7004 7005 7006 7007

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=20 \
    CMD redis-cli -h 127.0.0.1 -p 7003 CLUSTER INFO | grep 'cluster_state:ok'

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["redis-cluster"]