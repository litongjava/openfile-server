# openfile-server

1.upload file   
2.download file  
3.file list  

upload file

windows
```
curl http://localhost:8080/upload/litongjava/go/pdf --form file="@C:\Users\Administrator\Downloads\Go\Gin Bind Data and Validate.pdf"
```

macos
```
curl http://localhost:8080/upload/litongjava/go/pdf --form file="@1.txt"
```

file list
```
http://localhost:8080/s/
```