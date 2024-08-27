from pygoridge import create_relay, RPC, SocketRelay
import time 


rpc = RPC(SocketRelay("127.0.0.1", 6001))


start = time.perf_counter()

responce_query = (rpc("App.MakeRequest", (
    "postgres", 
    "host=localhost user=nikolay dbname=bi sslmode=disable password=12345", 
    "select * from bi_charts",
    "0",
    "0")))
rpc.close()

# print(responce_query)

delta = time.perf_counter() - start

print(f'''получение токена доступа и запроса к БД составило: {delta}s
количество записей:  {len(responce_query['rows'])}''')
print(responce_query["rows"][-1])