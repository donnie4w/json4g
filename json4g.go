// Copyright (c) 2016, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// donnie4w@gmail.com

package json4g

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Inode interface {
	ToJsonNode() (node *JsonNode)
}

type NODETYPE int8

const (
	STRUCT      NODETYPE = iota // {"a":1,"b":"2"}
	STRUCTARRAY                 // [{"a":1,"b":"2"},{"aa":1,"bb":"2"},1,"2"]
	NUMBERARRAY                 // [1,2,3]
	BOOLARRAY                   // [false,false,true]
	STRINGARRAY                 // ["a","b"]
	NUMBER                      // "nodename":1
	BOOL                        // "nodename":false
	STRING                      // "nodename":"abcd"
	NULL                        // null 值
)

type JsonNode struct {
	Name         string
	ValueNumber  float64
	ValueString  string
	ValueBool    bool
	ArraysNumber []float64
	ArraysString []string
	ArraysBool   []bool
	ArraysStruct []Inode
	StructNodes  map[string]Inode
	NType        NODETYPE
}

// NowJsonNodeByString 从JSON字符串创建命名节点
func NowJsonNodeByString(nodename, jsonstr string) (json *JsonNode) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println("NowJsonNodeByString error ", er)
		}
	}()
	var err error
	json, err = LoadByString(jsonstr)
	if err != nil || json == nil {
		json = new(JsonNode)
		json.NType = NULL
	}
	json.Name = nodename
	return
}

// NowJsonNode 创建一个命名节点，自动识别值类型
func NowJsonNode(nodename string, nodevalue interface{}) *JsonNode {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println("NowJsonNode error", er)
		}
	}()
	json := new(JsonNode)
	json.Name = nodename
	if err := json.SetValue(nodevalue); err != nil {
		json.NType = NULL
	}
	return json
}

// LoadByString 从JSON字符串加载根节点
func LoadByString(jsonstr string) (root *JsonNode, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("LoadByString error ", er))
		}
	}()
	var dat interface{}
	if err = json.Unmarshal([]byte(jsonstr), &dat); err == nil {
		root = new(JsonNode)
		root.SetValue(dat)
	}
	return
}

// AddNode 增加子节点，仅STRUCT类型有效
func (n *JsonNode) AddNode(node *JsonNode) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("AddChild error ", er))
		}
	}()
	if n == nil {
		return errors.New("AddChild: node is nil")
	}
	if n.NType != STRUCT {
		return errors.New("AddChild: current node is not STRUCT type")
	}
	if n.StructNodes == nil {
		n.StructNodes = make(map[string]Inode)
	}
	n.StructNodes[node.Name] = node
	return nil
}

// DelNode 删除子节点，仅STRUCT类型有效
func (n *JsonNode) DelNode(nodename string) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("DelNode error ", er))
		}
	}()
	if n.NType != STRUCT {
		return errors.New("DelNode: current node is not STRUCT type")
	}
	delete(n.StructNodes, nodename)
	return
}

func (n *JsonNode) ToJsonNode() *JsonNode {
	return n
}

// SetValue 设置节点值，自动识别类型
func (n *JsonNode) SetValue(value interface{}) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("set value error ", er))
		}
	}()
	switch v := value.(type) {
	case []byte:
		n.ValueString = string(v)
		n.NType = STRING
	case int:
		n.setJsonValue(float64(v))
	case int8:
		n.setJsonValue(float64(v))
	case int16:
		n.setJsonValue(float64(v))
	case int32:
		n.setJsonValue(float64(v))
	case int64:
		n.setJsonValue(float64(v))
	case uint:
		n.setJsonValue(float64(v))
	case uint8:
		n.setJsonValue(float64(v))
	case uint16:
		n.setJsonValue(float64(v))
	case uint32:
		n.setJsonValue(float64(v))
	case uint64:
		n.setJsonValue(float64(v))
	case float32:
		n.setJsonValue(float64(v))
	case []int:
		n.setJsonValue(numbers2floats(v))
	case []int8:
		n.setJsonValue(numbers2floats(v))
	case []int16:
		n.setJsonValue(numbers2floats(v))
	case []int32:
		n.setJsonValue(numbers2floats(v))
	case []int64:
		n.setJsonValue(numbers2floats(v))
	case []uint:
		n.setJsonValue(numbers2floats(v))
	case []uint16:
		n.setJsonValue(numbers2floats(v))
	case []uint64:
		n.setJsonValue(numbers2floats(v))
	case []float32:
		n.setJsonValue(numbers2floats(v))
	default:
		n.setJsonValue(value)
	}
	return
}

