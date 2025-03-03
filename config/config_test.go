package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	os.Setenv("LISTEN_ADDRESS", "0.0.0.0")
	os.Setenv("HTTP_PORT", "8000")

	type Port int
	type Cfg struct {
		ListenAddress     string `env:"LISTEN_ADDRESS" default:"0.0.0.0"`
		GrpcPort          Port   `env:"GRPC_PORT" default:"4040"`
		HttpPort          Port   `env:"HTTP_PORT" default:"8080"`
		HttpBasicUsername string `env:"HTTP_BASIC_USERNAME" default:"root"`
		HttpBasicPassword string `env:"HTTP_BASIC_PASSWORD" default:"toor"`
		Endpoint          string `env:"CLUSTER_ENDPOINT" default:"http://localhost:8888"`
	}

	want := &Cfg{
		ListenAddress:     "0.0.0.0",
		GrpcPort:          Port(4040),
		HttpPort:          Port(8000),
		HttpBasicUsername: "root",
		HttpBasicPassword: "toor",
		Endpoint:          "http://localhost:8888",
	}

	got, err := LoadConfig[Cfg]()
	if err != nil {
		t.Fatalf("cannot parse config: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadOtelProxyConfig() = %+v, want %+v", got, want)
	}

}
