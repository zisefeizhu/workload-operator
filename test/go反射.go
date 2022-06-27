package main

//
//import (
//	"fmt"
//	appv1 "k8s.io/api/apps/v1"
//	"reflect"
//)
//
//// 反射传入的结构体并实例化一个实例，根据结构体字段在CSV行中找对应数据，找到的就设值
//func instantiate(unknownSt interface{}, jsonData map[string]interface{}) interface{} {
//	// 获取结构体的类型反射
//	sType := reflect.TypeOf(unknownSt).Elem()
//	// 根据类型实例化一个新结构体
//	retSt := reflect.New(sType).Elem()
//	// 遍历每一个结构体类型成员
//	for i := 0; i < sType.NumField(); i++ {
//		//	// 成员类型
//		f := sType.Field(i)
//		//	// 新结构体的对应成员
//		fv := retSt.Field(i)
//		// 查找数据Map，赋值
//		if v, ok := jsonData[f.Tag.Get("json")]; ok && v != nil {
//			fv.Set(reflect.ValueOf(v))
//		}
//	}
//	return retSt
//}
//
////type tSt struct {
////	A string `json:"a"`
////	B int    `json:"b"`
////}
//
//func main() {
//	a := &appv1.Deployment{
//		Spec: appv1.DeploymentSpec{
//			Paused:          true,
//			MinReadySeconds: 2,
//		},
//	}
//
//	jd := make(map[string]interface{}, 0)
//	jd["spec"] = a
//	//jd["b"] = b
//	ret := instantiate(&appv1.Deployment{}, jd)
//	fmt.Println(ret)
//}
