sed -i "s QUEUE_ADDR $QUEUE_ADDR " /config.json
sed -i "s ENVIRONMENT $ENVIRONMENT " /config.json

counter=1
for url in $(echo $URLS | tr "," "\n"); do
	sed -i "s URL$counter $url " /config.json
	let counter=$counter+1
done
./app -config config.json
