package configs

type Config struct {
    App      AppConfig
    Security SecurityConfig // <--- Вот это поле ищет компилятор!
    Database DatabaseConfig
}

type AppConfig struct {
    Port string
}

type SecurityConfig struct {
    SecretKey string
}

type DatabaseConfig struct {
    DSN string
}

// Заглушка, чтобы просто запустилось
func Load() *Config {
    return &Config{
        App: AppConfig{Port: ":8080"},
        Security: SecurityConfig{
            SecretKey: "secret-key-123", // Любой ключ
        },
        Database: DatabaseConfig{DSN: ""},
    }
}