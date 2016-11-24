#!/bin/sh
sed -i "s KAFKA_PROXY_HOST $KAFKA_PROXY_HOST " /config.json
sed -i "s QUEUE_ADDR $QUEUE_ADDR " /config.json
sed -i "s CONTENT_URL $CONTENT_URL " /config.json
sed -i "s LISTS_URL $LISTS_URL " /config.json
sed -i "s NOTIFICATIONS_URL $NOTIFICATIONS_URL " /config.json
sed -i "s NOTIFICATIONS_PUSH_URL $NOTIFICATIONS_PUSH_URL " /config.json
sed -i "s METHODE_ARTICLE_TRANSFORMER_URL $METHODE_ARTICLE_TRANSFORMER_URL " /config.json

exec ./publish-availability-monitor -etcd-peers $ETCD_PEERS -config /config.json
