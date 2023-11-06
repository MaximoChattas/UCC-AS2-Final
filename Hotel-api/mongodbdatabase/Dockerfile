FROM mongo:bionic

COPY init-data.json /init-data.json

CMD mongoimport                                 \
    # Host (docker-compose service name)
    --host mongodatabase                        \
    # user and pass from .env file
    --username root      \
    --password pass      \
    --authenticationDatabase admin              \
    # initial db and collection from .env file
    --db test                 \
    --collection hotels       \
    --type json                                 \
    --file /init-data.json                      \
    --jsonArray                                 \
    # WARNING! drop collection before run import
    --drop                                      \
    # verbose
    -v                                          \