package slogTracer

import "strings"

type AccessChecker interface {
	Check(value string) bool
	Header() string
	Value() string
}

type accessChecker struct {
	header string
	value  string
}

func NewStaticAccessChecker(header, value string) AccessChecker {
	return &accessChecker{
		header: header,
		value:  value,
	}
}

func (ac *accessChecker) Check(value string) bool {
	return strings.EqualFold(ac.value, ac.value)
}

func (ac *accessChecker) Value() string {
	return ac.value
}

func (ac *accessChecker) Header() string {
	return ac.header
}
