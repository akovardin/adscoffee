package ads

import (
	"encoding/json"
	"net"
	"slices"
)

type Targeting struct {
	Bundle   ExcludeInclude
	Audience ExcludeInclude
	Bapp     ExcludeInclude
	IP       ExcludeIncludeIP
	Country  ExcludeInclude
	City     ExcludeInclude
	Region   ExcludeInclude
	Network  ExcludeInclude
}

func NewTargeting(data string) (Targeting, error) {
	t := Targeting{}

	if data == "" {
		return t, nil
	}

	if err := json.Unmarshal([]byte(data), &t); err != nil {
		return Targeting{}, err
	}

	return t, nil
}

type ExcludeInclude struct {
	IncludeOr  []string
	ExcludeOr  []string
	IncludeAnd []string
	ExcludeAnd []string
}

func (e ExcludeInclude) Validate(values []string) bool {
	if len(e.ExcludeAnd) > 0 && arrayInArrayAnd(values, e.ExcludeAnd) {
		return false // в установленных приложениях есть все из исключения
	}

	if len(e.IncludeAnd) > 0 && !arrayInArrayAnd(values, e.IncludeAnd) {
		return false // в установленных приложениях нет всех из включенных
	}

	// проверяем если заполнены вхождения по OR
	if len(e.ExcludeOr) > 0 && arrayInArrayOr(values, e.ExcludeOr) {
		return false // в установленных приложениях есть все из исключения
	}

	if len(e.IncludeOr) > 0 && !arrayInArrayOr(values, e.IncludeOr) {
		return false // в установленных приложениях нет всех из включенных
	}

	return true
}

type ExcludeIncludeIP struct {
	Include []*net.IPNet
	Exclude []*net.IPNet
}

func (e ExcludeIncludeIP) Validate(ip string) bool {
	if len(e.Exclude) > 0 && ipInNetworks(e.Exclude, ip) {
		return false // в ip есть в исключениях
	}

	if len(e.Include) > 0 && !ipInNetworks(e.Include, ip) {
		return false // ip нет во включенных
	}

	return true
}

func ipInNetworks(networks []*net.IPNet, val string) bool {
	ip := net.ParseIP(val)
	if ip == nil {
		return false
	}

	for _, network := range networks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// Возвращает true если есть хотя бы одно значение из search входит в массив items.
func arrayInArrayOr[T comparable](items []T, search []T) bool {
	for _, value := range search {
		if slices.Contains(items, value) {
			return true
		}
	}

	return false
}

// Возвращает true если есть все значения из search входят в массив items.
func arrayInArrayAnd[T comparable](items []T, search []T) bool {
	for _, value := range search {
		if !slices.Contains(items, value) {
			return false
		}
	}

	return true
}
