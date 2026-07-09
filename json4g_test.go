// Copyright (c) 2016, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// donnie4w@gmail.com

package json4g

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestJson4g(t *testing.T) {
	// ==================== 1. 解析 JSON ====================
	jsonStr := `{
		"user": {
			"name": "Alice",
			"age": 28,
			"active": true,
			"scores": [90, 85, 95, 88],
			"tags": ["go", "json", "fast"],
			"address": {
				"city": "Beijing",
				"detail": null
			}
		},
		"meta": {
			"total": 1,
			"page": 1
		}
	}`

	root, err := LoadByString(jsonStr)
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}
	fmt.Println("✅ 解析成功！根节点类型:", root.Type())

	// ==================== 2. 路径查询 ====================
	fmt.Println("\n--- 路径查询 ---")
	fmt.Println("用户姓名:", root.GetNodeByPath("user.name").ValueString)
	fmt.Println("用户年龄:", root.GetNodeByPath("user.age").ValueNumber)
	fmt.Println("第2个分数:", root.GetNodeByPath("user.scores[1]").ValueNumber)
	fmt.Println("城市:", root.GetNodeByPath("user.address.city").ValueString)
	fmt.Println("是否存在 tags[2]:", root.IsExist("user.tags[2]"))
	fmt.Println("不存在的字段:", root.GetNodeByPath("user.not_exist")) // nil

	// ==================== 3. 节点操作 ====================
	fmt.Println("\n--- 修改操作 ---")
	// 添加新字段
	root.GetNodeByPath("user").AddNode(NowJsonNode("email", "alice@example.com"))

	// 修改数组
	scores := root.GetNodeByPath("user.scores")
	scores.SetIndex(0, 99) // 修改第一个分数
	scores.Append(100)     // 追加
	scores.Insert(2, 77)   // 在索引2插入

	fmt.Println("修改后分数数组:", scores.ToString())

	// 删除字段
	root.GetNodeByPath("meta").DelNode("page")
	fmt.Println("删除 meta.page 后:", root.GetNodeByPath("meta").ToString())

	// ==================== 4. 原生 Go 类型转换 ====================
	fmt.Println("\n--- 原生类型转换 ---")
	arrNode := NowJsonNode("numbers", []int{1, 2, 3, 4, 5})
	fmt.Println("int切片转NUMBERARRAY:", arrNode.ToString())

	// ==================== 5. 序列化 ====================
	fmt.Println("\n--- 最终序列化结果 ---")
	finalJSON := root.ToString()
	fmt.Println(finalJSON)

	// 验证合法性
	var validate interface{}
	if jsonErr := json.Unmarshal([]byte(finalJSON), &validate); jsonErr == nil {
		fmt.Println("✅ 序列化结果是合法 JSON！")
	} else {
		fmt.Println("❌ 序列化结果非法:", jsonErr)
	}

	// ==================== 6. 其他辅助方法 ====================
	fmt.Println("\n--- 辅助方法 ---")
	fmt.Println("user 子节点数量:", root.GetNodeByPath("user").Len())
	fmt.Println("user 的子字段:", root.GetNodeByPath("user").GetChildNames())
	fmt.Println("根节点是否为 STRUCT:", root.Type() == STRUCT)
}

