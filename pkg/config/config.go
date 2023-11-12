package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// PGConfig representing a postgres configuration
type PGConfig struct {
	URI string
}

// WebConfig representing a web configuration
type WebConfig struct {
	Host string
	Port string
}

func (c WebConfig) Address() string {
	return fmt.Sprintf("%s%v", c.Host, c.Port)
}

// OtelConfig representing an open telemetry configuration
type OtelConfig struct {
	ServiceName      string
	ExporterEndpoint string
}

// IAMConfig representing an identity configuration
type IAMConfig struct {
	Domain   string
	Audience string
}

// AppConfig representing an application configuration
type AppConfig struct {
	Environment     string
	PG              PGConfig
	Web             WebConfig
	IAM             IAMConfig
	Otel            OtelConfig
	ShutdownTimeout time.Duration
}

// ReadConfigFromEnv reads all environment variables, validates it and parses it into AppConfig struct
func ReadConfigFromEnv() (AppConfig, error) {
	environment := strings.TrimSpace(os.Getenv("ENVIRONMENT"))
	if environment == "" {
		return AppConfig{}, errors.New("environment is invalid")
	}

	port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("APP_PORT")))
	if err != nil || (port < 0 || port > 9999) {
		return AppConfig{}, errors.New("port is invalid")
	}

	pgURI := strings.TrimSpace(os.Getenv("PG_URL"))
	if pgURI == "" {
		return AppConfig{}, errors.New("pg uri is required")
	}

	otelServiceName := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME"))
	if otelServiceName == "" {
		log.Print("open telemetry service name have not been set")
	}

	otelExporterEndpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if otelExporterEndpoint == "" {
		log.Print("open telemetry exporter endpoint have not been set")
	}

	iamDomain := strings.TrimSpace(os.Getenv("IAM_DOMAIN"))
	if iamDomain == "" {
		log.Print("iam domain have not been set")
	}

	iamAudience := strings.TrimSpace(os.Getenv("IAM_AUDIENCE"))
	if iamAudience == "" {
		log.Print("iam audience have not been set")
	}

	return AppConfig{
		Environment: environment,
		Web: WebConfig{
			Host: "0.0.0.0",
			Port: fmt.Sprintf(":%v", port),
		},
		Otel: OtelConfig{
			ServiceName:      otelServiceName,
			ExporterEndpoint: otelExporterEndpoint,
		},
		IAM: IAMConfig{
			Domain:   iamDomain,
			Audience: iamAudience,
		},
		PG: PGConfig{
			URI: pgURI,
		},
	}, nil
}
