// Copyright 2016 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen_helper.go

package ptyy

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var _ = fmt.Sprintf

// 医院信息
type HospitalInfo struct {
	Name string // 名称
	City string // 城市
}

// 医院列表
var All []HospitalInfo = _AllHospitalInfoList

// 查询列表
func Search(query string, limits int) []HospitalInfo {
	// 规范化: 删除前后空白
	query = strings.TrimSpace(query)

	// 如果为空的话, 返回全部
	if query == "" {
		if limits <= 0 || limits > len(All) {
			return All
		}
		return All[:limits]
	}

	// 根据关键字查询
	if results := searchByKeywords(query, limits); len(results) > 0 {
		return results
	}

	// 如果没有匹配的, 则尝试正则查询
	if re, err := regexp.Compile(query); err == nil {
		if results := searchByRegexp(re, limits); len(results) > 0 {
			return results
		}
	}

	// 没有匹配
	return nil
}

// 根据关键字查询
func SearchByKeywords(keywords string, limits int) []HospitalInfo {
	// 规范化: 删除前后空白
	keywords = strings.TrimSpace(keywords)

	// 如果为空的话, 返回全部
	if keywords == "" {
		if limits <= 0 || limits > len(All) {
			return All
		}
		return All[:limits]
	}

	// 根据关键字查询
	if results := searchByKeywords(keywords, limits); len(results) > 0 {
		return results
	}

	// 没有匹配
	return nil
}

// 根据关键字查询
func SearchByRegexp(query string, limits int) ([]HospitalInfo, error) {
	// 规范化: 删除前后空白
	query = strings.TrimSpace(query)

	// 如果为空的话, 返回全部
	if query == "" || query == ".*" {
		if limits <= 0 || limits > len(All) {
			return All, nil
		}
		return All[:limits], nil
	}

	// 如果没有匹配的, 则尝试正则查询
	re, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}
	if results := searchByRegexp(re, limits); len(results) > 0 {
		return results, nil
	}

	// 没有匹配
	return nil, nil
}

// 根据关键字查询
// TODO: 以后可扩展为多个关键字, 采用非字符字符分隔
func searchByKeywords(key string, limits int) (results []HospitalInfo) {
	result0Map := make(map[string]HospitalInfo)
	result1Map := make(map[string]HospitalInfo)

	for _, v := range All {
		if limits > 0 && len(result0Map)+len(result1Map) >= limits {
			break
		}
		if strings.HasPrefix(v.Name, key) || strings.HasPrefix(v.City, key) {
			result0Map[v.Name] = v
		}
		if strings.HasPrefix(_NamePinyinLongMap[v.Name], key) || strings.HasPrefix(_NamePinyinShortMap[v.Name], key) {
			result0Map[v.Name] = v
		}
		if strings.HasPrefix(_NamePinyinLongMap[v.City], key) || strings.HasPrefix(_NamePinyinShortMap[v.City], key) {
			result0Map[v.Name] = v
		}
	}
	for _, v := range All {
		if limits > 0 && len(result0Map)+len(result1Map) >= limits {
			break
		}
		if strings.Contains(v.Name, key) || strings.Contains(v.City, key) {
			if _, ok := result0Map[v.Name]; !ok {
				result1Map[v.Name] = v
			}
		}
		if strings.Contains(_NamePinyinLongMap[v.Name], key) || strings.Contains(_NamePinyinShortMap[v.Name], key) {
			if _, ok := result0Map[v.Name]; !ok {
				result1Map[v.Name] = v
			}
		}
		if strings.Contains(_NamePinyinLongMap[v.City], key) || strings.Contains(_NamePinyinShortMap[v.City], key) {
			if _, ok := result0Map[v.Name]; !ok {
				result1Map[v.Name] = v
			}
		}
	}

	var result0List []HospitalInfo
	var result1List []HospitalInfo

	for _, v := range result0Map {
		result0List = append(result0List, v)
	}
	for _, v := range result1Map {
		result1List = append(result1List, v)
	}

	sort.Sort(byHospitalInfo(result0List))
	sort.Sort(byHospitalInfo(result1List))

	results = append(results, result0List...)
	results = append(results, result1List...)

	return
}

// 根据正则表达式查询
func searchByRegexp(re *regexp.Regexp, limits int) []HospitalInfo {
	resultMap := make(map[string]HospitalInfo)

	for _, v := range All {
		if limits > 0 && len(resultMap) >= limits {
			break
		}
		if re.MatchString(v.Name) || re.MatchString(v.City) {
			resultMap[v.Name] = v
		}
		if re.MatchString(_NamePinyinLongMap[v.Name]) || re.MatchString(_NamePinyinShortMap[v.Name]) {
			resultMap[v.Name] = v
		}
		if re.MatchString(_NamePinyinLongMap[v.City]) || re.MatchString(_NamePinyinShortMap[v.City]) {
			resultMap[v.Name] = v
		}
	}

	var resultList []HospitalInfo
	for _, v := range resultMap {
		resultList = append(resultList, v)
	}
	sort.Sort(byHospitalInfo(resultList))
	return resultList
}

// 按unicode排序
type byHospitalInfo []HospitalInfo

func (d byHospitalInfo) Len() int           { return len(d) }
func (d byHospitalInfo) Less(i, j int) bool { return d[i].Name < d[j].Name }
func (d byHospitalInfo) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
