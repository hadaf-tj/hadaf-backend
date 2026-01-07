package configs

import "os"

var (
	AccessSecret  = os.Getenv("ACCESS_SECRET")
	AccessExpire  = os.Getenv("ACCESS_EXPIRE")
	RefreshSecret = os.Getenv("REFRESH_SECRET")
	RefreshExpire = os.Getenv("REFRESH_EXPIRE")
)

var (
	MinioAccessKey = os.Getenv("MINIO_ACCESS_KEY")
	MinioSecretKey = os.Getenv("MINIO_SECRET_KEY")
	MinioEndpoint  = os.Getenv("MINIO_ENDPOINT")
	MinioBucket    = os.Getenv("MINIO_BUCKET")
)

var (
	RedisHost      = os.Getenv("REDIS_HOST")
	RedisDefaultDB = os.Getenv("REDIS_DEFAULT_DB")
	RedisTimeout   = os.Getenv("REDIS_TIMEOUT")
)