// setJsonValue 核心值设置逻辑
func (n *JsonNode) setJsonValue(value interface{}) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("set error ", er))
		}
	}()
	switch v := value.(type) {
	case nil:
		n.NType = NULL
	case string:
		n.ValueString = v
		n.NType = STRING
	case bool:
		n.ValueBool = v
		n.NType = BOOL
	case float64:
		n.ValueNumber = v
		n.NType = NUMBER
	case []float64:
		n.ArraysNumber = v
		n.NType = NUMBERARRAY
	case []string:
		n.ArraysString = v
		n.NType = STRINGARRAY
	case []bool:
		n.ArraysBool = v
		n.NType = BOOLARRAY
	case map[string]interface{}:
		n.NType = STRUCT
		parseStruct(v, n)
	case []interface{}:
		if len(v) == 0 {
			n.ArraysStruct = []Inode{}
			n.NType = STRUCTARRAY
			return
		}
		switch v[0].(type) {
		case string:
			if res, er := interfaces2strings(v); er == nil {
				n.ArraysString = res
				n.NType = STRINGARRAY
			} else {
				n.ArraysStruct = interfaces2array(v)
				n.NType = STRUCTARRAY
			}
		case float64:
			if res, er := interfaces2floats(v); er == nil {
				n.ArraysNumber = res
				n.NType = NUMBERARRAY
			} else {
				n.ArraysStruct = interfaces2array(v)
				n.NType = STRUCTARRAY
			}
		case bool:
			if res, er := interfaces2bools(v); er == nil {
				n.ArraysBool = res
				n.NType = BOOLARRAY
			} else {
				n.ArraysStruct = interfaces2array(v)
				n.NType = STRUCTARRAY
			}
		default:
			n.ArraysStruct = interfaces2array(v)
			n.NType = STRUCTARRAY
		}
	}
	return
}

// ==================== 数组操作 ====================

// SetIndex 修改数组指定索引位置的元素
// 仅数组类型节点有效，索引越界或类型不匹配返回错误
// 注意：基本类型数组请使用本方法修改，直接修改 GetNodeByPath 返回的节点不生效
func (n *JsonNode) SetIndex(index int, value interface{}) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("SetIndex error ", er))
		}
	}()
	if index < 0 {
		return errors.New("SetIndex: index cannot be negative")
	}
	switch n.NType {
	case NUMBERARRAY:
		if index >= len(n.ArraysNumber) {
			return errors.New("SetIndex: index out of range")
		}
		f, ok := toFloat64(value)
		if !ok {
			return errors.New("SetIndex: value type mismatch for NUMBERARRAY")
		}
		n.ArraysNumber[index] = f

	case STRINGARRAY:
		if index >= len(n.ArraysString) {
			return errors.New("SetIndex: index out of range")
		}
		switch v := value.(type) {
		case string:
			n.ArraysString[index] = v
		case []byte:
			n.ArraysString[index] = string(v)
		default:
			return errors.New("SetIndex: value type mismatch for STRINGARRAY")
		}

	case BOOLARRAY:
		if index >= len(n.ArraysBool) {
			return errors.New("SetIndex: index out of range")
		}
		b, ok := value.(bool)
		if !ok {
			return errors.New("SetIndex: value type mismatch for BOOLARRAY")
		}
		n.ArraysBool[index] = b

	case STRUCTARRAY:
		if index >= len(n.ArraysStruct) {
			return errors.New("SetIndex: index out of range")
		}
		switch v := value.(type) {
		case *JsonNode:
			n.ArraysStruct[index] = v
		case Inode:
			n.ArraysStruct[index] = v
		default:
			node := new(JsonNode)
			if e := node.SetValue(v); e != nil {
				return e
			}
			n.ArraysStruct[index] = node
		}

	default:
		return errors.New("SetIndex: current node is not array type")
	}
	return nil
}

