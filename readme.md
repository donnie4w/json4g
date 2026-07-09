# Json Go库

轻量级、零依赖的 Go 语言 JSON 动态操作库，提供类 DOM 的节点式 API，支持解析、构建、查询、修改、序列化全流程操作，无需预定义结构体即可灵活处理任意 JSON 数据。

## ✨ 特性
- **零外部依赖**：仅基于 Go 标准库 `encoding/json`，无第三方包冲突风险
- **动态无结构体**：无需提前定义 `struct`，任意不规则 JSON 直接读写
- **智能路径查询**：支持 `a.b.c` 点分嵌套 + `list[0].name` 数组下标混合路径
- **全类型覆盖**：对象、同类型数组、混合数组、数字、字符串、布尔、null
- **完整数组操作**：追加、插入、修改索引、删除元素、清空数组全套 API
- **高性能序列化**：`bytes.Buffer` 流式拼接，规避字符串频繁拷贝，GC 压力更低
- **安全容错**：全方法内置 panic 捕获，非法下标/路径/类型不会直接崩溃服务
- **读写一体化**：区别于只读类 JSON 库，支持解析后修改并重新生成合法 JSON

## 📦 安装
```bash
go get github.com/donnie4w/json4g
```

## 🚀 快速开始
### 1. 解析 JSON 字符串，深度路径读取
```go
package main

import (
	"fmt"
	"github.com/donnie4w/json4g"
)

func main() {
	jsonStr := `{"name":"alice","age":25,"hobbies":["reading","coding"],"address":{"city":"beijing"}}`
	
	root, err := json4g.LoadByString(jsonStr)
	if err != nil {
		panic(err)
	}

	// 单层字段读取
	fmt.Println(root.GetNodeByName("name").ValueString) // alice

	// 多层嵌套路径
	fmt.Println(root.GetNodeByPath("address.city").ValueString) // beijing

	// 数组下标混合路径
	fmt.Println(root.GetNodeByPath("hobbies[1]").ValueString) // coding
}
```

### 2. 从零构建 JSON 对象
```go
func buildJSON() {
	root := new(json4g.JsonNode)
	_ = root.SetValue(map[string]interface{}{}) // 初始化为对象节点

	// 添加基础字段
	root.AddNode(json4g.NowJsonNode("name", "bob"))
	root.AddNode(json4g.NowJsonNode("age", 30))

	// 添加数组
	tagNode := json4g.NowJsonNode("tags", []string{"go", "json"})
	root.AddNode(tagNode)

	fmt.Println(root.ToString())
	// {"name":"bob","age":30,"tags":["go","json"]}
}
```

### 3. 对象/数组动态修改
```go
func modifyJSON() {
	root, _ := json4g.LoadByString(`{"name":"test","list":[1,2,3]}`)

	// 修改普通字段
	_ = root.GetNodeByName("name").SetValue("newName")

	// 数组操作
	list := root.GetNodeByName("list")
	_ = list.SetIndex(1, 200) // 修改下标1
	_ = list.Append(400)      // 尾部追加
	_ = list.Insert(0, 0)     // 指定位置插入
	_ = list.RemoveIndex(2)   // 删除下标元素

	// 删除对象字段
	_ = root.DelNode("name")
	fmt.Println(root.ToString())
}
```

### 4. 路径存在性判断
```go
func checkPath() {
	root, _ := json4g.LoadByString(`{"user":{"info":{"phone":"13800000000"}}}`)
	fmt.Println(root.IsExist("user.info.phone")) // true
	fmt.Println(root.IsExist("user.info.email")) // false
}
```

## 📚 API 文档
### 节点创建 & 解析
| 方法 | 说明 |
|------|------|
| `LoadByString(jsonStr string) (*JsonNode, error)` | JSON 字符串解析为根节点 |
| `NowJsonNode(name string, value interface{}) *JsonNode` | 创建命名节点，自动识别值类型 |
| `NowJsonNodeByString(name string, jsonStr string) *JsonNode` | JSON 片段生成子节点 |

### 对象(STRUCT)操作
| 方法 | 说明 |
|------|------|
| `AddNode(node *JsonNode) error` | 新增子字段 |
| `DelNode(name string) error` | 删除指定字段 |
| `GetNodeByName(name string) *JsonNode` | 获取一级子节点 |
| `GetNodeByPath(path string) *JsonNode` | 点分+下标混合路径深度查询 |
| `GetChildNames() []string` | 获取所有一级字段名 |
| `IsExist(path string) bool` | 判断路径是否存在 |

