package config

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type ExportAttribute map[string]interface{}

func (att ExportAttribute) MustDecode(target interface{}) {

	err := att.Decode(target)
	if err != nil {
		panic(err)
	}
}

func (att ExportAttribute) Decode(target interface{}) error {

	rv := reflect.TypeOf(target)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("decode param must be a ptr")
	}

	err := mapstructure.Decode(att, target)
	if err != nil {
		return err
	}
	return nil
}

// config for exporter
type ExportConfig struct {
	Name      string
	Attribute map[string]interface{}
}

// io exporter
type StdoutExporterConfig struct {
	// FileName specify the place where trace export to ,if empty default stdout
	FileName string

	PrettyPrint bool

	Timestamps bool
}

// jeagerExporter
type JeagerExporterConfig struct {
	Url string
}
