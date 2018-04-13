package main

import (
	"k8s.io/api/core/v1"
	"regexp"
)

type KeyPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ServiceAnnotation struct {
	Key         string  `json:"key"`
	Value       string  `json:"value"`
	Annotation  KeyPair `json:"annotation"`
	ServiceName string  `json:"service_name"`
}

type ServiceAnnotations struct {
	Annotations []*ServiceAnnotation `json:"annotations"`
}

func (sas *ServiceAnnotations) Add(annotations ...*ServiceAnnotation) *ServiceAnnotations {
	for _, annotation := range annotations {
		sas.Annotations = append(sas.Annotations, annotation)
	}
	return sas
}

type AnnotationFilter func(annotation *ServiceAnnotation) bool

func (sas *ServiceAnnotations) Filter(filters ...AnnotationFilter) *ServiceAnnotations {
	sa := &ServiceAnnotations{}
	for _, annotation := range sas.Annotations {
		add := true
		for _, filter := range filters {
			if !filter(annotation) {
				add = false
				break
			}
		}
		if add {
			sa.Add(annotation)
		}
	}
	sas.Annotations = sa.Annotations
	return sas
}

func FqdnFilter(fqdn string) AnnotationFilter {
	return func(annotation *ServiceAnnotation) bool {
		prefix, _ := Split2(annotation.Annotation.Key, "/")
		return prefix == fqdn
	}
}

func KeyFilter(filterKey string) AnnotationFilter {
	return func(annotation *ServiceAnnotation) bool {
		_, keyValue := Split2(annotation.Annotation.Key, "/")
		key, _ := Split2(keyValue, ".")
		return key == filterKey
	}
}

func KeyValueFilter(filterKey, filterValue string) AnnotationFilter {
	return func(annotation *ServiceAnnotation) bool {
		_, keyValue := Split2(annotation.Annotation.Key, "/")
		key, value := Split2(keyValue, ".")
		return key == filterKey && value == filterValue
	}
}

func ValidAtlasFilter() AnnotationFilter {
	pattern := `([a-zA-Z][-a-zA-Z0-9.]*)\/([-a-zA-Z0-9]+)\.([-a-zA-Z0-9.]*)`
	re := regexp.MustCompile(pattern)
	return func(annotation *ServiceAnnotation) bool {
		match := re.FindStringSubmatch(annotation.Annotation.Key)
		if len(match) == 0 {
			return false
		}

		// Mutates the annotation -- not sure if I like this
		annotation.Key = match[2]
		annotation.Value = match[3]
		return true
	}
}

func AnnotationsFromService(service v1.Service) []*ServiceAnnotation {
	annotations := &ServiceAnnotations{}
	for key, value := range service.Annotations {
		annotations.Add(&ServiceAnnotation{
			Annotation:  KeyPair{key, value},
			ServiceName: service.Name,
		})
	}
	return annotations.Annotations
}

func (c *Client) allAnnotations() *ServiceAnnotations {

	// this if block is done to allow testing
	if c.services == nil {
		c.GetServices()
	}
	annotations := &ServiceAnnotations{}
	for _, service := range c.services.Items {
		annotations.Add(AnnotationsFromService(service)...)
	}
	return annotations
}