### 数组专用操作（仅数组类型节点生效）
| 方法 | 说明 |
|------|------|
| `SetIndex(index int, value interface{}) error` | 修改指定下标元素 |
| `Append(value interface{}) error` | 数组尾部追加 |
| `Insert(index int, value interface{}) error` | 指定下标插入元素 |
| `RemoveIndex(index int) error` | 删除指定下标元素 |
| `ClearArray() error` | 清空整个数组 |

### 通用基础方法
| 方法 | 说明 |
|------|------|
| `SetValue(value interface{}) error` | 重置当前节点值 |
| `ToString() string` | 序列化输出合法 JSON |
| `Len() int` | 对象字段数 / 数组长度 |
| `Type() NODETYPE` | 获取当前节点类型枚举 |
| `IsNull() bool` | 判断节点是否为 null |

### 路径语法规则
1. 对象嵌套：`user.info.address`
2. 普通数组：`hobbies[0]`
3. 对象数组：`users[2].phone`
4. 多维数组：`matrix[0][1]`
5. 根节点为数组：`[0].id`

### 节点类型枚举 NODETYPE
```go
const (
	STRUCT      // JSON 对象 {}
	STRUCTARRAY // 混合/对象数组 [{},1,"abc"]
	NUMBERARRAY // 纯数字数组 [1,2,3]
	BOOLARRAY   // 纯布尔数组 [true,false]
	STRINGARRAY // 纯字符串数组 ["a","b"]
	NUMBER      // 数字值
	BOOL        // 布尔值
	STRING      // 字符串值
	NULL        // null
)
```

## 🧪 测试 & 基准测试
### 执行单元测试
```bash
go test -v
```
### 执行性能基准
```bash
go test -bench=. -benchmem
```

### 基准参考
| 测试场景 | 单次耗时 | 内存分配 |
|--------|----------|----------|
| 小型JSON解析 | ~300 ns | ~320 B |
| 复杂嵌套JSON解析 | ~1800 ns | ~2100 B |
| 浅层路径查询 | ~60 ns | ~0 B |
| 带下标混合路径查询 | ~180 ns | ~0 B |
| 大对象序列化 | ~900 ns | ~512 B |

> 数值随硬件、Go 版本略有浮动

## 适用场景 
### ✅ 推荐使用 json4g 的场景
1. **JSON 结构动态/不固定，不想定义大量 struct**
   Go 标准库 `encoding/json` 强依赖预定义结构体，不规则、动态字段场景需要写大量嵌套 struct；
   json4g 直接解析后通过路径读写，零模板代码。
   ```go
   // 标准库原生多层取值：大量类型断言，极易panic
   data, _ := json.Unmarshal([]byte(raw), &map[string]any{})
   u, _ := data["user"].(map[string]any)
   info, _ := u["info"].(map[string]any)
   phone := info["phone"].(string)

   // json4g 一行搞定，自动判空容错
   phone := root.GetNodeByPath("user.info.phone").ValueString
   ```
2. **需要修改 JSON 数组（插入/删除/替换）**
   原生 map+slice 操作数组需要手动处理切片截取、类型转换，代码臃肿；
   json4g 封装全套数组 API，自带越界校验与类型匹配校验。
3. **通用网关、爬虫、日志清洗、运维脚本**
   一次性工具、通用中间件，JSON 结构不可预知，追求快速开发。
4. **动态构造返回报文、动态配置生成**
   运行时按需增删字段，无需提前定义固定结构体。
5. **仅读取 JSON 局部字段，完整结构体性价比极低**
   只需提取深层某几个字段，不需要完整映射整个 JSON。

### ❌ 不推荐使用，优先标准库 struct json 的场景
1. **超高并发核心服务，JSON 结构永久固定**
   支付、存储、RPC 等百万 QPS 服务，标准库代码生成性能远超本库；
   json4g 内部维护独立节点树，存在一层额外内存与转换开销。
2. **需要强类型校验、JSON Schema、字段必填约束**
   标准库 tag 支持 `required`、`omitempty`、数值范围校验；
   json4g 为动态弱类型，无编译期校验，校验逻辑需要手动编写。
3. **业务核心是结构体与 JSON 互相映射转换**
   大量实体结构体流转场景，本库会增加额外转换成本。

