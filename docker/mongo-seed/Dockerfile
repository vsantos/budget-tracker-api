FROM mongo:4.4.0

COPY init.json /init.json

CMD mongoimport --host localhost --db reach-engine --collection users --type json --file /init.json