// Append 向数组末尾追加元素
// 仅数组类型节点有效，类型不匹配返回错误
func (n *JsonNode) Append(value interface{}) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("Append error ", er))
		}
	}()
	switch n.NType {
	case NUMBERARRAY:
		f, ok := toFloat64(value)
		if !ok {
			return errors.New("Append: value type mismatch for NUMBERARRAY")
		}
		n.ArraysNumber = append(n.ArraysNumber, f)

	case STRINGARRAY:
		switch v := value.(type) {
		case string:
			n.ArraysString = append(n.ArraysString, v)
		case []byte:
			n.ArraysString = append(n.ArraysString, string(v))
		default:
			return errors.New("Append: value type mismatch for STRINGARRAY")
		}

	case BOOLARRAY:
		b, ok := value.(bool)
		if !ok {
			return errors.New("Append: value type mismatch for BOOLARRAY")
		}
		n.ArraysBool = append(n.ArraysBool, b)

	case STRUCTARRAY:
		switch v := value.(type) {
		case *JsonNode:
			n.ArraysStruct = append(n.ArraysStruct, v)
		case Inode:
			n.ArraysStruct = append(n.ArraysStruct, v)
		default:
			node := new(JsonNode)
			if e := node.SetValue(v); e != nil {
				return e
			}
			n.ArraysStruct = append(n.ArraysStruct, node)
		}

	default:
		return errors.New("Append: current node is not array type")
	}
	return nil
}

// Insert 在数组指定索引位置插入元素
// 仅数组类型节点有效，索引越界或类型不匹配返回错误
func (n *JsonNode) Insert(index int, value interface{}) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("Insert error ", er))
		}
	}()
	if index < 0 {
		return errors.New("Insert: index cannot be negative")
	}
	switch n.NType {
	case NUMBERARRAY:
		if index > len(n.ArraysNumber) {
			return errors.New("Insert: index out of range")
		}
		f, ok := toFloat64(value)
		if !ok {
			return errors.New("Insert: value type mismatch for NUMBERARRAY")
		}
		n.ArraysNumber = append(n.ArraysNumber[:index], append([]float64{f}, n.ArraysNumber[index:]...)...)

	case STRINGARRAY:
		if index > len(n.ArraysString) {
			return errors.New("Insert: index out of range")
		}
		var s string
		switch v := value.(type) {
		case string:
			s = v
		case []byte:
			s = string(v)
		default:
			return errors.New("Insert: value type mismatch for STRINGARRAY")
		}
		n.ArraysString = append(n.ArraysString[:index], append([]string{s}, n.ArraysString[index:]...)...)

	case BOOLARRAY:
		if index > len(n.ArraysBool) {
			return errors.New("Insert: index out of range")
		}
		b, ok := value.(bool)
		if !ok {
			return errors.New("Insert: value type mismatch for BOOLARRAY")
		}
		n.ArraysBool = append(n.ArraysBool[:index], append([]bool{b}, n.ArraysBool[index:]...)...)

	case STRUCTARRAY:
		if index > len(n.ArraysStruct) {
			return errors.New("Insert: index out of range")
		}
		var node Inode
		switch v := value.(type) {
		case *JsonNode:
			node = v
		case Inode:
			node = v
		default:
			jsonNode := new(JsonNode)
			if e := jsonNode.SetValue(v); e != nil {
				return e
			}
			node = jsonNode
		}
		n.ArraysStruct = append(n.ArraysStruct[:index], append([]Inode{node}, n.ArraysStruct[index:]...)...)

	default:
		return errors.New("Insert: current node is not array type")
	}
	return nil
}

// RemoveIndex 删除数组指定索引位置的元素
// 仅数组类型节点有效，索引越界返回错误
func (n *JsonNode) RemoveIndex(index int) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("RemoveIndex error ", er))
		}
	}()
	if index < 0 {
		return errors.New("RemoveIndex: index cannot be negative")
	}
	switch n.NType {
	case NUMBERARRAY:
		if index >= len(n.ArraysNumber) {
			return errors.New("RemoveIndex: index out of range")
		}
		n.ArraysNumber = append(n.ArraysNumber[:index], n.ArraysNumber[index+1:]...)

	case STRINGARRAY:
		if index >= len(n.ArraysString) {
			return errors.New("RemoveIndex: index out of range")
		}
		n.ArraysString = append(n.ArraysString[:index], n.ArraysString[index+1:]...)

	case BOOLARRAY:
		if index >= len(n.ArraysBool) {
			return errors.New("RemoveIndex: index out of range")
		}
		n.ArraysBool = append(n.ArraysBool[:index], n.ArraysBool[index+1:]...)

	case STRUCTARRAY:
		if index >= len(n.ArraysStruct) {
			return errors.New("RemoveIndex: index out of range")
		}
		n.ArraysStruct = append(n.ArraysStruct[:index], n.ArraysStruct[index+1:]...)

	default:
		return errors.New("RemoveIndex: current node is not array type")
	}
	return nil
}

