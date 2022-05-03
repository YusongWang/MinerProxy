package test

import (
	"log"
	"miner_proxy/pack"
	"miner_proxy/utils"
	"strings"
	"testing"

	"github.com/buger/jsonparser"
)

var send_job []string
var feejob = []byte(`{"id":0,"version":"2.0","result":["0xf6067c77e474565f43b8c6b22afe7c12a178ec400e95974be5bad4f92657969d", "0xf87998fd030a4d04802b9f2ef04443bb4d2c105f13be5aa338ef49860c3c5425", "0x000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffff"]}`)

func BenchmarkSendJob(b *testing.B) {
	send_job := &pack.Job{}
	a := []string{"0xf6067c77e474565f43b8c6b22afe7c12a178ec400e95974be5bad4f92657969d", "0xf87998fd030a4d04802b9f2ef04443bb4d2c105f13be5aa338ef49860c3c5425", "0x000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffff"}
	send_job.Job = append(send_job.Job, a)
	worker := pack.NewWorker("1", "2", "3", "4")

	for i := 0; i < b.N; i++ {
		worker.FeeAdd()
		//Repeat(`b`, 5)
		FeeSend(send_job.Job, worker)
	}
}

func BenchmarkSendJobButter(b *testing.B) {
	//send_job := &pack.Job{}
	var Job [][]byte
	//a := []string{"0xf6067c77e474565f43b8c6b22afe7c12a178ec400e95974be5bad4f92657969d", "0xf87998fd030a4d04802b9f2ef04443bb4d2c105f13be5aa338ef49860c3c5425", "0x000000007fffffffffffffffffffffffffffffffffffffffffffffffffffffff"}
	Job = append(Job, feejob)
	worker := pack.NewWorker("1", "2", "3", "4")

	for i := 0; i < b.N; i++ {
		worker.FeeAdd()
		//Repeat(`b`, 5)
		FeeSendBetter(Job, worker)
	}
}

func FeeSend(job [][]string, worker *pack.Worker) {

	if len(job) > 0 {
		send_job = job[len(job)-1]
	} else {
		return
	}

	if len(send_job) == 0 {
		log.Println("当前job内容为空")
		return
	}

	diff := utils.TargetHexToDiff(send_job[2])
	worker.SetDevDiff(diff)

	//dev[send_job[0]] = true
	job_str := ConcatJobTostr(send_job)
	_ = ConcatToPushJob(job_str)

	// _, err := c.Write(job_byte)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return
	// }
}

var Dev = make(map[string]bool)
var job_id string

func FeeSendBetter(job [][]byte, worker *pack.Worker) {
	var send []byte
	if len(job) > 0 {
		send = job[len(job)-1]
	} else {
		return
	}

	if worker.Fee_idx == 1 {
		// diff := utils.TargetHexToDiff(send_job[2])
		// worker.SetDevDiff(diff)
	} else if worker.Fee_idx%100 == 0 {
		// diff := utils.TargetHexToDiff(send_job[2])
		// worker.SetDevDiff(diff)
	}

	job_id, err := jsonparser.GetString(send, "result", "[0]")
	if err != nil {
		log.Println("Err")
		return
	}

	Dev[job_id] = true
	//job_str := []byte(ConcatJobTostr(send_job))
	// _ = ConcatToPushJob(job_str)
	// a := append(golbal_jobb, job_str...)
	// _ = append(a, golbal_jobendb...)
	// _, err := c.Write(job_byte)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return
	// }
}

var golbal_job = `{"id":0,"jsonrpc":"2.0","result":`
var golbal_jobend = `}\n`
var golbal_jobb = []byte(`{"id":0,"jsonrpc":"2.0","result":`)
var golbal_jobendb = []byte(`}\n`)

func ConcatJobTostr(job []string) string {
	var builder strings.Builder
	builder.WriteString(`["`)

	job_len := len(job) - 1
	for i, j := range job {
		if i == job_len {
			builder.WriteString(j + `"]`)
			break
		}
		builder.WriteString(j + `","`)
	}

	return builder.String()
}

func ConcatToPushJob(job string) []byte {
	//inner_job := []byte(golbal_job + string(job) + golbal_jobend)
	var builder strings.Builder
	builder.WriteString(golbal_job)
	builder.WriteString(job)
	builder.WriteString(golbal_jobend)
	builder.WriteByte('\n')
	return []byte(builder.String())
}
