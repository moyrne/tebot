package commands

import (
	"github.com/moyrne/tebot/internal/logs"
	"github.com/moyrne/tractor/syncx"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"sync/atomic"
	"time"
)

var netNormal = int32(1)

// TODO 抽象出 看门狗功能

func StartCQHTTP() {
	syncx.Go(netChecking)
	for {
		syncx.Safe(func() {
			// 拉起cqhttp进程
			if err := startCQHTTP(); err != nil {
				logs.Error("start cqhttp", "error", err)
				time.Sleep(time.Second)
			}
			// 当进程退出时, 检测网络状态, 判断是否重新拉起
			for atomic.LoadInt32(&netNormal) != 1 {
				// 自旋等待网络通畅
				time.Sleep(time.Second)
			}
		})
	}
}

func netChecking() {
	window := [3]bool{true, true, true}
	tick := time.NewTicker(time.Second * 5)
	for {
		if !(window[0] && window[1] && window[2]) {
			atomic.StoreInt32(&netNormal, 0)
		}
		if window[0] && window[1] && window[2] {
			atomic.StoreInt32(&netNormal, 1)
		}

		window[0], window[1], window[2] = window[1], window[2], ping()
		<-tick.C
	}
}

func ping() bool {
	cmd := exec.Command("ping", "www.baidu.com", "-c", "1", "-W", "5")
	out, err := cmd.Output()
	if err != nil {
		logs.Error("ping", "error", err, "msg", string(out))
		return false
	}
	return true
}

func startCQHTTP() error {
	cqCmd := exec.Command("cqhttp")
	f, err := os.OpenFile("cqhttp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	w := io.MultiWriter(f, os.Stdout)
	cqCmd.Stdout = w
	cqCmd.Stderr = w
	done := make(chan struct{})
	syncx.Go(func() {
		cqHeartbeat(cqCmd, done)
	})
	defer func() {
		close(done)
	}()
	return errors.WithStack(cqCmd.Run())
}

var heartbeat = make(chan struct{}, 1)

func CQHeartBeat() {
	heartbeat <- struct{}{}
}

func cqHeartbeat(cqCmd *exec.Cmd, done chan struct{}) {
	du := time.Second * 20
	timer := time.NewTimer(du)
	window := [3]bool{true, true, true}
	for {
		timer.Reset(du)
		select {
		case <-done:
			// 主动退出
			return
		case <-heartbeat:
			window[0], window[1], window[2] = true, true, true
		case <-timer.C:
			logs.Error("cqhttp heartbeat overtime")
			// 心跳出错
			window[0], window[1], window[2] = window[1], window[2], false
			if !window[0] && !window[1] && !window[2] {
				killProcess(cqCmd)
				return
			}
		}
	}
}

func killProcess(cqCmd *exec.Cmd) {
	if cqCmd != nil && cqCmd.Process != nil {
		if err := cqCmd.Process.Kill(); err != nil {
			logs.Error("cqhttp kill", "error", err)
			return
		}
	}
}
