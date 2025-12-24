package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	Server          ServerConfig          `mapstructure:"server"`
	Database        DatabaseConfig        `mapstructure:"database"`
	JWT             JWTConfig             `mapstructure:"jwt"`
	Log             LogConfig             `mapstructure:"log"`
	CORS            CORSConfig            `mapstructure:"cors"`
	Swagger         SwaggerConfig         `mapstructure:"swagger"`
	Dify            DifyConfig            `mapstructure:"dify"`
	DeepSeek        DeepSeekConfig        `mapstructure:"deepseek"`
	VolcEngine      VolcEngineConfig      `mapstructure:"volcengine"`
	VolcEngineImage VolcEngineImageConfig `mapstructure:"volcengineImage"` // 图片生成配置
	BaiduOCR        BaiduOCRConfig        `mapstructure:"baiduocr"`
	ODAConverter    ODAConverterConfig    `mapstructure:"odaconverter"` // 兼容旧配置
	DWGConverter    DWGConverterConfig    `mapstructure:"dwgconverter"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Charset         string `mapstructure:"charset"`
	ParseTime       bool   `mapstructure:"parseTime"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"` // 秒
	TLS             bool   `mapstructure:"tls"`             // TiDB Cloud 等云数据库需要 TLS
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string `mapstructure:"secret"`
	AccessTokenExpire  int    `mapstructure:"accessTokenExpire"`  // 秒
	RefreshTokenExpire int    `mapstructure:"refreshTokenExpire"` // 秒
	Issuer             string `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxAge     int    `mapstructure:"maxAge"`
	MaxBackups int    `mapstructure:"maxBackups"`
	Compress   bool   `mapstructure:"compress"`
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allowOrigins"`
	AllowMethods     []string `mapstructure:"allowMethods"`
	AllowHeaders     []string `mapstructure:"allowHeaders"`
	ExposeHeaders    []string `mapstructure:"exposeHeaders"`
	AllowCredentials bool     `mapstructure:"allowCredentials"`
	MaxAge           int      `mapstructure:"maxAge"`
}

// SwaggerConfig Swagger配置
type SwaggerConfig struct {
	Title       string `mapstructure:"title"`
	Version     string `mapstructure:"version"`
	Description string `mapstructure:"description"`
	Host        string `mapstructure:"host"`
	BasePath    string `mapstructure:"basePath"`
}

// DifyConfig Dify配置
type DifyConfig struct {
	DatasetAPIKey string `mapstructure:"datasetApiKey"`
	ChatAPIKey    string `mapstructure:"chatApiKey"`
	BaseURL       string `mapstructure:"baseUrl"`
	Timeout       int    `mapstructure:"timeout"`
	DatasetID     string `mapstructure:"datasetId"`
}

// DeepSeekConfig DeepSeek配置
type DeepSeekConfig struct {
	APIKey      string  `mapstructure:"apiKey"`
	BaseURL     string  `mapstructure:"baseUrl"`
	Model       string  `mapstructure:"model"`
	Timeout     int     `mapstructure:"timeout"`
	MaxTokens   int     `mapstructure:"maxTokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// VolcEngineConfig 火山引擎AI配置 (豆包大模型)
type VolcEngineConfig struct {
	APIKey      string  `mapstructure:"apiKey"`
	BaseURL     string  `mapstructure:"baseUrl"`
	Model       string  `mapstructure:"model"` // 模型ID或接入点ID
	Timeout     int     `mapstructure:"timeout"`
	MaxTokens   int     `mapstructure:"maxTokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// BaiduOCRConfig 百度OCR配置
type BaiduOCRConfig struct {
	AppID     string `mapstructure:"appId"`
	APIKey    string `mapstructure:"apiKey"`
	SecretKey string `mapstructure:"secretKey"`
	Timeout   int    `mapstructure:"timeout"`
}

// VolcEngineImageConfig 火山引擎图片生成配置 (Doubao-Seedream)
// VolcEngine image generation config for CAD rendering
type VolcEngineImageConfig struct {
	APIKey  string `mapstructure:"apiKey"`
	BaseURL string `mapstructure:"baseUrl"`
	Model   string `mapstructure:"model"`   // Doubao-Seedream接入点ID
	Timeout int    `mapstructure:"timeout"` // 图片生成超时时间
	Size    string `mapstructure:"size"`    // 图片尺寸
}

// DWGConverterConfig DWG转换器配置
// DWG converter configuration for DWG to DXF conversion
type DWGConverterConfig struct {
	Type         string `mapstructure:"type"`         // 转换器类型: oda, libredwg / Converter type
	ODAPath      string `mapstructure:"odaPath"`      // ODA File Converter路径 / ODA path
	LibreDWGPath string `mapstructure:"libredwgPath"` // LibreDWG dwg2dxf路径 / LibreDWG path
}

// ODAConverterConfig 保留兼容旧配置
// Keep for backward compatibility
type ODAConverterConfig struct {
	Path string `mapstructure:"path"`
}

var GlobalConfig *Config

// LoadConfig 加载配置文件，支持环境变量覆盖
// Load config from file with environment variable override support
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 启用环境变量支持 / Enable environment variable support
	viper.AutomaticEnv()

	// 读取配置文件（如果存在）/ Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时使用默认值 / Use defaults if config file not found
		fmt.Printf("配置文件未找到，使用环境变量: %v\n", err)
	}

	// 绑定环境变量到配置项 / Bind environment variables to config keys
	bindEnvVariables()

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 从环境变量覆盖关键配置 / Override key configs from environment variables
	overrideFromEnv(&config)

	GlobalConfig = &config
	return &config, nil
}

