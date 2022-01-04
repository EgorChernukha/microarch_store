FROM debian:9-slim

RUN apt-get update && \
    apt-get install -y --no-install-suggests --no-install-recommends ca-certificates && \
    apt-get clean && \
    groupadd -g 1001 microuser && \
    useradd -u 1001 -r -g 1001 -s /sbin/nologin -c "go microservice user" microuser

ADD ./bin/order /app/bin/
WORKDIR /app

ADD ./data/mysql/migrations/order /data/mysql/migrations/order
ENV ORDER_MIGRATIONS_DIR=/data/mysql/migrations/order

EXPOSE 8000

USER microuser
CMD [ "/app/bin/order" ]
