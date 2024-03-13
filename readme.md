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
curl --location --request POST http://localhost:8080/upload/litongjava/images --form file=@"q4.png"
```
or

```
curl --request POST http://localhost:8080/upload/litongjava/images/ --form file=@"q4.png"
```
or
```
curl -X POST http://localhost:8080/upload/litongjava/images/ --form file=@"q4.png"
```
or
```
curl http://localhost:8080/upload/litongjava/images/ --form file=@"q4.png"
```


file list
```
http://localhost:8080/s/
```

