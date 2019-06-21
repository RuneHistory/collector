package main

import (
	"collector/internal/application/handler"
	"collector/internal/events"
	"collector/internal/migrate"
	"collector/internal/migrate/migrations"
	"database/sql"
	"errors"
	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	dsn := os.Getenv("DSN")
	dsn = dsn + "?multiStatements=true&parseTime=true"

	brokers := os.Getenv("KAFKA_BROKERS")
	brokerList := strings.Split(brokers, ",")

	wg := &sync.WaitGroup{}
	shutdownCh := make(chan struct{})
	errCh := make(chan error)
	go handleShutdownSignal(errCh)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	err = migrateDb(db)
	if err != nil {
		panic(err)
	}

	log.Println("creating kafka client")
	kafkaClient, err := createKafkaClient(brokerList)
	if err != nil {
		panic(err)
	}
	var subscriber events.Subscriber = events.NewKafkaSubscriber(brokerList, kafkaClient, "rh-collector")
	//accountRepo := mysql.NewAccountMySQL(db)

	// TODO: Make this better...

	eventHandlers := map[events.Event][]events.Handler{
		events.NewAccountEvent: {
			&handler.NewAccountEventHandler{},
		},
	}
	for event, handlers := range eventHandlers {
		for _, h := range handlers {
			log.Printf("subscribing %T to %T\n", h, event)
			err = subscriber.Subscribe(event, h)
			if err != nil {
				log.Printf("fatal err: %s\n", err)
			}
		}
	}

	// Start consumers/producers here
	//go http_transport.Start(address, r, wg, shutdownCh, errCh)

	err = <-errCh
	if err != nil {
		log.Printf("fatal err: %s\n", err)
	}

	log.Println("initiating graceful shutdown")
	close(shutdownCh)

	wg.Wait()
	log.Println("shutdown")
}

func handleShutdownSignal(errCh chan error) {
	quitCh := make(chan os.Signal)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGTERM)

	hit := false
	for {
		<-quitCh
		if hit {
			os.Exit(0)
		}
		if !hit {
			errCh <- errors.New("shutdown signal received")
		}
		hit = true
	}
}

func migrateDb(db *sql.DB) error {
	migrationArr := []migrate.Migration{
		&migrations.CreateBucketsTable{},
		&migrations.CreateAccountsTable{},
	}
	return migrate.Migrate(db, migrationArr)
}

func createKafkaClient(brokerList []string) (sarama.Client, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Timeout = time.Second * 3
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return sarama.NewClient(brokerList, config)
}
