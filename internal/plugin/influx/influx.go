package influx

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Plugin struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
}

type Config struct {
	Host   string
	Token  string `encrypt:"9rcWH9mkpE2ZUn3z8x2SnTcUHvpqBFYs"`
	Org    string
	Bucket string
}

type Measurement struct {
	plugin *Plugin
	name   string
	tags   map[string]string
	fields map[string]interface{}
	ts     time.Time
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
	return m.plugin.writeAPI.WritePoint(context.Background(), influxdb2.NewPoint(m.name, m.tags, m.fields, m.ts))
}

func New(config Config) *Plugin {
	client := influxdb2.NewClient(config.Host, config.Token)
	writeAPI := client.WriteAPIBlocking(config.Org, config.Bucket)
	return &Plugin{
		client:   client,
		writeAPI: writeAPI,
	}
}

func (p *Plugin) Measurement(name string) *Measurement {
	return &Measurement{
		ts:     time.Now(),
		plugin: p,
		name:   name,
		tags:   make(map[string]string),
		fields: make(map[string]interface{}),
	}
}

func (p *Plugin) Flush() error {
	return p.writeAPI.Flush(context.Background())
}

func (p *Plugin) Destroy() {
	p.client.Close()
}
