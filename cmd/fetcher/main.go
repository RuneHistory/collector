package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/RuneHistory/collector/internal/application/polling"
	"github.com/RuneHistory/collector/internal/application/service"
	"github.com/RuneHistory/collector/internal/application/validate"
	"github.com/RuneHistory/collector/internal/repository/mysql"
	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	dsn := os.Getenv("DSN")
	dsn = dsn + "?multiStatements=true&parseTime=true"

	lookupHost := os.Getenv("LOOKUP_HOST")
	if lookupHost == "" {
		panic(fmt.Errorf("LOOKUP_HOST not specified"))
	}
	//brokers := os.Getenv("KAFKA_BROKERS")
	//brokerList := strings.Split(brokers, ",")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	shutdownCh := make(chan struct{})
	errCh := make(chan error)
	go handleShutdownSignal(shutdownCh)
	go func() {
		select {
		case <-shutdownCh:
			break
		case err := <-errCh:
			log.Printf("fatal error: %s", err)
		}
		cancel()
	}()

	bucketRepo := mysql.NewBucketMySQL(db)
	bucketRules := validate.NewBucketRules(bucketRepo)
	bucketValidator := validate.NewBucketValidator(bucketRules)
	bucketService := service.NewBucketService(bucketRepo, bucketValidator)

	accountRepo := mysql.NewAccountMySQL(db)
	accountRules := validate.NewAccountRules(accountRepo, bucketRepo)
	accountValidator := validate.NewAccountValidator(accountRules)
	accountService := service.NewAccountService(accountRepo, accountValidator)

	poller := polling.NewHighScorePoller(accountService, bucketService, lookupHost)

	wg.Add(1)
	go func() {
		for err := range poller.Errors() {
			log.Printf("got err: %s", err)
		}
		log.Println("ended err worker")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for highScore := range poller.HighScores() {
			log.Printf("got highscore: %v", *highScore)
		}
		log.Println("ended highscore worker")
		wg.Done()
	}()

	poller.Poll(ctx)

	// doneCh will be closed once wg is done
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		// we're finished so start the shutdown
		log.Println("all services finished")
	case <-ctx.Done():
		break
		// break out and wait for shutdown
	}

	log.Println("waiting for shutdown")

	select {
	case <-time.After(time.Second * 10):
		log.Println("killed - took too long to shutdown")
	case <-doneCh:
		log.Println("all services shutdown")
	}
}

func handleShutdownSignal(shutdownCh chan struct{}) {
	quitCh := make(chan os.Signal)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGTERM)

	startedShutdown := false
	for {
		<-quitCh
		if startedShutdown {
			os.Exit(0)
		}
		close(shutdownCh)
		startedShutdown = true
	}
}

func createSaramaClient(brokerList []string) (sarama.Client, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.Producer.Timeout = time.Second
	config.Version = sarama.V1_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewClient(brokerList, config)
}
