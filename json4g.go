//   donnie4w@gmail.com
//   1.0.0

package json4g

import (
	"encoding/json"
	"errors"
	"fmt"
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

func NowJsonNodeByString(nodename, jsonstr string) (json *JsonNode) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println("NowJsonNodeByString error ", er)
		}
	}()
	json, _ = LoadByString(jsonstr)
	json.Name = nodename
	return
}

func NowJsonNode(nodename string, nodevalue interface{}) (json *JsonNode) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println("NowJsonNode error ", er)
		}
	}()
	json = new(JsonNode)
	json.Name = nodename
	json.setJsonValue(nodevalue)
	return
}

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

//增加子节点
func (n *JsonNode) AddNode(node *JsonNode) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("AddNode error ", er))
		}
	}()
	if n.NType == STRUCT {
		if n.StructNodes == nil {
			n.StructNodes = make(map[string]Inode, 0)
		}
		n.StructNodes[node.Name] = node
	}
	return
}

//删除子节点
func (n *JsonNode) DelNode(nodename string) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("DelNode error ", er))
		}
	}()
	if n.NType == STRUCT {
		delete(n.StructNodes, nodename)
	}
	return
}

func (n *JsonNode) ToJsonNode() (node *JsonNode) {
	node = n
	return
}

func (n *JsonNode) SetValue(value interface{}) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("set value error ", er))
		}
	}()
	switch value.(type) {
	case []byte:
		n.ValueString = string(value.([]byte))
		n.NType = STRING
	case int:
		n.setJsonValue(float64(value.(int)))
	case int32:
		n.setJsonValue(float64(value.(int32)))
	case int64:
		n.setJsonValue(float64(value.(int64)))
	case float32:
		n.setJsonValue(float64(value.(float32)))
	case uint8:
		n.setJsonValue(float64(value.(uint8)))
	case int8:
		n.setJsonValue(float64(value.(int8)))
	case int16:
		n.setJsonValue(float64(value.(int16)))
	case uint16:
		n.setJsonValue(float64(value.(uint16)))
	case uint64:
		n.setJsonValue(float64(value.(uint64)))
	case []int:
		n.setJsonValue(numbers2floats(value.([]int)))
	case []int32:
		n.setJsonValue(numbers2floats(value.([]int32)))
	case []int64:
		n.setJsonValue(numbers2floats(value.([]int64)))
	case []float32:
		n.setJsonValue(numbers2floats(value.([]float32)))
	case []int8:
		n.setJsonValue(numbers2floats(value.([]int8)))
	case []int16:
		n.setJsonValue(numbers2floats(value.([]int16)))
	case []uint16:
		n.setJsonValue(numbers2floats(value.([]uint16)))
	case []uint64:
		n.setJsonValue(numbers2floats(value.([]uint64)))
	default:
		n.setJsonValue(value)
	}
	return
}

func (n *JsonNode) setJsonValue(value interface{}) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("set error ", er))
		}
	}()
	switch value.(type) {
	case string:
		n.ValueString = value.(string)
		n.NType = STRING
	case map[string]interface{}:
		n.NType = STRUCT
		praseStruct(value.(map[string]interface{}), n)
	case []interface{}:
		arr := value.([]interface{})
		v := arr[0]
		switch v.(type) {
		case string:
			var er error
			n.ArraysString, er = interfaces2strings(arr)
			if er != nil {
				n.ArraysStruct = interfaces2array(arr)
				n.NType = STRUCTARRAY
			} else {
				n.NType = STRINGARRAY
			}
		case float64:
			var er error
			n.ArraysNumber, er = interfaces2floats(arr)
			if er != nil {
				n.ArraysStruct = interfaces2array(arr)
				n.NType = STRUCTARRAY
			} else {
				n.NType = NUMBERARRAY
			}
		case bool:
			var er error
			n.ArraysBool, er = interfaces2bools(arr)
			if er != nil {
				n.ArraysStruct = interfaces2array(arr)
				n.NType = STRUCTARRAY
			} else {
				n.NType = BOOLARRAY
			}
		case map[string]interface{}:
			var er error
			n.ArraysStruct, er = interfaces2struct(arr)
			if er != nil {
				n.ArraysStruct = interfaces2array(arr)
				n.NType = STRUCTARRAY
			} else {
				n.NType = STRUCTARRAY
			}
		case []interface{}:
			n.ArraysStruct = interfaces2array(arr)
			n.NType = STRUCTARRAY
		}
	case float64:
		n.ValueNumber = value.(float64)
		n.NType = NUMBER
	case bool:
		n.ValueBool = value.(bool)
		n.NType = BOOL
	}
	return
}

func numbers2floats(inters interface{}) (floats []float64) {
	switch inters.(type) {
	case []int:
		ints := inters.([]int)
		floats = make([]float64, len(ints))
		for i, v := range floats {
			floats[i] = v
		}
	case []int32:
		ints := inters.([]int32)
		floats = make([]float64, len(ints))
		for i, v := range floats {
			floats[i] = v
		}
	case []int64:
		ints := inters.([]int64)
		floats = make([]float64, len(ints))
		for i, v := range floats {
			floats[i] = v
		}
	case []float32:
		ints := inters.([]float32)
		floats = make([]float64, len(ints))
		for i, v := range floats {
			floats[i] = v
		}
	case []int8:
		ints := inters.([]int8)
		floats = make([]float64, len(ints))
		for i, v := range floats {
			floats[i] = v
		}
	case []int16:
		ints := inters.([]int16)
		floats = make([]float64, len(ints))
		for i, v := range floats {
			floats[i] = v
		}
	case []uint16:
		ints := inters.([]uint16)
		floats = make([]float64, len(ints))
		for i, v := range floats {
			floats[i] = v
		}
	case []uint64:
		ints := inters.([]uint64)
		floats = make([]float64, len(ints))
		for i, v := range floats {
			floats[i] = v
		}
	}
	return
}

