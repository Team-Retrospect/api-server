echo "(Re)-Building..."
go build

echo "Killing Cassandra..."
sudo pkill -f './cassandra-connector'

echo "Reviving Cassandra..."
sudo ./cassandra-connector &

