# tbl-cli

打表代码生成工具

## 打表代码客户端生成

```bash
tbl-cli --server dev:Zqa5Q5pYgvmT@192.168.110.112:3306 --db treasure_base --config ./config/ignore-client.yaml --pkg cdb --go-pkg gitlab.walrus.cs/treasure/treasure-pb/cbpb --out ../tbl/internal/gdb/cdb
```

## 打表代码服务端生成

```bash
tbl-cli --server dev:Zqa5Q5pYgvmT@192.168.110.112:3306 --db treasure_base --config ./config/ignore-server.yaml --pkg sdb --go-pkg gitlab.walrus.cs/treasure/treasure-pb/sbpb --out ../tbl/internal/gdb/sdb
```