func interfaces2floats(inters []interface{}) (floats []float64, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint("interfaces2strings error ", er))
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
		praseStruct(v.(map[string]interface{}), jsonnode)
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

func interfaces2oArray(inters []interface{}) (n []Inode) {
	for i, v := range inters {
		jsonnode := new(JsonNode)
		jsonnode.SetValue(v)
		n[i] = jsonnode
	}
	return
}

func (n *JsonNode) ToString() (s string) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println("ToString error ", er)
		}
	}()
	return _toString(n)
}

func _toString(n *JsonNode) (s string) {
	switch n.NType {
	case STRUCT:
		if n.Name != "" {
			s = fmt.Sprint(s, "\"", n.Name, "\":", "{")
		} else {
			s = fmt.Sprint(s, "{")
		}

		dats := n.StructNodes
		i := 0
		for _, dat := range dats {
			s = fmt.Sprint(s, _toString(dat.ToJsonNode()))
			if i < len(dats)-1 {
				s = fmt.Sprint(s, ",")
			}
			i++
		}
		s = fmt.Sprint(s, "}")
	case NUMBERARRAY:
		if n.Name != "" {
			s = fmt.Sprint(s, "\"", n.Name, "\":[")
		} else {
			s = fmt.Sprint(s, "[")
		}

		for i, v := range n.ArraysNumber {
			if i < len(n.ArraysNumber)-1 {
				s = fmt.Sprint(s, v, ",")
			} else {
				s = fmt.Sprint(s, v)
			}
		}
		s = fmt.Sprint(s, "]")
	case STRING:
		if n.Name != "" {
			s = fmt.Sprint(s, "\"", n.Name, "\":", "\"", n.ValueString, "\"")
		} else {
			s = fmt.Sprint(s, "\"", n.ValueString, "\"")
		}
	case NUMBER:
		if n.Name != "" {
			s = fmt.Sprint(s, "\"", n.Name, "\":", "", n.ValueNumber, "")
		} else {
			s = fmt.Sprint(s, "", n.ValueNumber, "")
		}
	case BOOL:
		if n.Name != "" {
			s = fmt.Sprint(s, "\"", n.Name, "\":", "", n.ValueBool, "")
		} else {
			s = fmt.Sprint(s, "", n.ValueBool, "")
		}
	case STRUCTARRAY:
		if n.Name != "" {
			s = fmt.Sprint(s, "\"", n.Name, "\":[")
		} else {
			s = fmt.Sprint(s, "[")
		}
		for i, v := range n.ArraysStruct {
			s = fmt.Sprint(s, _toString(v.ToJsonNode()))
			if i < len(n.ArraysStruct)-1 {
				s = fmt.Sprint(s, ",")
			}
		}
		s = fmt.Sprint(s, "]")
	case STRINGARRAY:
		if n.Name != "" {
			s = fmt.Sprint(s, "\"", n.Name, "\":[")
		} else {
			s = fmt.Sprint(s, "[")
		}
		for i, v := range n.ArraysString {
			if i < len(n.ArraysString)-1 {
				s = fmt.Sprint(s, "\"", v, "\"", ",")
			} else {
				s = fmt.Sprint(s, "\"", v, "\"")
			}
		}
		s = fmt.Sprint(s, "]")
	case BOOLARRAY:
		if n.Name != "" {
			s = fmt.Sprint(s, "\"", n.Name, "\":[")
		} else {
			s = fmt.Sprint(s, "[")
		}
		for i, v := range n.ArraysBool {
			if i < len(n.ArraysBool)-1 {
				s = fmt.Sprint(s, v, ",")
			} else {
				s = fmt.Sprint(s, v)
			}
		}
		s = fmt.Sprint(s, "]")
	}
	return
}

func (n *JsonNode) GetChildNames() (names []string) {
	if n.StructNodes != nil && len(n.StructNodes) > 0 {
		names = make([]string, len(n.StructNodes))
		i := 0
		for k, _ := range n.StructNodes {
			names[i] = k
			i++
		}
	}
	return
}

func (n *JsonNode) IsExist(path string) bool {
	if n.GetNodeByPath(path) == nil {
		return false
	} else {
		return true
	}
}

func isTypeStruct(value interface{}) bool {
	switch value.(type) {
	case map[string]interface{}:
		return true
	}
	return false
}

func praseStruct(dat map[string]interface{}, superNode *JsonNode) {
	superNode.StructNodes = make(map[string]Inode)
	nodes := superNode.StructNodes
	for k, v := range dat {
		json := new(JsonNode)
		json.Name = k
		json.setJsonValue(v)
		nodes[k] = json
		if isTypeStruct(v) {
			praseStruct(v.(map[string]interface{}), json)
		}
	}
}

//根据路径查询子节点  如：a.b  b节点为a节点的子节点
func (n *JsonNode) GetNodeByPath(path string) (node *JsonNode) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println("GetNode error ", path)
		}
	}()
	paths := strings.Split(path, ".")
	node = n
	for _, p := range paths {
		name := strings.TrimSpace(p)
		node = getChildJsonNode(node, name)
	}
	return
}

//根据节点名字查询
func (n *JsonNode) GetNodeByName(name string) *JsonNode {
	defer func() {
		if er := recover(); er != nil {
			fmt.Println(er)
		}
	}()
	if n.NType == STRUCT {
		chl := n.StructNodes
		if node, ok := chl[name]; ok {
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
	nodes := n.StructNodes
	if nodes != nil {
		node, ok := nodes[name]
		if ok {
			jnode := node.(*JsonNode)
			return jnode
		}
	}
	return nil
}
