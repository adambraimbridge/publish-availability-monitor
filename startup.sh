sed -i "s QUEUE_ADDR $QUEUE_ADDR " /config.json
sed -i "s URL1 $URL1 " /config.json
sed -i "s URL2 $URL2 " /config.json
sed -i "s URL3 $URL3 " /config.json
sed -i "s URL4 $URL4 " /config.json

./app -config config.json