func TestJson4g_Full(t *testing.T) {
	t.Run("Parse_Basic", func(t *testing.T) {
		node, err := LoadByString(`{"name":"Bob","age":30,"active":true,"score":95.5,"nil":null}`)
		if err != nil {
			t.Fatal(err)
		}
		if node.GetNodeByPath("name").ValueString != "Bob" ||
			node.GetNodeByPath("age").ValueNumber != 30 ||
			!node.GetNodeByPath("active").ValueBool ||
			node.GetNodeByPath("score").ValueNumber != 95.5 ||
			!node.GetNodeByPath("nil").IsNull() {
			t.Error("基本类型解析失败")
		}
	})

	t.Run("Parse_Array", func(t *testing.T) {
		node, _ := LoadByString(`[1, "two", true, null, {"id":999}]`)
		if node.NType != STRUCTARRAY {
			t.Error("混合数组应为STRUCTARRAY")
		}
		if node.Len() != 5 {
			t.Error("数组长度错误")
		}
	})

	t.Run("Path_Query", func(t *testing.T) {
		node, _ := LoadByString(`{"user":{"scores":[10,20,30],"info":{"city":"SH"}}}`)

		tests := []struct {
			path     string
			expected interface{}
		}{
			{"user.scores[1]", 20.0},
			{"user.info.city", "SH"},
			{"user.scores[5]", nil}, // 越界
		}

		for _, tc := range tests {
			n := node.GetNodeByPath(tc.path)
			if tc.expected == nil && n != nil {
				t.Errorf("路径 %s 应返回 nil", tc.path)
			}
		}
	})

	t.Run("Modify", func(t *testing.T) {
		node, _ := LoadByString(`{"arr":[1,2,3]}`)

		arr := node.GetNodeByPath("arr")
		arr.SetIndex(0, 99)
		arr.Append(4)
		arr.Insert(1, 88)

		if arr.Len() != 5 || arr.GetNodeByPath("[0]").ValueNumber != 99 {
			t.Error("数组修改失败")
		}

		// 添加字段
		node.GetNodeByPath("").AddNode(NowJsonNode("newkey", "value"))
		if !node.IsExist("newkey") {
			t.Error("AddChild 失败")
		}
	})

	t.Run("RoundTrip", func(t *testing.T) {
		original := `{"a":1,"b":[true,false],"c":{"d":"e"}}`
		node1, _ := LoadByString(original)
		str := node1.ToString()

		var v interface{}
		if err := json.Unmarshal([]byte(str), &v); err != nil {
			t.Fatalf("序列化非法: %v", err)
		}

		node2, _ := LoadByString(str)
		if node2.Len() != node1.Len() {
			t.Error("往返不一致")
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		// 空对象、空数组
		_, err := LoadByString(`{}`)
		if err != nil {
			t.Error(err)
		}
		_, err = LoadByString(`[]`)
		if err != nil {
			t.Error(err)
		}

		// 深层路径
		deep, _ := LoadByString(`{"a":{"b":{"c":{"d":42}}}}`)
		if deep.GetNodeByPath("a.b.c.d").ValueNumber != 42 {
			t.Error("深层路径失败")
		}
	})
}

// ========== 基础类型解析测试 ==========
func TestLoadByString_BasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		jsonStr  string
		wantType NODETYPE
		check    func(*testing.T, *JsonNode)
	}{
		{
			name:     "string value",
			jsonStr:  `"hello world"`,
			wantType: STRING,
			check: func(t *testing.T, n *JsonNode) {
				if n.ValueString != "hello world" {
					t.Errorf("want 'hello world', got '%s'", n.ValueString)
				}
			},
		},
		{
			name:     "number value",
			jsonStr:  `12345.67`,
			wantType: NUMBER,
			check: func(t *testing.T, n *JsonNode) {
				if n.ValueNumber != 12345.67 {
					t.Errorf("want 12345.67, got %v", n.ValueNumber)
				}
			},
		},
		{
			name:     "bool true",
			jsonStr:  `true`,
			wantType: BOOL,
			check: func(t *testing.T, n *JsonNode) {
				if !n.ValueBool {
					t.Error("want true, got false")
				}
			},
		},
		{
			name:     "bool false",
			jsonStr:  `false`,
			wantType: BOOL,
			check: func(t *testing.T, n *JsonNode) {
				if n.ValueBool {
					t.Error("want false, got true")
				}
			},
		},
		{
			name:     "null value",
			jsonStr:  `null`,
			wantType: NULL,
			check: func(t *testing.T, n *JsonNode) {
				if !n.IsNull() {
					t.Error("want null type")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := LoadByString(tt.jsonStr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if node.NType != tt.wantType {
				t.Errorf("type mismatch: want %d, got %d", tt.wantType, node.NType)
			}
			tt.check(t, node)
		})
	}
}

// ========== 数组解析测试 ==========
func TestLoadByString_Arrays(t *testing.T) {
	t.Run("string array", func(t *testing.T) {
		node, err := LoadByString(`["a","b","c"]`)
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != STRINGARRAY {
			t.Errorf("want STRINGARRAY, got %d", node.NType)
		}
		if len(node.ArraysString) != 3 {
			t.Errorf("want len 3, got %d", len(node.ArraysString))
		}
		if node.ArraysString[1] != "b" {
			t.Errorf("want 'b', got '%s'", node.ArraysString[1])
		}
	})

	t.Run("number array", func(t *testing.T) {
		node, err := LoadByString(`[1, 2, 3.5]`)
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != NUMBERARRAY {
			t.Errorf("want NUMBERARRAY, got %d", node.NType)
		}
		if node.ArraysNumber[2] != 3.5 {
			t.Errorf("want 3.5, got %v", node.ArraysNumber[2])
		}
	})

	t.Run("bool array", func(t *testing.T) {
		node, err := LoadByString(`[true, false, true]`)
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != BOOLARRAY {
			t.Errorf("want BOOLARRAY, got %d", node.NType)
		}
		if !node.ArraysBool[0] || node.ArraysBool[1] {
			t.Error("bool array value mismatch")
		}
	})

	t.Run("empty array", func(t *testing.T) {
		node, err := LoadByString(`[]`)
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != STRUCTARRAY {
			t.Errorf("empty array should be STRUCTARRAY, got %d", node.NType)
		}
		if node.Len() != 0 {
			t.Errorf("empty array len should be 0, got %d", node.Len())
		}
	})

	t.Run("mixed array", func(t *testing.T) {
		node, err := LoadByString(`[1, "two", true, null]`)
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != STRUCTARRAY {
			t.Errorf("mixed array should be STRUCTARRAY, got %d", node.NType)
		}
		if node.Len() != 4 {
			t.Errorf("want len 4, got %d", node.Len())
		}
		if node.ArraysStruct[0].ToJsonNode().ValueNumber != 1 {
			t.Error("first element should be 1")
		}
		if node.ArraysStruct[1].ToJsonNode().ValueString != "two" {
			t.Error("second element should be 'two'")
		}
		if !node.ArraysStruct[3].ToJsonNode().IsNull() {
			t.Error("fourth element should be null")
		}
	})

	t.Run("struct array", func(t *testing.T) {
		node, err := LoadByString(`[{"id":1,"name":"a"},{"id":2,"name":"b"}]`)
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != STRUCTARRAY {
			t.Errorf("want STRUCTARRAY, got %d", node.NType)
		}
		first := node.ArraysStruct[0].ToJsonNode()
		if first.GetNodeByName("id").ValueNumber != 1 {
			t.Error("id mismatch")
		}
	})

	t.Run("nested array", func(t *testing.T) {
		node, err := LoadByString(`[[1,2],[3,4]]`)
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != STRUCTARRAY {
			t.Error("nested array should be STRUCTARRAY")
		}
		inner := node.ArraysStruct[0].ToJsonNode()
		if inner.NType != NUMBERARRAY {
			t.Error("inner array should be NUMBERARRAY")
		}
		if inner.ArraysNumber[1] != 2 {
			t.Error("inner value mismatch")
		}
	})
}

// ========== 结构体解析测试 ==========
func TestLoadByString_Struct(t *testing.T) {
	t.Run("flat struct", func(t *testing.T) {
		node, err := LoadByString(`{"name":"alice","age":30,"active":true}`)
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != STRUCT {
			t.Error("root should be STRUCT type")
		}
		if node.GetNodeByName("name").ValueString != "alice" {
			t.Error("name mismatch")
		}
		if node.GetNodeByName("age").ValueNumber != 30 {
			t.Error("age mismatch")
		}
	})

	t.Run("nested struct", func(t *testing.T) {
		jsonStr := `{"a":{"b":{"c":123}}}`
		node, err := LoadByString(jsonStr)
		if err != nil {
			t.Fatal(err)
		}
		target := node.GetNodeByPath("a.b.c")
		if target == nil {
			t.Fatal("path a.b.c not found")
		}
		if target.ValueNumber != 123 {
			t.Errorf("want 123, got %v", target.ValueNumber)
		}
	})

	t.Run("empty struct", func(t *testing.T) {
		node, err := LoadByString(`{}`)
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != STRUCT {
			t.Error("empty object should be STRUCT type")
		}
		if node.Len() != 0 {
			t.Errorf("empty struct len should be 0, got %d", node.Len())
		}
	})
}

// ========== 路径查询测试（含下标） ==========
func TestGetNodeByPath(t *testing.T) {
	jsonStr := `{"level1":{"level2":{"level3":"deep_value"}}}`
	node, _ := LoadByString(jsonStr)

	t.Run("valid deep path", func(t *testing.T) {
		n := node.GetNodeByPath("level1.level2.level3")
		if n == nil || n.ValueString != "deep_value" {
			t.Error("deep path query failed")
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		n := node.GetNodeByPath("level1.not_exist")
		if n != nil {
			t.Error("invalid path should return nil")
		}
	})

	t.Run("empty path", func(t *testing.T) {
		n := node.GetNodeByPath("")
		if n != node {
			t.Error("empty path should return self")
		}
	})
}

func TestGetNodeByPath_Index(t *testing.T) {
	t.Run("struct array index", func(t *testing.T) {
		str := `{"list":[{"id":1,"name":"a"},{"id":2,"name":"b"}]}`
		node, _ := LoadByString(str)
		n := node.GetNodeByPath("list[1].name")
		if n == nil || n.ValueString != "b" {
			t.Errorf("want 'b', got %v", n)
		}
	})

	t.Run("root array index", func(t *testing.T) {
		str := `[{"id":1},{"id":2}]`
		node, _ := LoadByString(str)
		n := node.GetNodeByPath("[0].id")
		if n == nil || n.ValueNumber != 1 {
			t.Error("root array index query failed")
		}
	})

	t.Run("number array index", func(t *testing.T) {
		str := `{"scores":[90,85,95]}`
		node, _ := LoadByString(str)
		n := node.GetNodeByPath("scores[2]")
		if n == nil || n.ValueNumber != 95 {
			t.Error("number array index query failed")
		}
	})

	t.Run("string array index", func(t *testing.T) {
		str := `{"tags":["go","json","lib"]}`
		node, _ := LoadByString(str)
		n := node.GetNodeByPath("tags[0]")
		if n == nil || n.ValueString != "go" {
			t.Error("string array index query failed")
		}
	})

	t.Run("bool array index", func(t *testing.T) {
		str := `{"flags":[true,false,true]}`
		node, _ := LoadByString(str)
		n := node.GetNodeByPath("flags[1]")
		if n == nil || n.ValueBool {
			t.Error("bool array index query failed")
		}
	})

	t.Run("multi-dimensional array index", func(t *testing.T) {
		str := `{"matrix":[[1,2],[3,4]]}`
		node, _ := LoadByString(str)
		n := node.GetNodeByPath("matrix[1][0]")
		if n == nil || n.ValueNumber != 3 {
			t.Error("multi-dimensional index query failed")
		}
	})

	t.Run("index out of bound", func(t *testing.T) {
		str := `{"arr":[1,2,3]}`
		node, _ := LoadByString(str)
		n := node.GetNodeByPath("arr[10]")
		if n != nil {
			t.Error("out of bound should return nil")
		}
	})

	t.Run("invalid index format", func(t *testing.T) {
		str := `{"arr":[1,2,3]}`
		node, _ := LoadByString(str)
		if node.GetNodeByPath("arr[abc]") != nil {
			t.Error("invalid index format should return nil")
		}
		if node.GetNodeByPath("arr[-1]") != nil {
			t.Error("negative index should return nil")
		}
		if node.GetNodeByPath("arr[0") != nil {
			t.Error("unclosed bracket should return nil")
		}
	})

	t.Run("non-array node with index", func(t *testing.T) {
		str := `{"name":"test"}`
		node, _ := LoadByString(str)
		if node.GetNodeByPath("name[0]") != nil {
			t.Error("non-array node with index should return nil")
		}
	})

	t.Run("IsExist with index path", func(t *testing.T) {
		str := `{"list":[{"a":1}]}`
		node, _ := LoadByString(str)
		if !node.IsExist("list[0].a") {
			t.Error("list[0].a should exist")
		}
		if node.IsExist("list[1].a") {
			t.Error("list[1].a should not exist")
		}
	})
}

// ========== 节点增删测试 ==========
func TestAddNode(t *testing.T) {
	root := new(JsonNode)
	root.NType = STRUCT
	root.StructNodes = make(map[string]Inode)

	t.Run("add string node", func(t *testing.T) {
		err := root.AddNode(NowJsonNode("key1", "value1"))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(root.ToString())
		if root.GetNodeByName("key1").ValueString != "value1" {
			t.Error("add string node failed")
		}
	})

	t.Run("add number node", func(t *testing.T) {
		err := root.AddNode(NowJsonNode("key2", 999))
		if err != nil {
			t.Fatal(err)
		}
		t.Log(root.ToString())
		if root.GetNodeByName("key2").ValueNumber != 999 {
			t.Error("add number node failed")
		}
	})

	t.Run("add struct node", func(t *testing.T) {
		child := NowJsonNodeByString("key3", `{"a":1}`)
		err := root.AddNode(child)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(root.ToString())
		if root.GetNodeByPath("key3.a").ValueNumber != 1 {
			t.Error("add struct node failed")
		}
	})

	t.Run("add to non-struct node should fail", func(t *testing.T) {
		numNode := NowJsonNode("num", 123)
		err := numNode.AddNode(NowJsonNode("x", 1))
		t.Log(numNode.NType == NUMBER)
		t.Log(numNode.ToString())
		if err == nil {
			t.Error("should return error when add to non-struct node")
		}
	})
}

func TestDelNode(t *testing.T) {
	root, _ := LoadByString(`{"a":1,"b":2,"c":3}`)

	t.Run("delete exist node", func(t *testing.T) {
		err := root.DelNode("b")
		if err != nil {
			t.Fatal(err)
		}
		if root.IsExist("b") {
			t.Error("node 'b' should be deleted")
		}
		if root.Len() != 2 {
			t.Errorf("want len 2, got %d", root.Len())
		}
	})

	t.Run("delete non-exist node no error", func(t *testing.T) {
		err := root.DelNode("not_exist")
		if err != nil {
			t.Error("delete non-exist node should not return error")
		}
	})

	t.Run("delete on non-struct node should fail", func(t *testing.T) {
		numNode := NowJsonNode("num", 123)
		err := numNode.DelNode("x")
		if err == nil {
			t.Error("should return error when delete on non-struct node")
		}
	})
}

// ========== 原生数组SetValue测试 ==========
func TestSetValue_NativeArrays(t *testing.T) {
	t.Run("[]int to NUMBERARRAY", func(t *testing.T) {
		node := new(JsonNode)
		err := node.SetValue([]int{1, 2, 3})
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != NUMBERARRAY {
			t.Errorf("want NUMBERARRAY, got %d", node.NType)
		}
		if node.ArraysNumber[2] != 3 {
			t.Errorf("want 3, got %v", node.ArraysNumber[2])
		}
	})

	t.Run("[]string to STRINGARRAY", func(t *testing.T) {
		node := new(JsonNode)
		err := node.SetValue([]string{"x", "y"})
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != STRINGARRAY {
			t.Errorf("want STRINGARRAY, got %d", node.NType)
		}
		if node.ArraysString[0] != "x" {
			t.Error("string array value mismatch")
		}
	})

	t.Run("[]bool to BOOLARRAY", func(t *testing.T) {
		node := new(JsonNode)
		err := node.SetValue([]bool{true, false})
		if err != nil {
			t.Fatal(err)
		}
		if node.NType != BOOLARRAY {
			t.Errorf("want BOOLARRAY, got %d", node.NType)
		}
		if !node.ArraysBool[0] {
			t.Error("bool array value mismatch")
		}
	})
}

// ========== 序列化往返测试 ==========
func TestToString_RoundTrip(t *testing.T) {
	cases := []string{
		`{"name":"test","age":20,"active":true}`,
		`{"arr":[1,2,3],"nested":{"a":"b"}}`,
		`[1,"two",false,null]`,
		`{"empty_obj":{},"empty_arr":[]}`,
	}

	for i, original := range cases {
		t.Run(string(rune('0'+i)), func(t *testing.T) {
			node1, err := LoadByString(original)
			if err != nil {
				t.Fatalf("first parse error: %v", err)
			}
			str := node1.ToString()

			var tmp interface{}
			if err := json.Unmarshal([]byte(str), &tmp); err != nil {
				t.Fatalf("generated json is invalid: %v\njson: %s", err, str)
			}

			node2, err := LoadByString(str)
			if err != nil {
				t.Fatalf("second parse error: %v", err)
			}
			if node2.Len() != node1.Len() {
				t.Error("round trip length mismatch")
			}
		})
	}
}

func TestToString_Escape(t *testing.T) {
	t.Run("special characters", func(t *testing.T) {
		node := NowJsonNode("msg", `he said "hello" \ world`)
		str := node.ToString()

		var result map[string]string
		if err := json.Unmarshal([]byte("{"+str+"}"), &result); err != nil {
			t.Fatalf("escaped json parse failed: %v\nstr: %s", err, str)
		}
		if result["msg"] != `he said "hello" \ world` {
			t.Errorf("round trip string mismatch: %s", result["msg"])
		}
	})

	t.Run("control characters", func(t *testing.T) {
		node := NowJsonNode("text", "line1\nline2\tend")
		str := node.ToString()

		var result map[string]string
		if err := json.Unmarshal([]byte("{"+str+"}"), &result); err != nil {
			t.Fatalf("control char json parse failed: %v", err)
		}
		if result["text"] != "line1\nline2\tend" {
			t.Error("control characters round trip failed")
		}
	})
}

// ========== 辅助方法测试 ==========
func TestIsExist(t *testing.T) {
	node, _ := LoadByString(`{"a":{"b":1},"c":2}`)

	if !node.IsExist("a") {
		t.Error("a should exist")
	}
	if !node.IsExist("a.b") {
		t.Error("a.b should exist")
	}
	if node.IsExist("a.x") {
		t.Error("a.x should not exist")
	}
	if node.IsExist("d") {
		t.Error("d should not exist")
	}
}

func TestGetChildNames(t *testing.T) {
	node, _ := LoadByString(`{"a":1,"b":2,"c":3}`)
	names := node.GetChildNames()
	if len(names) != 3 {
		t.Errorf("want 3 names, got %d", len(names))
	}
	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}
	if !nameSet["a"] || !nameSet["b"] || !nameSet["c"] {
		t.Error("child names incomplete")
	}
}

// ========== 错误场景测试 ==========
func TestErrorCases(t *testing.T) {
	t.Run("invalid json string", func(t *testing.T) {
		_, err := LoadByString(`{invalid json}`)
		if err == nil {
			t.Error("invalid json should return error")
		}
	})

	t.Run("NowJsonNodeByString with bad json", func(t *testing.T) {
		node := NowJsonNodeByString("bad", `{wrong}`)
		if node == nil {
			t.Fatal("should not return nil")
		}
		if !node.IsNull() {
			t.Error("bad json should result in NULL type node")
		}
	})
}

// ========== 边界测试 ==========
func TestEdgeCases(t *testing.T) {
	t.Run("deep nested struct", func(t *testing.T) {
		deep := `{"a":{"b":{"c":{"d":{"e":{"f":1}}}}}}`
		node, err := LoadByString(deep)
		if err != nil {
			t.Fatal(err)
		}
		if node.GetNodeByPath("a.b.c.d.e.f").ValueNumber != 1 {
			t.Error("deep nested query failed")
		}
	})

	t.Run("large number array", func(t *testing.T) {
		node := new(JsonNode)
		arr := make([]int, 1000)
		for i := 0; i < 1000; i++ {
			arr[i] = i
		}
		_ = node.SetValue(arr)
		if node.Len() != 1000 {
			t.Errorf("want len 1000, got %d", node.Len())
		}
		if node.ArraysNumber[999] != 999 {
			t.Error("last element mismatch")
		}
	})

	t.Run("empty string value", func(t *testing.T) {
		node, err := LoadByString(`""`)
		if err != nil {
			t.Fatal(err)
		}
		if node.ValueString != "" {
			t.Error("empty string value mismatch")
		}
	})
}

// ========== 修复验证测试 ==========
func TestBugfix_AddNodeNil(t *testing.T) {
	root := new(JsonNode)
	root.NType = STRUCT
	root.StructNodes = make(map[string]Inode)
	// 传入 nil 节点不 panic
	err := root.AddNode(nil)
	if err != nil {
		t.Error("AddChild(nil) should return error:", err.Error())
	}
}

func TestBugfix_ClearArrayNil(t *testing.T) {
	var node *JsonNode = nil
	// nil 接收者调用不 panic
	err := node.ClearArray()
	if err != nil {
		t.Error("ClearArray on nil node should return error:")
	}
}

func TestBugfix_ArrayValueCopy(t *testing.T) {
	root, _ := LoadByString(`{"scores":[1,2,3]}`)
	// 直接修改返回的副本不应生效
	node := root.GetNodeByPath("scores[0]")
	node.ValueNumber = 999
	t.Log(root.ToString())
	if root.GetNodeByPath("scores[0]").ValueNumber == 999 {
		t.Error("副本修改不应影响原数组")
	}
	// SetIndex 应当生效
	_ = root.GetNodeByName("scores").SetIndex(0, 999)
	t.Log(root.ToString())
	if root.GetNodeByPath("scores[0]").ValueNumber != 999 {
		t.Error("SetIndex should modify original array")
	}
}

func TestBugfix_FloatFormat(t *testing.T) {
	node := NowJsonNode("num", 123456789.0)
	str := node.ToString()
	t.Log(str)
	// 不应出现科学计数法
	if strings.Contains(str, "e") || strings.Contains(str, "E") {
		t.Errorf("数字不应输出科学计数法: %s", str)
	}
}

func TestBugfix_ByteArrayAppend(t *testing.T) {
	node := NowJsonNode("tags", []string{"go"})
	err := node.Append([]byte("json"))
	if err != nil {
		t.Fatal("[]byte 追加到字符串数组失败:", err)
	}
	if node.ArraysString[1] != "json" {
		t.Error("[]byte 转换失败")
	}
	t.Log(node.ToString())
}

// ======================================
// ========== 基准测试 Benchmark =========
// ======================================

const smallJSON = `{"name":"test","age":25,"active":true,"tags":["go","json"]}`
const largeJSON = `{"data":[{"id":1,"name":"user1","score":99.5},{"id":2,"name":"user2","score":88.0},{"id":3,"name":"user3","score":76.3}],"meta":{"total":3,"page":1,"size":10},"status":"ok"}`

// BenchmarkLoadByString_Small 小JSON解析性能
func BenchmarkLoadByString_Small(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = LoadByString(smallJSON)
	}
}

