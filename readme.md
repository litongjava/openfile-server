# openfile-server

## introduction

1.upload file   
2.download file  
3.file list

## install
### docker
```
docker build -t litongjava/openfile-server:1.0.0 .
```

run

```
docker run -dit --name=openfile-server -p 8080:80 litongjava/openfile-server:1.0.0
```
### cmd
```shell
go build
```
### openfile-server.service

```
vi /lib/systemd/system/openfile-server.service
```

```
[Unit]
Description=HTTP Server
After=network.target

[Service]
Type=simple
User=root
Restart=on-failure
RestartSec=5s
WorkingDirectory=/data/apps/openfile-server
ExecStart=/usr/local/bin/openfile-server

[Install]
WantedBy=multi-user.target
```

```
systemctl enable openfile-server
systemctl start openfile-server
systemctl status openfile-server
```

### 配置nginx
/etc/nginx/common_file_locations.conf
```shell
common_file_locations.conf
location /file {
  root /data/apps/openfile-server;
}

location /s {
  root /data/apps/openfile-server;
  autoindex on;
  autoindex_exact_size on;
  autoindex_localtime on;
  charset utf-8,gbk;
}
  
location /ping {
  proxy_pass http://file_server;
  proxy_pass_header Set-Cookie;
  proxy_set_header Host $host:$server_port;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
  error_log  /var/log/nginx/backend.error.log;
  access_log  /var/log/nginx/backend.access.log;
}

location /upload {
  proxy_pass http://file_server;
  proxy_pass_header Set-Cookie;
  proxy_set_header Host $host:$server_port;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
  error_log  /var/log/nginx/backend.error.log;
  access_log  /var/log/nginx/backend.access.log;
}

location /uploadImg {
  proxy_pass http://file_server;
  proxy_pass_header Set-Cookie;
  proxy_set_header Host $host:$server_port;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
  error_log  /var/log/nginx/backend.error.log;
  access_log  /var/log/nginx/backend.access.log;
}

location /uploadVideo {
  proxy_pass http://file_server;
  proxy_pass_header Set-Cookie;
  proxy_set_header Host $host:$server_port;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
  error_log  /var/log/nginx/backend.error.log;
  access_log  /var/log/nginx/backend.access.log;
}

location /uploadMp3 {
  proxy_pass http://file_server;
  proxy_pass_header Set-Cookie;
  proxy_set_header Host $host:$server_port;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
  error_log  /var/log/nginx/backend.error.log;
  access_log  /var/log/nginx/backend.access.log;
}

location /uploadDoc {
  proxy_pass http://file_server;
  proxy_pass_header Set-Cookie;
  proxy_set_header Host $host:$server_port;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_set_header Upgrade $http_upgrade;
  proxy_set_header Connection "upgrade";
  error_log  /var/log/nginx/backend.error.log;
  access_log  /var/log/nginx/backend.access.log;
}
```

```shell
upstream file_server {
  server 127.0.0.1:9000;
}


server {
  listen 8568;
  server_name localhost;
  index index.htm index.html;

  include common_file_locations.conf;
}
```

## usage

upload file

```shell
curl http://localhost:8080/upload/litongjava/images/ -F file=@"q4.png"
```

file list

```
http://localhost:8080/s/
```