// ClearArray 清空数组
// 仅数组类型节点有效
func (n *JsonNode) ClearArray() (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("ClearArray error ", er))
		}
	}()
	switch n.NType {
	case NUMBERARRAY:
		n.ArraysNumber = n.ArraysNumber[:0]
	case STRINGARRAY:
		n.ArraysString = n.ArraysString[:0]
	case BOOLARRAY:
		n.ArraysBool = n.ArraysBool[:0]
	case STRUCTARRAY:
		n.ArraysStruct = n.ArraysStruct[:0]
	default:
		return errors.New("ClearArray: current node is not array type")
	}
	return nil
}

// toFloat64 数值类型统一转换
func toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	default:
		return 0, false
	}
}

func numbers2floats(inters interface{}) (floats []float64) {
	switch v := inters.(type) {
	case []int:
		floats = make([]float64, len(v))
		for i, val := range v {
			floats[i] = float64(val)
		}
	case []int8:
		floats = make([]float64, len(v))
		for i, val := range v {
			floats[i] = float64(val)
		}
	case []int16:
		floats = make([]float64, len(v))
		for i, val := range v {
			floats[i] = float64(val)
		}
	case []int32:
		floats = make([]float64, len(v))
		for i, val := range v {
			floats[i] = float64(val)
		}
	case []int64:
		floats = make([]float64, len(v))
		for i, val := range v {
			floats[i] = float64(val)
		}
	case []uint:
		floats = make([]float64, len(v))
		for i, val := range v {
			floats[i] = float64(val)
		}
	case []uint16:
		floats = make([]float64, len(v))
		for i, val := range v {
			floats[i] = float64(val)
		}
	case []uint64:
		floats = make([]float64, len(v))
		for i, val := range v {
			floats[i] = float64(val)
		}
	case []float32:
		floats = make([]float64, len(v))
		for i, val := range v {
			floats[i] = float64(val)
		}
	}
	return
}

func interfaces2floats(inters []interface{}) (floats []float64, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("interfaces2floats error ", er))
			floats = nil
		}
	}()
	floats = make([]float64, len(inters))
	for i, v := range inters {
		floats[i] = v.(float64)
	}
	return
}

func interfaces2strings(inters []interface{}) (strings []string, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("interfaces2strings error ", er))
			strings = nil
		}
	}()
	strings = make([]string, len(inters))
	for i, v := range inters {
		strings[i] = v.(string)
	}
	return
}

func interfaces2bools(inters []interface{}) (bools []bool, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("interfaces2bools error ", er))
			bools = nil
		}
	}()
	bools = make([]bool, len(inters))
	for i, v := range inters {
		bools[i] = v.(bool)
	}
	return
}

func interfaces2struct(inters []interface{}) (n []Inode, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("interfaces2struct error ", er))
			n = nil
		}
	}()
	n = make([]Inode, len(inters))
	for i, v := range inters {
		jsonnode := new(JsonNode)
		parseStruct(v.(map[string]interface{}), jsonnode)
		n[i] = jsonnode
	}
	return
}

func interfaces2array(inters []interface{}) (n []Inode) {
	n = make([]Inode, len(inters))
	for i, v := range inters {
		jsonnode := new(JsonNode)
		jsonnode.SetValue(v)
		n[i] = jsonnode
	}
	return
}

// writeEscapedString JSON字符串转义写入buffer
func writeEscapedString(buf *bytes.Buffer, s string) {
	for _, r := range s {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		default:
			if r < 0x20 {
				_, _ = fmt.Fprintf(buf, `\u%04x`, r)
			} else {
				buf.WriteRune(r)
			}
		}
	}
}

