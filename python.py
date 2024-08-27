from pygoridge import create_relay, RPC, SocketRelay
import time 


rpc = RPC(SocketRelay("127.0.0.1", 6001))


start = time.perf_counter()

responce_query = (rpc("App.MakeRequest", (
    "postgres", 
    "host=192.168.32.100 user=sppradmin dbname=etl_datamart sslmode=disable password=TQmGSShiqv_rNUVPeT06", 
    "select * from pden_vol_summary_fact",
    "0",
    "0")))
rpc.close()

# print(responce_query)

delta = time.perf_counter() - start

print(f'''получение токена доступа и запроса к БД составило: {delta}s
количество записей:  {len(responce_query['rows'])}''')
print(responce_query["rows"][-1])