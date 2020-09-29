package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"log"
	"mgr-agent/app"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	PID_FILE = "mgr-agent.pid"
	version  = "0.01"
)

var (
	daemon   bool
	confFile string
	workDir  string
)

func init() {
	absMgrAgent, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	workDir = filepath.Dir(absMgrAgent)
	os.Chdir(workDir)
}

func main() {
	var rootCmd = &cobra.Command{
		Use:     "mgr-agent",
		Short:   "mgr-agent \n\nA Virtual IP failover agent tool for MySQL Group Replication(Single Primary Mode): ",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.PersistentFlags().StringVarP(&confFile, "config", "c", "conf/agent.system", "configuration file to use")
	rootCmd.AddCommand(StartCmd())
	rootCmd.AddCommand(StopCmd())
	rootCmd.Execute()
}

func StartCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "start",
		Short: "start mgr-agent",
		Long:  "Use this to start mgr-agent",
		Run: func(cmd *cobra.Command, args []string) {
			if daemon {
				command := exec.Command("./bin/mgr-agent", "start")
				LockFile()
				command.Start()
				log.Printf("./bin/mgr-agent start, [PID] %d running...\n", command.Process.Pid)
				ioutil.WriteFile(PID_FILE, []byte(fmt.Sprintf("%d", command.Process.Pid)), 0666)
				daemon = false
				os.Exit(0)
			} else {
				log.Println("mgr-agent start")
			}
			_ = Run()
		},
	}
	cmd.Flags().BoolVarP(&daemon, "deamon", "d", false, "is daemon?")
	return
}

func Run() (err error) {
	if err = app.InitConfig(confFile); err != nil {
		return
	}
	if err = LogInit(); err != nil {
		return
	}
	if err = app.InitDbMeta(); err != nil {
		return
	}
	go ListenSignal()
	mysqlCheck := app.NewMysqlCheckHandler()
	mysqlCheck.CheckIsPrimaryLoop()
	return
}

func StopCmd() (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop mgr-agent",
		Run: func(cmd *cobra.Command, args []string) {
			strPid, _ := ioutil.ReadFile(PID_FILE)
			command := exec.Command("kill", string(strPid))
			if err := command.Run(); err == nil {
				logger.Info("mgr-agent stopped")
			}
			os.Remove(PID_FILE)
		},
	}
	return
}

func LogInit() (err error) {
	//logLevel := app.GetCfgString("system", "log_level" )
	if err = logger.SetLogger(`{
		"TimeFormat":"2006-01-02 15:04:05",
		"File": {                   
			"filename": "logs/mgr-agent.log",  
			"level": "TRAC",       
			"daily": true,          
			"maxlines": 100000000000,   
			"maxsize": 25,           
			"maxdays": 15,         
			"append": true,        
			"permit": "0660"       
	}
}`); err != nil {
		logger.Error("logger set Looger error desc=%v", err)
	}
	return
}

func LockFile() {
	lockFile := filepath.Join(workDir, PID_FILE)
	lock, err := os.Open(lockFile)
	defer lock.Close()
	if err == nil {
		filePid, err := ioutil.ReadAll(lock)
		if err == nil {
			pidStr := fmt.Sprintf("%s", filePid)
			pid, _ := strconv.Atoi(pidStr)
			_, err := os.FindProcess(pid)
			if err == nil {
				fmt.Printf("[ERROR] mgr agent might be already running...exiting[%s]\n", pidStr)
				os.Exit(1)
			}
		}
	}
}

func ListenSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	select {
	case <-sig:
		logger.Info("mgr-agent exit by signal %v:", sig)
		os.Remove(PID_FILE)
		os.Exit(0)
	}
}
