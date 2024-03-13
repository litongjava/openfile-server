# openfile-server
## introduction
1.upload file   
2.download file  
3.file list  

## install
```
docker build -t litongjava/openfile-server:1.0.0 .
```
run
```
docker run -dit --name=openfile-server -p 8080:80 litongjava/openfile-server:1.0.0
```



## usage
upload file

windows
```
curl --location --request POST http://192.168.3.9/upload/litongjava/images --form file=@"q4.png"
```
linux
```shell
curl --location --request POST http://192.168.3.9:8080/upload/litongjava/images --form file=@"graalvm-jdk-21_linux-x64_bin.tar.gz"
```
macos
```
curl http://localhost:8080/upload/litongjava/go/pdf --form file="@1.txt"
```

file list
```
http://localhost:8080/s/
```