// bindEnvVariables 绑定环境变量
// Bind environment variables to viper keys
func bindEnvVariables() {
	// 服务器配置 / Server config
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.mode", "GIN_MODE")

	// 数据库配置 / Database config
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.database", "DB_NAME")
	viper.BindEnv("database.username", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")

	// JWT配置 / JWT config
	viper.BindEnv("jwt.secret", "JWT_SECRET")

	// 火山引擎配置 / VolcEngine config
	viper.BindEnv("volcengine.apiKey", "VOLCENGINE_API_KEY")
	viper.BindEnv("volcengineImage.apiKey", "VOLCENGINE_IMAGE_API_KEY")
}

// overrideFromEnv 从环境变量覆盖配置
// Override config values from environment variables
func overrideFromEnv(config *Config) {
	if port := viper.GetInt("PORT"); port != 0 {
		config.Server.Port = port
	}
	if mode := viper.GetString("GIN_MODE"); mode != "" {
		config.Server.Mode = mode
	}
	if host := viper.GetString("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := viper.GetInt("DB_PORT"); port != 0 {
		config.Database.Port = port
	}
	if name := viper.GetString("DB_NAME"); name != "" {
		config.Database.Database = name
	}
	if user := viper.GetString("DB_USER"); user != "" {
		config.Database.Username = user
	}
	if pass := viper.GetString("DB_PASSWORD"); pass != "" {
		config.Database.Password = pass
	}
	if tls := viper.GetString("DB_TLS"); tls == "true" {
		config.Database.TLS = true
	}
	if secret := viper.GetString("JWT_SECRET"); secret != "" {
		config.JWT.Secret = secret
	}
	if apiKey := viper.GetString("VOLCENGINE_API_KEY"); apiKey != "" {
		config.VolcEngine.APIKey = apiKey
	}
	if apiKey := viper.GetString("VOLCENGINE_IMAGE_API_KEY"); apiKey != "" {
		config.VolcEngineImage.APIKey = apiKey
	}
}

// GetDSN 获取数据库连接字符串
// 支持 TiDB Cloud 等需要 TLS 的云数据库
func (c *DatabaseConfig) GetDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Charset,
		c.ParseTime,
	)
	// TiDB Cloud 需要 TLS 连接 / TiDB Cloud requires TLS
	if c.TLS {
		dsn += "&tls=true"
	}
	return dsn
}

// GetAccessTokenDuration 获取访问令牌有效期
func (c *JWTConfig) GetAccessTokenDuration() time.Duration {
	return time.Duration(c.AccessTokenExpire) * time.Second
}

// GetRefreshTokenDuration 获取刷新令牌有效期
func (c *JWTConfig) GetRefreshTokenDuration() time.Duration {
	return time.Duration(c.RefreshTokenExpire) * time.Second
}
