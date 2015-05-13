golang 的 json处理库
json4g 提供了json的简便处理方法

方法介绍
1，LoadByString(string)    参数为json字符串，返回 JsonNode 对象指针
2，NowJsonNode(string,interface{})    参数节点名
3，NowJsonNodeByString(string,string)     参数 json字符串 如： {"a":"b","c":123}
4，AddNode(*JsonNode)     为某节点增加子节点 
5，DelNode(string)        某节点 删除指定名称子节点
6，ToJsonNode()          转JsonNode 指针对象
7，SetValue()            节点设置 节点值，可以为 数字，字符串，bool值，数组
8，GetNodeByPath(string) 通过路径查询 节点 如： {"a":{"b":{"c":123,"d":true}}} 节点c对象GetNodeByPath("a.b.c")
9，GetNodeByName(string) 通过节点名 查询 节点对象
10,ToString()            JsonNode对象转换为字符串
具体操作请参考 测试类 json4g_test.go

部分方法示例：

	jsonStr := `{"STRING": "abcd","NUMBER":123456,"BOOL":false,"STRUCT":{"a":{"b":12345}},"array":["a","b","c"]}`
	node, err := LoadByString(jsonStr)
	if err == nil {
		jnode := node.GetNodeByPath("STRUCT.a.b")
		fmt.Println("STRUCT.a.b>>>>>>>", jnode.ValueNumber)
	}
	node.DelNode("NUMBER")  //删除
	node.AddNode(NowJsonNode("name", "value"))  //增加
	node.ToString()          // 转化为字符串

有问题或建议欢迎 email : donnie4w@gmail.com