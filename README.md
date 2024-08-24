# go-easymq
easymq golang client

## **Example**
```go
client,err:=goeasymq.Connect("easymq","localhost:7777")
if err!=nil{
  panic(err)
}
err=client.Push(context.Background(), "testGo", "hello easymq")
if err!=nil{
  panic(err)
}
content,err:=client.ReadLatest(context.Background(),"testGo")
if err!=nil{
  panic(err)
}
println(content)
```