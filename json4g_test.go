//   donnie4w@gmail.com
//   1.0.0

package json4g

import (
	"fmt"
	"testing"
)

//基本操作
func Test(t *testing.T) {
	jsonStr := `{"STRING": "abcd","NUMBER":123456,"BOOL": false,"STRUCT":{"a":{"b":12345}},"array":["a","b","c"]}`
	node, err := LoadByString(jsonStr)

	if err == nil {
		jnode := node.GetNodeByPath("STRUCT.a.b")
		fmt.Println("STRUCT.a.b>>>>>>>", jnode.ValueNumber)
		jnode.SetValue("1")
		fmt.Println("STRUCT.a.b 修改后的值>>>>>>>", jnode.ValueString)
		fmt.Println("STRING>>>>>>>>>", node.GetNodeByName("STRING").ValueString)
		fmt.Println("NUMBER>>>>>>>>>", node.GetNodeByName("NUMBER").ValueNumber)
		fmt.Println("BOOL>>>>>>>>>>>", node.GetNodeByName("BOOL").ValueBool)
		fmt.Println("array>>>>>>>>>>", node.GetNodeByName("array").ArraysString)

	}
	jsonstr2 := `[{ "id": 1, "name": "n_1" }, { "id": 2, "name": "n_2"}]`
	node, err = LoadByString(jsonstr2)
	if err == nil {
		fmt.Println("jsonstr2>>>>>>>>>", node.ArraysStruct[0].ToJsonNode().GetNodeByName("id").ValueNumber)
	}

	jsonstr3 := `[11,22,33,44,55]` //json数组类型相同时用法
	node, err = LoadByString(jsonstr3)
	fmt.Println("jsonstr3>>>>>>>>>", node.ArraysNumber[4])

	jsonstr4 := `[11,22,33,44,"aa"]` //json数组类型不相同时用法，统一为类型ArraysStruct
	node, err = LoadByString(jsonstr4)
	fmt.Println("jsonstr4>>>>>>>>>", node.ArraysStruct[0].ToJsonNode().ValueNumber, node.ArraysStruct[4].ToJsonNode().ValueString)
}

//获取子节点
func TestGetNode(t *testing.T) {
	str := `{
    "structnodes": [
		{ "node1":"abcd1" , "node2":"abcd2" ,"struct2":{"node11":12345}},
		{ "node3":"abcd3" , "node4":"abcd4" }
		]
	}`
	node, _ := LoadByString(str)
	fmt.Println("struct2>>>>>>>>>", node.GetNodeByPath("structnodes").ArraysStruct[0].ToJsonNode().GetNodeByPath("struct2.node11").ValueNumber)
}

//获取子节点 ：支持路径查询, 节点间用 . 分隔
func TestGetNode2(t *testing.T) {
	jsonStr := `{"STRING": "abcd","NUMBER":123456,"BOOL": false,"STRUCT":{"a":{"b":12345}},"array":["a","b","c"]}`
	node, err := LoadByString(jsonStr)
	if err == nil {
		jnode := node.GetNodeByPath("STRUCT.a.b")
		fmt.Println("STRUCT.a.b>>>>>>>", jnode.ValueNumber)
	}
	fmt.Println("node.AddNode>>>>>>>", node.ToString())
}

//增加节点
func TestAddNode(t *testing.T) {
	str := `{
    "structnodes": [
		{ "node1":"abcd1" , "node2":"abcd2" },
		{ "node3":"abcd3" , "node4":"abcd4" }
		]
	}`
	root, _ := LoadByString(str)
	addnodeString := new(JsonNode)
	addnodeString.Name = "newNodeNameString"
	addnodeString.SetValue("newNodeValue") // 注意，对象值必须通过SetValue()设置
	root.AddNode(addnodeString)
	addnodeNumber := new(JsonNode)
	addnodeNumber.Name = "newNodeNameNumber"
	addnodeNumber.SetValue(123456)
	root.AddNode(addnodeNumber)
	addnodeBool := new(JsonNode)
	addnodeBool.Name = "newNodeNameBool"
	addnodeBool.SetValue(false)
	root.AddNode(addnodeBool)
	fmt.Println("AddNode 后json字符串>>>>>>>>>>>>", root.ToString())
}

//删除节点
func TestDelNode(t *testing.T) {
	str := `{
    "structnodes": [
		{ "node1":"abcd1" , "node2":"abcd2" },
		{ "node3":"abcd3" , "node4":"abcd4" }
		],
	 "delNode":"wuxiaodong"			
	}`
	root, _ := LoadByString(str)
	fmt.Println("删除前>>>>>>>>>>>>", root.ToString())
	root.DelNode("delNode")
	fmt.Println("删除后>>>>>>>>>>>>", root.ToString())
}

//new 节点
func TestNewNode(t *testing.T) {
	jsonroot := new(JsonNode)
	newnode := NowJsonNode("name", "value")
	jsonroot.AddNode(newnode)
	fmt.Println("newnode>>>>>>", jsonroot.ToString())
	newnode2 := NowJsonNode("name2", "value2")
	jsonroot.AddNode(newnode2)
	fmt.Println("newnode2>>>>>", jsonroot.ToString())
	newnode3 := NowJsonNodeByString("name3", `{"node":123456,"node2":"abcd"}`)
	jsonroot.AddNode(newnode3)
	fmt.Println("newnode3>>>>>", jsonroot.ToString())
	newnodeerr := NowJsonNodeByString("name4", `{"node":123456,"node2":"abcd" json error format}`)
	fmt.Println("newnodeerr>>>>", newnodeerr)
	err := jsonroot.AddNode(newnodeerr)
	fmt.Println("err>>>>>>>>", err)
	fmt.Println("jsonroot>>>>>", jsonroot.ToString())
}

//对象转字符串
func TestToString(t *testing.T) {
	str := `{
    "structnodes": [
		{ "node1":"abcd1" , "node2":"abcd2" },
		{ "node3":"abcd3" , "node4":"abcd4" }
		],
	 "delNode":"wuxiaodong",
	  "array":[123,"aabc",{"node1":"abcd1","node2":"abcd2"}]		
	}`
	root, err := LoadByString(str)
	root, err = LoadByString(root.ToString())
	if err == nil {
		fmt.Println("TestToString>>>>>>>>", root.ToString())
	} else {
		fmt.Println("TestToString>>>>>>>>", err)
	}
}
