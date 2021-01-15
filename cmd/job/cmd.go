package job

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/RichardKnop/machinery/v1/tasks"
	jobLog "github.com/RichardKnop/machinery/v1/log"
	// tracers "github.com/RichardKnop/machinery/example/tracers"

	"abs/pkg/conf"
	"abs/pkg/job"
	jobTasks "abs/internal/job/tasks"
)

var env string
var queue string

// Cmd run job once or periodically
var Cmd = &cobra.Command{
	Use:   "job",
	Short: "Run job",
	Long:  "运行异步队列任务进程",
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化各项服务
		initStep()
		// Register tasks
		taskLists := map[string]interface{}{
			"add":               jobTasks.Add,
			"multiply":          jobTasks.Multiply,
			"sum_ints":          jobTasks.SumInts,
			"sum_floats":        jobTasks.SumFloats,
			"concat":            jobTasks.Concat,
			"split":             jobTasks.Split,
			"panic_task":        jobTasks.PanicTask,
			"long_running_task": jobTasks.LongRunningTask,
		}
		// 注册任务函数
		if err := job.Machinery.RegisterTasks(taskLists); err != nil {
			log.Fatalf("Job Machinery Register tasks err: %v", err)
		}

		// The second argument is a consumer tag
		// Ideally, each worker should have a unique tag (worker1, worker2 etc)
		worker := job.Machinery.NewWorker(queue, 0)

		// Here we inject some custom code for error handling,
		// start and end of task hooks, useful for metrics for example.
		errorhandler := func(err error) {
			jobLog.ERROR.Println("I am an error handler:", err)
		}

		pretaskhandler := func(signature *tasks.Signature) {
			jobLog.INFO.Println("I am a start of task handler for:", signature.Name)
		}

		posttaskhandler := func(signature *tasks.Signature) {
			jobLog.INFO.Println("I am an end of task handler for:", signature.Name)
		}

		worker.SetPostTaskHandler(posttaskhandler)
		worker.SetErrorHandler(errorhandler)
		worker.SetPreTaskHandler(pretaskhandler)

		// 运行队列消费服务
		if err := worker.Launch(); err != nil {
			log.Fatalf("Job Machinery Worker Run err: %v", err)
		}
	},
}

// init server cmd
func init() {
	Cmd.Flags().StringVar(&env, "env", "local", "conf environmental science")
	Cmd.Flags().StringVar(&queue, "queue", "abs_machinery_tasks", "job default queue")
}

// init job cmd
func initStep() {
	// 初始化各项服务
	// 配置加载
	conf.Init(env)
	// 初始化队列服务
	job.MachineryStartServer(queue)
}