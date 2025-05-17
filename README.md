# hello_firewall

#### 介绍

可视化ip防火墙管理系统，使用go语言编写相关接口，底层防火墙驱动是iptables，提供API调用

#### 安装

```
go mod tidy
```

#### 运行
使用go运行

```
go run main.go serve
```
使用gin运行（热重载）
```
gin serve run main.go
```