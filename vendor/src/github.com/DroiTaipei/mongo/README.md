# Droi mongo package
## Dependency
* github.com/DroiTaipei/droictx
* github.com/DroiTaipei/droipkg
* github.com/DroiTaipei/mgo

## Initialize
mongo package 預設會使用 droipkg.Logger 紀錄 access log，所以需要在初始化設定 droipkg.Logger 。

## Config
注意：authdatabase 和 database 的差異，database是操作用的資料庫。
* Example Config: 
```
[database]
[database.mgo.default]
  name = "Mgo"
  user = "bass"
  password = "F4mIfVIYGUBYBiQE"
  authdatabase = "admin"
  database = "baas"
  max_conn = 10
  max_idle = 15
  host_num = 2
  timeout = 30
  sharduser= "sharder"
  shardpassword = "izihGm5UD7KZnCRzbAkFX6PLdtUXK7n3"
  direct = false
```