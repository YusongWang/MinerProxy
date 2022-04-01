package cmd

import (
	"bufio"
	"log"
	"miner_proxy/handles/eth"
	"miner_proxy/network"
	"miner_proxy/serve"
	"net"

	"github.com/spf13/cobra"
)

func init() {
	// rootCmd.Flags().StringVarP(&cc.contentType, "kind", "k", "", "content type to create")
	// rootCmd.Flags().StringVar(&cc.contentEditor, "editor", "", "edit new content with this editor, if provided")

	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动MinerProxy核心，提供转发服务。",
	Long:  `无UI界面启动。`,
	Run: func(cmd *cobra.Command, args []string) {
		net, err := network.NewTcp(":38888")
		if err != nil {
			log.Panicln("can't bind to addr ", ":38888")
		}
		handle := eth.Handle{}
		for {
			//循环接入所有客户端得到专线连接
			conn, err := net.Accept()
			if err != nil {
				log.Panicln(err.Error())
				continue
			}
			log.Println("Tcp Accept Concent From ", conn.RemoteAddr())
			handle.OnConnect(conn.RemoteAddr().String())
			//开辟独立协程与该客聊天
			go start_loop(conn, handle)
		}

		s := serve.NewServe()
		s.StartLoop()
	},
}

func start_loop(conn net.Conn, handle eth.Handle) {
	/// TODO 封装到serve层
	reader := bufio.NewReader(conn)
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err.Error())
			handle.OnClose()
			return
		}

		ret, err := handle.OnMessage(conn, buf)
		if err != nil {
			log.Println(err.Error())
			handle.OnClose()
			return
		}

		_, err = conn.Write(ret)
		if err != nil {
			log.Println(err.Error())
			handle.OnClose()
			return
		}
	}
}
