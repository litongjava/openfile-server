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
location /file {
  # 处理 OPTIONS 请求
  if ($request_method = 'OPTIONS') {
    add_header Access-Control-Allow-Origin *;
    add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
    add_header Access-Control-Allow-Headers 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range';
    add_header Content-Length 0;
    add_header Content-Type text/plain;
    return 204;
  }

  add_header Access-Control-Allow-Origin *;
  add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
  add_header Access-Control-Allow-Headers 'DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range';
  add_header Access-Control-Expose-Headers 'Content-Length,Content-Range';
  
  root /data/apps/openfile-server;
  
  # 开启缓存
  proxy_cache file_cache;
  proxy_cache_valid 200 1h;
  proxy_cache_valid 404 1m;
  add_header X-Cache-Status $upstream_cache_status;

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
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=file_cache:10m max_size=10g inactive=60m use_temp_path=off;

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

## Thumbnail Configuration - Using Nginx `image_filter` Module

This document explains how to use Nginx's `image_filter` module to dynamically generate image thumbnails by specifying width and height in the request.

### 1. Check if the `image_filter` module is installed

To handle image resizing in Nginx, the `image_filter` module is required. You can check if this module is installed with the following command:

```bash
nginx -V 2>&1 | grep image_filter
```

If the output includes `--with-http_image_filter_module`, the module is installed. If not, you can recompile Nginx with this module enabled or use the official precompiled package provided by Nginx.
### 2. load ngx_http_image_filter_module
```cgo
--with-http_image_filter_module=dynamic
```
if dynamic please load module
```cgo
load_module modules/ngx_http_image_filter_module.so;

```
### 3. Nginx Configuration

Below is an example configuration using the `image_filter` module to resize images:

```nginx
location ~* /file/(.*)_(\d+)x(\d+)\.(jpg|jpeg|png|gif)$ {
  # Extract the image name, width, and height
  set $image_name $1;
  set $width $2;
  set $height $3;

  # Apply image resizing
  image_filter resize $width $height;
  
  # Set JPEG quality to 70 to reduce file size (default is 75)
  image_filter_jpeg_quality 70;

  # Configure buffer size to avoid memory overflow when processing large images
  image_filter_buffer 100M;

  # Set the returned image type based on the file extension
  default_type image/$4;

  # Use alias instead of root to match the actual image path
  alias /data/apps/openfile-server/file/$image_name.$4;
}
```

#### Explanation:

- **location matching rule**: Matches image requests starting with `/file/`, where the name includes width and height parameters (e.g., `_100x100`). In the regular expression, `$1` represents the image name, `$2` and `$3` represent width and height, and `$4` represents the image extension (`jpg`, `jpeg`, `png`, `gif`).

- **`set` directive**: Assigns the captured groups from the regular expression to `$image_name` (original image name), `$width` (width), and `$height` (height).

- **`image_filter resize`**: Uses the `image_filter` module to resize the image to the specified width and height.

- **`image_filter_jpeg_quality`**: Adjusts the JPEG image compression quality. It's set to 70 here to balance image quality and file size.

- **`image_filter_buffer`**: Specifies the buffer size to prevent memory overflow when processing large images. A larger value, such as `100M`, is recommended.

- **`default_type`**: Returns the correct `MIME` type based on the image extension, ensuring the browser renders the image correctly.

- **`alias` directive**: Maps to the actual image path on the server using `alias` instead of `root` to avoid path confusion. In this case, images are stored in the `/data/apps/openfile-server/file/` directory.

### 4. Testing Access

You can test the setup using the following URLs:

- **Original image access**:  
  `http://192.168.3.9:8568/file/images/200-dpi.png`

- **Thumbnail access (100x100 size)**:  
  `http://192.168.3.9:8568/file/images/200-dpi_100x100.png`

### 5. Notes

- The `image_filter` module dynamically processes images, which may consume significant CPU and memory resources, especially in high-concurrency environments. It is recommended to implement a caching mechanism or pre-generate thumbnails to reduce server load.
- For large-scale image resizing, consider caching frequently used thumbnails in a CDN or on disk to further reduce server load.

With this configuration, Nginx can dynamically generate thumbnails based on the requested width and height, allowing users to access images in different sizes efficiently.