package app

import (
	"fmt"
	"github.com/setcreed/store-core/pkg/util/log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	v1 "github.com/setcreed/store-core/api/store_service/v1"
	"github.com/setcreed/store-core/cmd/app/options"
)

func NewServerCommand(version string) *cobra.Command {
	opts := options.NewOptions()
	cmd := &cobra.Command{
		Use: "db-server",
		Run: func(cmd *cobra.Command, args []string) {

			if err := opts.Validate(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := opts.Complete(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := opts.InitHttp(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}

			if err := Run(opts); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}
	// 绑定参数
	opts.BindFlags(cmd)

	verCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  "Print version and exit.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	reloadCmd := &cobra.Command{
		Use:   "reload",
		Short: "reload sysconfig(app.yaml)",
		Run: func(cmd *cobra.Command, args []string) {
			err := opts.InitConfig()
			if err != nil {
				fmt.Println(err)
				return
			}
			rsp, err := http.Get(fmt.Sprintf("http://localhost:%d/reload", opts.ComponentConfig.Default.App.HttpPort))
			if err != nil {
				log.Infoln(err)
			}
			defer rsp.Body.Close()
			if rsp.StatusCode == 200 {
				log.Infoln("配置文件重载成功")
			} else {
				log.Errorln("配置文件重载失败")
			}
		},
	}

	opts.BindFlags(reloadCmd)

	cmd.AddCommand(verCmd, reloadCmd)

	return cmd
}

func Run(opts *options.Options) error {
	server := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	// 注册服务
	v1.RegisterDBServiceServer(server, opts.Store)
	// 创建一个通道来监听系统中断信号
	sigChan := make(chan os.Signal, 1)
	// 监听 SIGINT 和 SIGTERM 信号
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 使用 goroutine 来启动服务器
	go func() {
		fmt.Printf("启动对外grpc服务,端口是:%v\n", opts.ComponentConfig.Default.App.RpcPort)
		if err = server.Serve(lis); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}()

	// 阻塞等待接收到信号
	sig := <-sigChan
	fmt.Printf("Received signal %s, initiating graceful shutdown\n", sig)

	// 接收到信号后优雅地停止服务器
	server.GracefulStop()
	fmt.Println("Server has shut down")

	return nil
}
