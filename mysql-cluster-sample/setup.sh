docker-compose up -d

sleep 5

SCHEMA_SQL=$(cat schema.sql)
SCHEMA_SQL=$(echo -n $SCHEMA_SQL)
docker-compose exec mysql-server-1 mysql --user=root --password=pass test_db -e "$SCHEMA_SQL"
