package main

import (
	"fmt"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var testAnnotations = map[string]string{
	"fqdn.io/key.value":       "endpoint",
	"other.fqdn.io/key.value": "endpoint",
}

func makeAnnotations(annotations map[string]string) *ServiceAnnotations {
	sas := &ServiceAnnotations{}
	for key, value := range annotations {
		sas.Add(&ServiceAnnotation{
			Annotation:  KeyPair{key, value},
			ServiceName: "irrelevant"})
	}
	return sas
}

func TestServiceAnnotations_Add(t *testing.T) {
	sa := &ServiceAnnotation{}

	sas := ServiceAnnotations{}
	sas.Add(sa)

	if len(sas.Annotations) != 1 {
		t.Fail()
	}

	sas.Add(sa, sa, sa)
	if len(sas.Annotations) != 4 {
		t.Fail()
	}
}

func boolFilter(rv bool) AnnotationFilter {
	return func(sa *ServiceAnnotation) bool {
		return rv
	}
}

func TestServiceAnnotations_Filter(t *testing.T) {
	sas := makeAnnotations(testAnnotations)
	trueFilter := boolFilter(true)
	falseFilter := boolFilter(false)

	singleFilterResult := sas.Filter(trueFilter)
	if len(singleFilterResult.Annotations) != len(testAnnotations) {
		t.Fail()
	}

	manyTrueFilterResult := sas.Filter(trueFilter, trueFilter)
	if len(manyTrueFilterResult.Annotations) != len(testAnnotations) {
		t.Fail()
	}

	eventuallyFalseFilterResult := sas.Filter(trueFilter, falseFilter)
	if len(eventuallyFalseFilterResult.Annotations) != 0 {
		t.Fail()
	}

}

type filterSubject struct {
	sa       *ServiceAnnotation
	fn       AnnotationFilter
	expected bool
}

func keyAnt(key string) *ServiceAnnotation {
	return &ServiceAnnotation{Annotation: KeyPair{key, "value"}}
}

var testFqdnFilterSubjects = []filterSubject{
	{keyAnt("go.test.io/key.value"), FqdnFilter("go.test.io"), true},
	{keyAnt("go.test.io/key.value"), FqdnFilter("go.test"), false},
	{keyAnt("go.test.io/key.value"), FqdnFilter("test.io"), false},
	{keyAnt("namespace/key.value"), FqdnFilter("namespace"), true},
}

var testKeyFilterSubjects = []filterSubject{
	{keyAnt("n/key.value"), KeyFilter("key"), true},
	{keyAnt("n/key.value.value"), KeyFilter("key"), true},
	{keyAnt("n/key.value.value"), KeyFilter("key.value"), false},
	{keyAnt("n/key-key.value.value"), KeyFilter("key-key"), true},
}

var testKeyValueFilterSubjects = []filterSubject{
	{keyAnt("n/key.value"), KeyValueFilter("key", "value"), true},
	{keyAnt("n/key.value.value"), KeyValueFilter("key", "value"), false},
	{keyAnt("n/key-key.value"), KeyValueFilter("key-key", "value"), true},
}

var testValidAtlasFilterSubjects = []filterSubject{
	{keyAnt("namespace/key.value"), ValidAtlasFilter(), true},
	{keyAnt("namespace.tld/key.value"), ValidAtlasFilter(), true},
	{keyAnt("sub.namespace.tld/key.value"), ValidAtlasFilter(), true},
	{keyAnt("namespace/kevalue"), ValidAtlasFilter(), false},
	{keyAnt("namespace.key.value"), ValidAtlasFilter(), false},
}

func testFilterSubjects(subjects []filterSubject, t *testing.T) {
	for _, s := range subjects {
		result := s.fn(s.sa)
		if result != s.expected {
			t.Logf("Failed: %v\n\t expected %v got %v", s, s.expected, result)
			t.Fail()
		}
	}
}

func TestFqdnFilter(t *testing.T) {
	testFilterSubjects(testFqdnFilterSubjects, t)
}

func TestKeyFilter(t *testing.T) {
	testFilterSubjects(testKeyFilterSubjects, t)
}

func TestKeyValueFilter(t *testing.T) {
	testFilterSubjects(testKeyValueFilterSubjects, t)
}

func TestValidAtlasFilter(t *testing.T) {
	testFilterSubjects(testValidAtlasFilterSubjects, t)

	// verify filter mutation -- may change
	for _, annotation := range testValidAtlasFilterSubjects {
		fmt.Print(annotation.sa.Key, annotation.sa.Value)
		if annotation.expected && annotation.sa.Key == "" {
			t.Fail()
		}
	}
}

func serviceWithAnnotations(annotations map[string]string) v1.Service {
	obj := v1.Service{
		ObjectMeta: v12.ObjectMeta{
			Name:        "servicename",
			Namespace:   v1.NamespaceDefault,
			Annotations: annotations}}
	return obj
}

func TestAnnotationsFromService(t *testing.T) {
	ant := map[string]string{
		"go.test.io/a.b.c": "value",
		"go.test.io/d.e.f": "value",
	}
	svc := serviceWithAnnotations(ant)

	annotations := AnnotationsFromService(svc)

	if len(annotations) != 2 {
		t.Fail()
	}
}

func TestClient_allAnnotations(t *testing.T) {
	ant := map[string]string{
		"go.test.io/a.b.c": "value",
		"go.test.io/d.e.f": "value",
	}
	svc := serviceWithAnnotations(ant)
	svcList := &v1.ServiceList{
		Items: []v1.Service{svc},
	}
	c := Client{
		services: svcList,
	}
	sas := c.allAnnotations()
	if len(sas.Annotations) != 2 {
		t.Fail()
	}
}
