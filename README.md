# Python-Golang RPC

my own expirience with database executor on golang with rpc call from python

# Variables

## For retrieve system_user token
```
SYSTEM_USER
SYSTEM_USER_PASSWORD
AUTH_TOKEN_URL
```

# Setup

## Compile and run golang part
```
(venv) nikolay@nikolay-BOD-WXX9:~/projects/python-go-rpc$ cd go_sql
(venv) nikolay@nikolay-BOD-WXX9:~/projects/python-go-rpc/go_sql$ go build goGetSQL.go 
(venv) nikolay@nikolay-BOD-WXX9:~/projects/python-go-rpc/go_sql$ ./goGetSQL 
```
## Test it with python
```
(venv) nikolay@nikolay-BOD-WXX9:~/projects/python-go-rpc$ python3 python.py 
получение токена доступа и запроса к БД составило: 0.0050461950013414025s
количество записей:  12
```