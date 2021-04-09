rm ../out.*
go build ../../depcarto.go
./depcarto server -c server.json &
export pid=`echo $!`
curl -s -d'@front1.json' localhost:8088/server -X POST
curl -s -d'@back1.json' localhost:8088/server -X POST
curl -s -d'@back2.json' localhost:8088/server -X POST
curl -s -d'@db.json' localhost:8088/server -X POST
curl -s -d'@es.json' localhost:8088/server -X POST
curl -s localhost:8088/servers > ../out.json
curl -s localhost:8088/servers?format=dot > ../out.dot
curl -s localhost:8088/servers?format=jpg > ../out.jpg
kill $pid
rm depcarto