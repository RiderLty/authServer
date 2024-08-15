# 简易认证服务器

配合nginx实现同一根域名下网站统一认证

设置环境变量DOMAIN为跟域名地址，例如 ```.abc.com```

## 接口

### getcode接口

获取cookie

### auth接口

如果没有cookie或者cookie错误，返回403

成功则302到redirect参数的地址

## Nginx & Cloudflare Tunnel

authServer使用nginx的基本认证，配置参考如下

```Nginx
server {
    listen 80;
    server_name auth.abc.com;
    location / {
        auth_basic "Restricted Access";
        auth_basic_user_file /etc/nginx/passwds.txt;
        proxy_pass http://192.168.3.3:8080;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header Host $host:$server_port;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Host $host;
        proxy_redirect http:// $scheme://;
    }
}
```

本地服务使用nginx反向代理，设置域名为 \[server\].abc.com

cloudflare Tunnel创建对应的隧道

配置参考如下
``` Nginx
server {
    listen 80;
    server_name rsshub.abc.com;
    location / {
        proxy_pass http://192.168.3.3:1200;
        auth_request /auth;
        error_page 403 = @error403;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;
        proxy_set_header Host $host:$server_port;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Host $host;
        proxy_redirect http:// $scheme://;
    }
    
    location @error403 {
      return 302 "https://auth.abc.com/getcode?redirect=$scheme://$host$request_uri";
    }

    location /auth {
      internal;
      proxy_pass http://192.168.3.3:8080/auth;
      proxy_set_header Content-Length "";
      proxy_set_header X-Original-URI $request_uri;
      proxy_set_header Cookie $http_cookie;
    }
}```