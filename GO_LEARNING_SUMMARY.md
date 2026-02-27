# Go 语言基础语法速查大纲

这份大纲总结了你在开发 go-zero 项目时遇到的核心语法概念，适合有其他语言背景（如 JS/TS/Java）的开发者快速对照。

## 1. 变量与基础类型

### 变量声明
Go 只有两种主要的声明方式。
```go
// 1. 标准声明 (带类型或自动推导)
var x int = 10
var y = "hello"

// 2. 短变量声明 (只能在函数内部使用，最常用)
name := "zhangsan" 
age := 18
```

### 基础类型
*   **数值**: `int`, `int64`, `float64` (注意：`int` 和 `int64` 不能直接运算，必须强转)。
*   **字符串**: `string` (双引号)。
*   **布尔**: `bool` (`true`, `false`)。
*   **字节**: `byte` (相当于 `uint8`，常用于处理二进制数据，如密码加密 `[]byte`).

### 类型转换 (Type Conversion)
Go 不会隐式转换类型，必须显式转换。
```go
i := 10
f := float64(i)
b := []byte("hello") // string 转 byte 切片
s := string(b)       // byte 切片转 string
```

---

## 2. 复杂数据结构

### 切片 (Slice) - 动态数组
Go 的数组是定长的，但在开发中我们几乎只用切片（Slice）。
```go
// 定义
var nums []int
names := []string{"Alice", "Bob"}

// 追加 (append 会返回新的切片)
nums = append(nums, 1)

// 长度与容量
len(nums) // 长度
cap(nums) // 容量
```

### 映射 (Map) - 字典/哈希表
类似于 JS 的 Object 或 Java 的 HashMap。
```go
// 定义: map[Key类型]Value类型
user := make(map[string]interface{})
user["name"] = "Alice"
user["age"] = 18

// 检查 Key 是否存在 (Comma-ok idiom)
val, ok := user["name"] 
if ok {
    // 存在
}
```

---

## 3. 控制流

### 条件判断 (if)
不需要括号，左花括号必须跟在条件后面。
```go
if x > 10 {
    // ...
} else {
    // ...
}

// 特技: if 可以在判断前执行一句简单的语句 (常用于错误处理)
if err := doSomething(); err != nil {
    return err
}
```

### 循环 (for)
Go 只有 `for`，没有 `while`。
```go
// 1. 经典 C 风格
for i := 0; i < 10; i++ { }

// 2. 相当于 while
for i < 10 { }

// 3. 死循环 (相当于 while(true))
for { }

// 4. 遍历 Slice 或 Map (range)
// i 是索引(key), v 是值
for i, v := range items {
    fmt.Println(i, v)
}
// 只需要值，忽略索引 (使用下划线 _)
for _, v := range items {
}
```

---

## 4. 函数 (Function)

### 基础定义
```go
func Add(a int, b int) int {
    return a + b
}
```

### 多返回值
这是 Go 的一大特色，常用于返回 `(结果, 错误)`。
```go
func Div(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

// 调用
res, err := Div(10, 0)
```

---

## 5. 指针 (Pointer)

在 Go 中，指针主要用于**修改值**或者**避免大对象复制**。

*   `&` (取地址): 获取变量的内存地址。
*   `*` (解引用): 获取地址指向的值。

**为什么 GORM 查询要用 `&user`?**
```go
var user User
// 必须传指针，函数内部才能修改 user 变量的值
db.First(&user) 
```

---

## 6. 结构体 (Struct) & 方法 (Method)

Go 没有 `class`，而是用 `struct` + `method` 来实现面向对象。

### 定义结构体
```go
type User struct {
    Name string
    Age  int
}
```

### 定义方法 (Receiver)
给结构体添加函数。注意 `(u *User)` 称为**接收者**。
```go
// 指针接收者: 可以修改 u 内部的值 (推荐默认使用)
func (u *User) Birthday() {
    u.Age++
}

// 值接收者: 只能读取，修改不会影响原对象
func (u User) SayHello() {
    fmt.Println(u.Name)
}
```

---

## 7. 接口 (Interface) - Duck Typing

Go 的接口是**隐式实现**的。只要一个结构体实现了接口定义的所有方法，它就自动实现了该接口。

```go
type Speaker interface {
    Speak()
}

type Cat struct{}
func (c *Cat) Speak() { fmt.Println("Meow") }

// Cat 指针自动成为了 Speaker
var s Speaker = &Cat{}
```

---

## 8. 错误处理 (Error Handling)

Go 没有 `try-catch`。错误被视为普通的值。

```go
f, err := os.Open("filename.txt")
if err != nil {
    // 处理错误 (通常是返回或打印日志)
    log.Println("打开失败:", err)
    return
}
// 处理成功逻辑
```

---

## 9. JSON 处理 (Tag)

在 web 开发中非常常见。通过反引号 ``` `...` ``` 给结构体字段打标签，告诉 JSON 解析器如何转换。

```go
type LoginReq struct {
    // 解析 JSON 时，去找 "username" 字段赋值给 Username
    Username string `json:"username"` 
    Password string `json:"password"`
}
```

## 10. 模块与包 (Modules & Packages)

*   **Public/Private**: 首字母**大写**是公有的 (Public)，可以在其他包使用；首字母**小写**是私有的 (Private)，只能在当前包使用。
*   **Import**: 引入包的路径。
*   **go.mod**: 类似 `package.json`，管理项目依赖。
