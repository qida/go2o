/**
 * Copyright 2014 @ ops Inc.
 * name : delivery
 * author : newmin
 * date : 2014-10-06 14:21 :
 * description :
 * history :
 */
package domain

import(
    "regexp"
    "errors"
)


var(
    areaRegexp = regexp.MustCompile("(市)((.+)(区|县))")
    errNotMatch = errors.New("未识别的地址")
    cityRegexp = regexp.MustCompile("(省|自治区|行政区)((.+)市)")
)


// 获取地区名称
func GetAreaName(addr string)(string,error){
    var matches [][]string = areaRegexp.FindAllStringSubmatch(addr,-1)
    if len(matches) == 0 {
        return "",errNotMatch
    }
    return matches[0][2]
}

// 获取城市名称
func GetCityName(addr string)(string,error){
    var matches [][]string = cityRegexp.FindAllStringSubmatch(addr,-1)
    if len(matches) == 0 {
        return "",errNotMatch
    }
    return matches[0][2]
}