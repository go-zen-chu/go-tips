package externallibrary

import (
	"fmt"

	"github.com/go-zen-chu/go-tips/unittest/externalpackage"
)

func ActualPattern() {
	s := NewSample(&actualExternalLibrary{})
	s.MethodNeedToBeMocked()
}

func TestPattern() {
	s := NewSample(&mockExternalLibrary{})
	s.MethodNeedToBeMocked()
}

type Sample struct {
	el ExternalLibrary
}

// 外部ライブラリの interface を DI させることが大切
func NewSample(el ExternalLibrary) *Sample {
	return &Sample{
		el: el,
	}
}

func (s *Sample) MethodNeedToBeMocked() error {
	if err := s.el.ExternalFunc1(); err != nil {
		return fmt.Errorf("external func1: %w", err)
	}
	if err := s.el.ExternalFunc2(); err != nil {
		return fmt.Errorf("external func2: %w", err)
	}
	if err := s.el.ExternalFunc3(); err != nil {
		return fmt.Errorf("external func3: %w", err)
	}
	return nil
}

type ExternalLibrary interface {
	ExternalFunc1() error
	ExternalFunc2() error
	ExternalFunc3() error
}

// actual struct uses actual external library's methods
type actualExternalLibrary struct{}

func (a *actualExternalLibrary) ExternalFunc1() error {
	return externalpackage.ExternalFunc1()
}

func (a *actualExternalLibrary) ExternalFunc2() error {
	return externalpackage.ExternalFunc2()
}
func (a *actualExternalLibrary) ExternalFunc3() error {
	return externalpackage.ExternalFunc3()
}

// mock struct is defined for unittesting
type mockExternalLibrary struct{}

func (m *mockExternalLibrary) ExternalFunc1() error {
	return nil
}
func (m *mockExternalLibrary) ExternalFunc2() error {
	return nil
}
func (m *mockExternalLibrary) ExternalFunc3() error {
	return nil
}
