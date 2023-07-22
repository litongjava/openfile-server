# openfile-server

1.upload file   
2.download file  
3.file list  

upload file
```
curl --location --request POST 'http://localhost:8080/upload/litongjava/go/pdf' \
--form 'file=@"C:\\Users\\Administrator\\Downloads\\Go\\Gin Bind Data and Validate.pdf"'
```

file list
```
http://localhost:8080/s/
```