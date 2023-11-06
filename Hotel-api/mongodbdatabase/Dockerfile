FROM mongo:bionic

COPY hotel-data.json /hotel-data.json
COPY amenity-data.json /amenity-data.json

CMD mongoimport --host mongodatabase --username root --password pass --authenticationDatabase admin --db test --collection hotels --type json --file /hotel-data.json --jsonArray --drop -v && \
    mongoimport --host mongodatabase --username root --password pass --authenticationDatabase admin --db test --collection amenities --type json --file /amenity-data.json --jsonArray --drop -v
