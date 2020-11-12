package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sceneryback/kafka-best-practices/scram"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sceneryback/kafka-best-practices/consumer"
	"github.com/sceneryback/kafka-best-practices/producer"
)

const (
	ProduceMode    = "produce"
	SyncMode       = "sync"
	BatchMode      = "batch"
	MultiAsyncMode = "multiAsync"
	MultiBatchMode = "multiBatch"
)

var mode string
var broker string
var kafkaUser string
var kafkaPassword string
var certFile string
var caFile string
var keyFile string
var topic string

func init() {
	flag.StringVar(&mode, "m", "", "cmd mode, 'produce', 'sync', 'batch' or 'multiBatch'")
	flag.StringVar(&broker, "h", "127.0.0.1:9092", "kafka broker host:port")

	flag.StringVar(&kafkaUser, "u", "user1", "kafka user")
	flag.StringVar(&kafkaPassword, "p", "user1", "kafka password")
	flag.StringVar(&certFile, "cert", "./jks/ca.crt", "crt file")
	flag.StringVar(&caFile, "ca", "./jks/ca.p12", "p12 file")
	flag.StringVar(&keyFile, "key", "./jks/ca.key", "key file")
	flag.StringVar(&topic, "t", "topic-1", "topic to use")
}

func main() {
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	if mode == "" {
		flag.Usage()
		return
	}

	var done = make(chan struct{})
	defer close(done)

	config := sarama.NewConfig()

	config.Version = sarama.V2_6_0_0

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.Retention = 10000000 * time.Minute
	//config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Fetch.Max = 100
	config.Consumer.MaxWaitTime = 100 * time.Millisecond

	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Idempotent = false
	config.Producer.RequiredAcks = sarama.WaitForAll
	//config.Producer.Retry.Max = 100
	config.Producer.Flush.MaxMessages = 100
	config.Producer.Flush.Frequency = time.Millisecond
	//config.Producer.Return.Successes = true
	//config.Producer.Return.Errors = true

	// SASL
	config.Net.SASL.Enable = true
	config.Net.SASL.User = kafkaUser
	config.Net.SASL.Password = kafkaPassword
	config.Net.SASL.Handshake = true
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &scram.XDGSCRAMClient{HashGeneratorFcn: scram.SHA512} }
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = createTLSConfiguration(certFile, caFile, keyFile)

	switch mode {
	case ProduceMode:
		p, err := producer.NewProducer(broker, config)
		if err != nil {
			panic(err)
		}
		defer p.Close()
		go p.StartProduce(done, topic)
	case SyncMode:
		// 1. sync cs
		cs, err := consumer.StartSyncConsumer(broker, topic, config)
		if err != nil {
			panic(err)
		}
		defer cs.Close()
	case BatchMode:
		// 2. batch cs
		cs, err := consumer.StartBatchConsumer(broker, topic, config)
		if err != nil {
			panic(err)
		}
		defer cs.Close()
	case MultiAsyncMode:
		// 3. multi async cs
		cs, err := consumer.StartMultiAsyncConsumer(broker, topic, config)
		if err != nil {
			panic(err)
		}
		defer cs.Close()
	case MultiBatchMode:
		// 4. multi batch cs
		cs, err := consumer.StartMultiBatchConsumer(broker, topic, config)
		if err != nil {
			panic(err)
		}
		defer cs.Close()
	default:
		flag.Usage()
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Println("received signal", <-c)
}

func createTLSConfiguration(certFile, caFile, keyFile string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		panic(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
}