// writeFloat 格式化数字输出，避免不必要的科学计数法
func writeFloat(buf *bytes.Buffer, f float64) {
	buf.WriteString(strconv.FormatFloat(f, 'f', -1, 64))
}

// ToString 序列化为JSON字符串
func (n *JsonNode) ToString() string {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	writeJsonNode(buf, n)
	return buf.String()
}

func writeJsonNode(buf *bytes.Buffer, n *JsonNode) {
	switch n.NType {
	case STRUCT:
		if n.Name != "" {
			buf.WriteByte('"')
			writeEscapedString(buf, n.Name)
			buf.WriteString(`":{`)
		} else {
			buf.WriteByte('{')
		}
		first := true
		for _, dat := range n.StructNodes {
			if !first {
				buf.WriteByte(',')
			}
			first = false
			writeJsonNode(buf, dat.ToJsonNode())
		}
		buf.WriteByte('}')

	case NUMBERARRAY:
		if n.Name != "" {
			buf.WriteByte('"')
			writeEscapedString(buf, n.Name)
			buf.WriteString(`":[`)
		} else {
			buf.WriteByte('[')
		}
		for i, v := range n.ArraysNumber {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeFloat(buf, v)
		}
		buf.WriteByte(']')

	case STRING:
		if n.Name != "" {
			buf.WriteByte('"')
			writeEscapedString(buf, n.Name)
			buf.WriteString(`":"`)
		} else {
			buf.WriteByte('"')
		}
		writeEscapedString(buf, n.ValueString)
		buf.WriteByte('"')

	case NUMBER:
		if n.Name != "" {
			buf.WriteByte('"')
			writeEscapedString(buf, n.Name)
			buf.WriteString(`":`)
		}
		writeFloat(buf, n.ValueNumber)

	case BOOL:
		if n.Name != "" {
			buf.WriteByte('"')
			writeEscapedString(buf, n.Name)
			buf.WriteString(`":`)
		}
		if n.ValueBool {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}

	case NULL:
		if n.Name != "" {
			buf.WriteByte('"')
			writeEscapedString(buf, n.Name)
			buf.WriteString(`":null`)
		} else {
			buf.WriteString("null")
		}

	case STRUCTARRAY:
		if n.Name != "" {
			buf.WriteByte('"')
			writeEscapedString(buf, n.Name)
			buf.WriteString(`":[`)
		} else {
			buf.WriteByte('[')
		}
		for i, v := range n.ArraysStruct {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeJsonNode(buf, v.ToJsonNode())
		}
		buf.WriteByte(']')

	case STRINGARRAY:
		if n.Name != "" {
			buf.WriteByte('"')
			writeEscapedString(buf, n.Name)
			buf.WriteString(`":[`)
		} else {
			buf.WriteByte('[')
		}
		for i, v := range n.ArraysString {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte('"')
			writeEscapedString(buf, v)
			buf.WriteByte('"')
		}
		buf.WriteByte(']')

	case BOOLARRAY:
		if n.Name != "" {
			buf.WriteByte('"')
			writeEscapedString(buf, n.Name)
			buf.WriteString(`":[`)
		} else {
			buf.WriteByte('[')
		}
		for i, v := range n.ArraysBool {
			if i > 0 {
				buf.WriteByte(',')
			}
			if v {
				buf.WriteString("true")
			} else {
				buf.WriteString("false")
			}
		}
		buf.WriteByte(']')
	}
}

// GetChildNames 获取直接子节点名称列表
func (n *JsonNode) GetChildNames() (names []string) {
	if n.StructNodes != nil && len(n.StructNodes) > 0 {
		names = make([]string, 0, len(n.StructNodes))
		for k := range n.StructNodes {
			names = append(names, k)
		}
	}
	return
}

// IsExist 判断路径是否存在，支持下标路径
func (n *JsonNode) IsExist(path string) bool {
	return n.GetNodeByPath(path) != nil
}

func isTypeStruct(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	return ok
}

func parseStruct(dat map[string]interface{}, superNode *JsonNode) {
	superNode.StructNodes = make(map[string]Inode, len(dat))
	nodes := superNode.StructNodes
	for k, v := range dat {
		json := new(JsonNode)
		json.Name = k
		json.setJsonValue(v)
		nodes[k] = json
	}
}

