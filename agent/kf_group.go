package agent

import (
	"sync"

	"fmt"

	"github.com/barockok/camelia/misc"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ErrorKFGroupSetupConsumer struct {
	originError error
	kfConf      map[string]interface{}
	defConsConf map[string]interface{}
}

func (err *ErrorKFGroupSetupConsumer) Error() string {
	return fmt.Sprintf(`ErrorKFGroupSetupConsumer
	kafka-conf : %v
	consumer-default : %v
	`, err.kfConf, err.defConsConf)
}

type KFGroup struct {
	Identifier    string
	MemberCount   int
	State         string
	agentList     []*KFAgent
	wg            sync.WaitGroup
	Topics        []string
	MsgErrHandler MessageErrorHandler
	MsgHandler    MessageHandler
}

func (g *KFGroup) Init(kfConf, defConsConf map[string]interface{}) {

}
func (g *KFGroup) Run() {
	var wg sync.WaitGroup
	for i := 0; i < g.MemberCount; i++ {
		wg.Add(1)
		var ag *KFAgent
		var agentID string
		var consumer *kafka.Consumer
		agentID = fmt.Sprintf("KF-AGENT-%s-%d", g.Identifier, i)
		ag = NewKFAgent(agentID, g.Topics, consumer, g.MsgHandler, g.MsgErrHandler)
		g.agentList = append(g.agentList, ag)

		go func(cWg *sync.WaitGroup, cAg *KFAgent) {
			cAg.Start()
			cWg.Done()
		}(&wg, ag)
	}
	wg.Wait()
}

func setupKFConsumer(kfConf, defConsConf map[string]interface{}) (*kafka.Consumer, error) {
	convKfConf := misc.ToI(kfConf).(kafka.ConfigMap)
	convKfConf["go.events.channel.enable"] = true
	convKfConf["go.application.rebalance.enable"] = true
	convKfConf["default.topic.config"] = misc.ToI(defConsConf).(kafka.ConfigMap)
	return kafka.NewConsumer(&convKfConf)
}