// BenchmarkLoadByString_Large 大JSON解析性能
func BenchmarkLoadByString_Large(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = LoadByString(largeJSON)
	}
}

// BenchmarkGetNodeByPath_Shallow 浅层路径查询性能
func BenchmarkGetNodeByPath_Shallow(b *testing.B) {
	node, _ := LoadByString(largeJSON)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.GetNodeByPath("status")
	}
}

// BenchmarkGetNodeByPath_Deep 深层路径查询性能
func BenchmarkGetNodeByPath_Deep(b *testing.B) {
	node, _ := LoadByString(largeJSON)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.GetNodeByPath("meta.total")
	}
}

// BenchmarkGetNodeByPath_WithIndex 带下标路径查询性能
func BenchmarkGetNodeByPath_WithIndex(b *testing.B) {
	node, _ := LoadByString(largeJSON)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.GetNodeByPath("data[1].name")
	}
}

// BenchmarkToString_Small 小对象序列化性能
func BenchmarkToString_Small(b *testing.B) {
	node, _ := LoadByString(smallJSON)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.ToString()
	}
}

// BenchmarkToString_Large 大对象序列化性能
func BenchmarkToString_Large(b *testing.B) {
	node, _ := LoadByString(largeJSON)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.ToString()
	}
}

// BenchmarkAddNode 节点添加性能
func BenchmarkAddNode(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := new(JsonNode)
		root.NType = STRUCT
		root.StructNodes = make(map[string]Inode)
		for j := 0; j < 100; j++ {
			_ = root.AddNode(NowJsonNode("key"+string(rune('0'+j%10)), j))
		}
	}
}

// BenchmarkSetValue_IntArray 整数数组转换性能
func BenchmarkSetValue_IntArray(b *testing.B) {
	arr := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		arr[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node := new(JsonNode)
		_ = node.SetValue(arr)
	}
}

// BenchmarkNumbers2Floats 数字切片转换底层性能
func BenchmarkNumbers2Floats(b *testing.B) {
	arr := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		arr[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = numbers2floats(arr)
	}
}

// BenchmarkRoundTrip 解析+序列化往返性能
func BenchmarkRoundTrip(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node, _ := LoadByString(largeJSON)
		_ = node.ToString()
	}
}
