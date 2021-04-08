### socks5-client
将http(s)请求转发到socks-server上

#### 适用场景
在某些httpclient不支持设置代理的情况下，这些请求就不能走代理通道，
这时可以通过改host的方式，将这个域名的请求指向socks-client,
socks-client会将这部分流量转发到socks5 server上

#### 使用方式
1. go build 
2. sudo ./socks-client targetIp
3. 示例 sudo ./socks-client 

待办列表
- [x] 支持将单一域名的请求转发到socks5 tunnel上
- [ ] 支持多域名的请求转发到socks5上
- [ ] 支持配置文件指定特定域名的实际目标IP
- [ ] 支持socks-server 和socks-client在单机情况下可配制socks-server的dns resolver 

