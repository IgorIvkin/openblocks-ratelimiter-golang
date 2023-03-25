package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type Application struct {
	Config  *ApplicationConfig
	Buckets map[string]*RateLimiterBucket
}

type RateLimiterBucket struct {
	Limit        int64
	Tokens       int64
	Window       int64
	LastCallTime time.Time
	sync.Mutex
}

type RateLimitConfig struct {
	Limit int64  `yaml:"limit"`
	Unit  string `yaml:"unit"`
}

type ApplicationConfig struct {
	Limiters map[string]RateLimitConfig `yaml:"limiters"`
}

func NewApplication() *Application {
	config := getConfig()
	buckets := initializeBuckets(config)
	app := Application{
		Config:  config,
		Buckets: buckets,
	}
	return &app
}

// Получает новый токен из бакета.
func (bucket *RateLimiterBucket) GetToken() int64 {
	bucket.Lock()
	defer bucket.Unlock()

	bucket.UpdateTokens()
	bucket.Tokens = bucket.Tokens - 1
	return bucket.Tokens
}

// Пополняет токены в бакете в том случае, если временной диапазон
// позволяет это сделать, так, для временного диапазона в минуту
// новые токены пополняются раз в минуту.
func (bucket *RateLimiterBucket) UpdateTokens() {
	currentTime := time.Now()
	fmt.Println(currentTime.Sub(bucket.LastCallTime))
	fmt.Println(time.Duration(bucket.Window))
	if currentTime.Sub(bucket.LastCallTime) > time.Duration(bucket.Window) {
		bucket.Tokens = bucket.Limit
		bucket.LastCallTime = currentTime
	}
}

// Инициализирует бакеты по конфигурации, которая читается из yml-файла приложения,
// в каждый бакет записывается текущее количество токенов, лимиты по времени и
// текущее время.
func initializeBuckets(config *ApplicationConfig) map[string]*RateLimiterBucket {
	buckets := make(map[string]*RateLimiterBucket)

	limiterConfig := config.Limiters
	for limiterName, limiterConfig := range limiterConfig {
		window, err := getWindowByUnit(limiterConfig.Unit)
		if err != nil {
			log.Fatal(err)
		}
		bucket := RateLimiterBucket{
			Limit:        limiterConfig.Limit,
			Tokens:       limiterConfig.Limit,
			Window:       window,
			LastCallTime: time.Now(),
		}
		buckets[limiterName] = &bucket
	}

	return buckets
}

// Преобразует переданные в строковом виде единицы времени во временной
// диапазон, в течение которого проверяется ограничение в рейт-лимитере.
// Данные возвращаются в нано-секундах.
func getWindowByUnit(unit string) (int64, error) {
	if unit == "SECONDS" {
		return 1000 * 1_000_000, nil
	} else if unit == "MINUTES" {
		return 1000 * 60 * 1_000_000, nil
	} else if unit == "HOURS" {
		return 1000 * 60 * 60 * 1_000_000, nil
	} else {
		return 0, fmt.Errorf("unknown time unit - %s", unit)
	}
}

// Возвращает конфигурацию, заданную в yml-файле приложения,
// в ней можно задать множественные рейт-лимитеры.
func getConfig() *ApplicationConfig {
	configFile, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	var config ApplicationConfig
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
