package options

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/setcreed/store-core/cmd/app/config"
	"github.com/setcreed/store-core/pkg/controller"
	"github.com/setcreed/store-core/pkg/data"
	"github.com/setcreed/store-core/pkg/util/cfg"
)

const (
	defaultListen = 8090

	maxIdleConns = 10
	maxOpenConns = 100
	maxLifeTime  = 1800

	defaultConfigFile = "/etc/dbcore/config.yaml"
)

type Options struct {
	ComponentConfig *config.Config
	// 数据库接口
	Factory data.ShareDaoFactory

	Store *controller.DBService

	ConfigFile string
}

func NewOptions() *Options {
	return &Options{
		ConfigFile: defaultConfigFile,
	}
}

func (o *Options) Complete() error {
	if len(o.ConfigFile) == 0 {
		if cfgFile := os.Getenv("ConfigFile"); cfgFile != "" {
			o.ConfigFile = cfgFile
		} else {
			o.ConfigFile = defaultConfigFile
		}
	}

	// 加载配置
	if err := o.InitConfig(); err != nil {
		return err
	}

	if o.ComponentConfig.Default.App.RpcPort == 0 {
		o.ComponentConfig.Default.App.RpcPort = 8090
	}

	// 注册依赖组件
	if err := o.register(); err != nil {
		return err
	}

	// 注册controller依赖
	o.Store = controller.NewDBService(o.ComponentConfig, o.Factory)
	return nil
}

func (o *Options) BindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.ConfigFile, "configfile", "", "The location of the dbcore configuration file")
}

func (o *Options) register() error {
	if err := o.registerDatabase(); err != nil {
		return err
	}
	return nil
}

func (o *Options) registerDatabase() error {
	dsn := o.ComponentConfig.DBConfig.DSN
	opt := &gorm.Config{}
	if o.ComponentConfig.Default.Mode == "debug" {
		opt.Logger = logger.Default.LogMode(logger.Info)
	}

	DB, err := gorm.Open(mysql.Open(dsn), opt)
	if err != nil {
		return err
	}
	// 设置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(o.ComponentConfig.DBConfig.MaxOpenConn)
	sqlDB.SetMaxIdleConns(o.ComponentConfig.DBConfig.MinIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(o.ComponentConfig.DBConfig.MaxLifeSecond) * time.Second)

	o.Factory, err = data.NewShareDaoFactory(DB)
	if err != nil {
		return err
	}

	return nil
}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) InitConfig() error {
	c := cfg.New()
	c.SetConfigType("yaml")
	c.SetConfigFile(o.ConfigFile)
	if err := c.Binding(&o.ComponentConfig); err != nil {
		return err
	}
	return nil
}

func (o *Options) InitHttp() error {
	go func() {
		http.HandleFunc("/reload", func(writer http.ResponseWriter, request *http.Request) {
			err := o.InitConfig()
			if err != nil {
				writer.Write([]byte(fmt.Sprintf("reload config error:%s", err.Error())))
			}
		})
		fmt.Println(fmt.Sprintf("启动内部http服务,端口是:%d", o.ComponentConfig.Default.App.HttpPort))
		http.ListenAndServe(fmt.Sprintf(":%d", o.ComponentConfig.Default.App.HttpPort), nil)
	}()
	return nil
}
