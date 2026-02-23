FROM postgres:18-alpine

ARG SQL_TAP_VERSION=0.2.3
ARG TARGETARCH

ADD https://github.com/mickamy/sql-tap/releases/download/v${SQL_TAP_VERSION}/sql-tap_${SQL_TAP_VERSION}_linux_${TARGETARCH}.tar.gz /tmp/sql-tap.tar.gz
RUN tar -xzf /tmp/sql-tap.tar.gz -C /usr/local/bin sql-tapd && rm /tmp/sql-tap.tar.gz

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]
EXPOSE 5432 5433 9091 8081
