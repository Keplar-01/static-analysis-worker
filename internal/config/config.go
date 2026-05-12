package config

import "os"

type Config struct {
	KafkaBrokers   string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	AnalyzerBinary string
}

func Load() *Config {
	return &Config{
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),

		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),

		AnalyzerBinary: getEnv("ANALYZER_BINARY", "/usr/local/bin/analyzer"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
