## echo-server

Simple http-server with only two endpoints  
&nbsp;   
Build:
```
git clone git@github.com:awnzl/echo-server.git
cd echo-server
make
```
&nbsp;  
To run just execute
```
./bin/server
```
&nbsp;  
Usage:  
```
curl localhost:8080
```
```
curl -X POST localhost:8080/echo -d '{"word": "ha-ha-ha"}'
```
