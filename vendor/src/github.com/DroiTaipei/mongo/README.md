# Droi mongo package

## Dependency

* github.com/DroiTaipei/droictx
* github.com/DroiTaipei/droipkg
* github.com/DroiTaipei/mgo

## Initialize

mongo package 預設會使用 droipkg.Logger。

```go
err := mongo.Initialize(mgoConfig)
```

## Config

* Example:

```toml
[database]
[database.mgo.default]
  name = "Mgo"
  user = "bass"
  password = "F4mIfVIYGUBYBiQE"
  authdatabase = "admin"
  max_conn = 5
  host_num = 2
  timeout = 30
  direct = false
  secondary = false

[database.mgo.instance.0]
  host = "10.128.112.181"
  port = 7379
  enabled = true

[database.mgo.instance.1]
  host = "10.128.112.181"
  port = 7380
  enabled = true
```