package polling

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/RuneHistory/collector/internal/application/domain"
	"github.com/RuneHistory/collector/internal/application/service"
	"github.com/RuneHistory/collector/internal/mapper"
	"net/http"
	"net/url"
	"time"
)

func NewHighScorePoller(accountService service.Account, bucketService service.Bucket, lookupHost string) *HighScorePoller {
	return &HighScorePoller{
		accountService: accountService,
		bucketService:  bucketService,
		lookupHost:     lookupHost,
		client: http.Client{
			Timeout: time.Second * 10,
		},
		highScoreCh: make(chan *domain.HighScore),
		errCh:       make(chan error),
	}
}

type HighScorePoller struct {
	accountService service.Account
	bucketService  service.Bucket
	lookupHost     string
	client         http.Client
	highScoreCh    chan *domain.HighScore
	errCh          chan error
}

func (p *HighScorePoller) HighScores() <-chan *domain.HighScore {
	return p.highScoreCh
}

func (p *HighScorePoller) Errors() <-chan error {
	return p.errCh
}

func (p *HighScorePoller) Poll(ctx context.Context) {
	buckets, err := p.bucketService.Get()
	if err != nil {
		p.errCh <- err
		return
	}
	for _, bucket := range buckets {
		err := p.pollBucket(bucket)
		if err != nil {
			p.errCh <- err
		}
	}
	close(p.highScoreCh)
	close(p.errCh)
}

func (p *HighScorePoller) pollBucket(bucket *domain.Bucket) error {
	accounts, err := p.accountService.GetByBucketId(bucket.ID)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		data, err := p.getHighScoreData(account.Nickname)
		if err != nil {
			p.errCh <- err
			continue
		}
		highScoreTransport := &mapper.HighScoreTransport{}
		err = json.Unmarshal([]byte(data), highScoreTransport)
		if err != nil {
			p.errCh <- err
			continue
		}
		highScore := mapper.HighScoreFromTransport(highScoreTransport)
		p.highScoreCh <- highScore
	}

	return nil
}

func (p *HighScorePoller) getHighScoreData(nickname string) (string, error) {
	lookupUrl := fmt.Sprintf("%s/osrs/highscore/%s", p.lookupHost, url.QueryEscape(nickname))
	req, err := http.NewRequest("GET", lookupUrl, nil)
	if err != nil {
		return "", err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBuffer := new(bytes.Buffer)
	_, err = bodyBuffer.ReadFrom(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return bodyBuffer.String(), fmt.Errorf("nickname not found on lookup: %s, %d", nickname, res.StatusCode)
	}
	return bodyBuffer.String(), nil
}
