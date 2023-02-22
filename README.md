



# thf

### <u>t</u>ypical <u>h</u>ttp.Handler <u>f</u>unc

对`http`标准库中 `func(http.ResponseWriter, *http.Request)` 函数进行泛型处理

开发者只需关注输入，输出对象的声明，及错误处理方式

## Installation

golang version 1.81+

`go get github.com/otk-final/thf`

## Usage & Example

### 默认处理方式

```go
func api(writer http.ResponseWriter, request *http.Request){
  //读取数据
  ioutil.ReadAll(request.Body)
  //反序列化 struct
  json.Unmarshal(body, &t)
  
  
  //业务处理
  //...
  
  
  //序列化
  bytes,_ := json.Marshal(&t)
  //响应数据
  writer.Write(bytes)
}

//路由
mux := http.NewServeMux()
mux.HandleFunc("/a",apiA)
mux.HandleFunc("/b",apiB)
mux.HandleFunc("/c",apiC)

```

不同api需要独立处理参数的序列化，和反序列化，通常一个单体应用参数序列化方式是同一的，如`json.Unmarshal` or `json.Marshal`

现将繁琐的操作交由 `thf` 框架层处理，其支持 `7` 种不同方法声明

1. ```go
   func(http.ResponseWriter, *http.Request, T) (R, error)		
   ```

2. ```go
   func(http.ResponseWriter, *http.Request, T) R
   ```

3. ```go
   func(http.ResponseWriter, *http.Request, T) error
   ```

4. ```go
   func(http.ResponseWriter, *http.Request, T)
   ```

5. ```go
   func(http.ResponseWriter, *http.Request) (R, error)
   ```

6. ```go
   func(http.ResponseWriter, *http.Request) R
   ```

7. ```go
   func(http.ResponseWriter, *http.Request) error
   ```

   [^注意]: 函数如未声明任何返回类型，框架不会做任何响应处理，需由开发者在函数体中自己实现响应输出过程

### thf处理方式

```go
// typical http.Handler func
func typicalApi(writer http.ResponseWriter, request *http.Request, in *Foo) (*Bar, error) {
	//TODO		
}

mux := http.NewServeMux()
mux.HandleFunc("/a",thf.Wrap(typicalApi).Func())
mux.HandleFunc("/b",thf.WrapIO(typicalApi).Func())
mux.HandleFunc("/c",thf.WrapIE(typicalApi).Func())
mux.HandleFunc("/d",thf.WrapI(typicalApi).Func())
mux.HandleFunc("/e",thf.WrapOE(typicalApi).Func())
mux.HandleFunc("/f",thf.WrapO(typicalApi).Func())
mux.HandleFunc("/g",thf.WrapE(typicalApi).Func())

// IOE 是 in out error 缩写
```

### 自定义编解码

```go
//实现编解码接口

type Decoder[T any] interface {
  //参数解析
  Decode(*http.Request) (T, error)
}

type Encoder[R any] interface {
  //正常响应
  Out(http.ResponseWriter, *http.Request, R)
  //异常响应
  Error(http.ResponseWriter, *http.Request, error)
}


mux.HandleFunc("/custom",thf.Wrap(typicalApi).FuncBy(xxxDecoder,xxxEncoder))
```

