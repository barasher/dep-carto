rm out.*
go build ../depcarto.go
./depcarto server -c servers/server.json &
export pid=`echo $!`
curl -s -d'@servers/front1.json' localhost:8088/server -X POST
curl -s -d'@servers/back1.json' localhost:8088/server -X POST
curl -s -d'@servers/back2.json' localhost:8088/server -X POST
curl -s -d'@servers/db.json' localhost:8088/server -X POST
curl -s -d'@servers/es.json' localhost:8088/server -X POST
curl -s localhost:8088/servers > out.json
curl -s localhost:8088/servers?format=dot > out.dot
curl -s localhost:8088/servers?format=jpg > out.jpg
#kill $pid
rm depcarto