// parsePathSegment 解析路径段，提取字段名和数组下标列表
func parsePathSegment(seg string) (name string, indexes []int, ok bool) {
	seg = strings.TrimSpace(seg)
	if seg == "" {
		return "", nil, false
	}
	bracketIdx := strings.Index(seg, "[")
	if bracketIdx == -1 {
		return seg, nil, true
	}
	name = seg[:bracketIdx]
	rest := seg[bracketIdx:]

	for len(rest) > 0 {
		if rest[0] != '[' {
			return "", nil, false
		}
		closeIdx := strings.Index(rest, "]")
		if closeIdx == -1 {
			return "", nil, false
		}
		numStr := rest[1:closeIdx]
		num, err := strconv.Atoi(numStr)
		if err != nil || num < 0 {
			return "", nil, false
		}
		indexes = append(indexes, num)
		rest = rest[closeIdx+1:]
	}
	return name, indexes, true
}

// getArrayElement 根据索引获取数组元素，包装为JsonNode返回
// 注意：对于 NUMBER/STRING/BOOL 数组，返回的是值副本，修改不影响原数组；修改请使用 SetIndex
func getArrayElement(n *JsonNode, index int) *JsonNode {
	if n == nil {
		return nil
	}
	switch n.NType {
	case STRUCTARRAY:
		if index < 0 || index >= len(n.ArraysStruct) {
			return nil
		}
		return n.ArraysStruct[index].ToJsonNode()
	case NUMBERARRAY:
		if index < 0 || index >= len(n.ArraysNumber) {
			return nil
		}
		return &JsonNode{
			ValueNumber: n.ArraysNumber[index],
			NType:       NUMBER,
		}
	case STRINGARRAY:
		if index < 0 || index >= len(n.ArraysString) {
			return nil
		}
		return &JsonNode{
			ValueString: n.ArraysString[index],
			NType:       STRING,
		}
	case BOOLARRAY:
		if index < 0 || index >= len(n.ArraysBool) {
			return nil
		}
		return &JsonNode{
			ValueBool: n.ArraysBool[index],
			NType:     BOOL,
		}
	default:
		return nil
	}
}

// GetNodeByPath 根据路径查询子节点，支持点分路径和数组下标
// 示例：a.b.c、a.b[0].d、[0].name、matrix[0][1]
// 注意：基本类型数组下标返回值副本，修改请使用 SetIndex 方法
func (n *JsonNode) GetNodeByPath(path string) (node *JsonNode) {
	defer func() {
		if er := recover(); er != nil {
			node = nil
		}
	}()
	if path == "" {
		return n
	}
	paths := strings.Split(path, ".")
	node = n
	for _, p := range paths {
		name, indexes, ok := parsePathSegment(p)
		if !ok {
			return nil
		}
		if name != "" {
			node = getChildJsonNode(node, name)
			if node == nil {
				return nil
			}
		}
		for _, idx := range indexes {
			node = getArrayElement(node, idx)
			if node == nil {
				return nil
			}
		}
	}
	return
}

// GetNodeByName 根据节点名查询直接子节点
func (n *JsonNode) GetNodeByName(name string) *JsonNode {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println(er)
		}
	}()
	if n.NType == STRUCT && n.StructNodes != nil {
		if node, ok := n.StructNodes[name]; ok {
			return node.(*JsonNode)
		}
	}
	return nil
}

func getChildJsonNode(n *JsonNode, name string) *JsonNode {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println(er)
		}
	}()
	if n.StructNodes != nil {
		if node, ok := n.StructNodes[name]; ok {
			return node.(*JsonNode)
		}
	}
	return nil
}

// Len 返回子节点数量或数组长度
func (n *JsonNode) Len() int {
	switch n.NType {
	case STRUCT:
		return len(n.StructNodes)
	case STRUCTARRAY:
		return len(n.ArraysStruct)
	case NUMBERARRAY:
		return len(n.ArraysNumber)
	case STRINGARRAY:
		return len(n.ArraysString)
	case BOOLARRAY:
		return len(n.ArraysBool)
	default:
		return 0
	}
}

// Type 返回节点类型
func (n *JsonNode) Type() NODETYPE {
	return n.NType
}

// IsNull 判断是否为null值
func (n *JsonNode) IsNull() bool {
	return n.NType == NULL
}
