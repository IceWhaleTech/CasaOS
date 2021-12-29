package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

var ctx context.Context
var cancel context.CancelFunc

func TestNewInteruptReader(t *testing.T) {
	ctx, cancel = context.WithCancel(context.Background())

	go func() {
		// 在初始上下文的基础上创建一个有取消功能的上下文
		//	ctx, cancel := context.WithCancel(ctx)
		fmt.Println("开始")
		fIn, err := os.Open("/Users/liangjianli/Downloads/demo_data.tar.gz")
		if err != nil {

		}
		defer fIn.Close()
		fmt.Println("创建新文件")
		fOut, err := os.Create("/Users/liangjianli/Downloads/demo_data1.tar.gz")
		if err != nil {
			fmt.Println(err)
		}

		defer fOut.Close()

		fmt.Println("准备复制")
		//	_, err = io.Copy(out, NewReader(ctx, f))
		//	time.Sleep(time.Second * 2)
		//ctx.Done()
		//	cancel()

		// interrupt context after 500ms

		// interrupt context with SIGTERM (CTRL+C)
		//sigs := make(chan os.Signal, 1)
		//signal.Notify(sigs, os.Interrupt)

		if err != nil {
			log.Fatal(err)
		}

		// Reader that fails when context is canceled
		in := NewReader(ctx, fIn)
		// Writer that fails when context is canceled
		out := NewWriter(ctx, fOut)

		//time.Sleep(2 * time.Second)

		//cancel()

		n, err := io.Copy(out, in)
		log.Println(n, "bytes copied.")
		if err != nil {
			fmt.Println("Err:", err)
		}

		fmt.Println("Closing.")
	}()

	go func() {
		//<-sigs
		time.Sleep(time.Second)
		fmt.Println("退出")
		ddd()
	}()
	time.Sleep(time.Second * 10)
}

func ddd() {
	cancel()
}
