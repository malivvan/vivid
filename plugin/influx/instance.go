package influx

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type PluginInstance struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
}

type Measurement struct {
	instance *PluginInstance
	name     string
	tags     map[string]string
	fields   map[string]interface{}
	ts       time.Time
}

func (m *Measurement) Tag(key, value string) *Measurement {
	m.tags[key] = value
	return m
}

func (m *Measurement) Field(key string, value interface{}) *Measurement {
	m.fields[key] = value
	return m
}

func (m *Measurement) Timestamp(t time.Time) *Measurement {
	m.ts = t
	return m
}

func (m *Measurement) Write() error {
	return m.instance.writeAPI.WritePoint(context.Background(), influxdb2.NewPoint(m.name, m.tags, m.fields, m.ts))
}

func (pi *PluginInstance) Measurement(name string) *Measurement {
	return &Measurement{
		ts:       time.Now(),
		instance: pi,
		name:     name,
		tags:     make(map[string]string),
		fields:   make(map[string]interface{}),
	}
}

func (pi *PluginInstance) Flush() error {
	return pi.writeAPI.Flush(context.Background())
}

func (pi *PluginInstance) Destroy() error {
	pi.client.Close()
	return nil
}
