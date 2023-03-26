package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mileusna/useragent"
	"github.com/sirupsen/logrus"
	"go_ping_kube/internal/domain/models"
	"go_ping_kube/internal/infrastructure/repository"
	"net"
	"strings"
	"time"

	"github.com/oschwald/maxminddb-golang"
)

type PingService struct {
	storage repository.IPingRepository
}

func NewPingService(storage *repository.PingRedisStore) *PingService {
	return &PingService{storage: storage}
}

func (p *PingService) Add(ctx context.Context, create *models.CreatePingData) (*models.PingData, error) {
	ping := &models.PingData{
		Id:         uuid.New(),
		CreatedAt:  time.Now(),
		IP:         create.IP,
		RequestURI: create.RequestURI,
		UserAgent:  create.UserAgent,
		Headers:    create.Headers,
		Client:     "",
		Device:     "",
		OS:         "",
		Lang:       nil,
		Country:    "",
	}

	ua := useragent.Parse(create.UserAgent)
	logrus.Info(strings.Repeat("=", len(ua.String)))
	logrus.Info(ua.String)
	logrus.Info(strings.Repeat("=", len(ua.String)))

	ping.Client = fmt.Sprintf("%s v: %s", ua.Name, ua.Version)
	if ua.Mobile {
		ping.Device = "Mobile"
	}
	if ua.Tablet {
		ping.Device = "Tablet"
	}
	if ua.Desktop {
		ping.Device = "Desktop"
	}
	if ua.Bot {
		ping.Device = "Bot"
	}

	ping.OS = fmt.Sprintf("%s v: %s", ua.OS, ua.OSVersion)

	lang, ok := create.Headers["Accept-Language"]
	if ok {
		ping.Lang = lang
	}

	country, err := getLocationByMaxmind(ctx, create.IP)
	ping.Country = country.CountryCode

	err = p.storage.Save(ctx, ping)
	if err != nil {
		return nil, err
	}

	return ping, nil
}

func (p *PingService) Get(ctx context.Context, uuid uuid.UUID) (*models.PingData, error) {
	return p.storage.Get(ctx, uuid)
}

func (p *PingService) All(ctx context.Context) ([]*models.PingData, error) {
	return p.storage.All(ctx)
}

func getLocationByMaxmind(_ context.Context, userIP string) (*models.Location, error) {
	host, _, err := net.SplitHostPort(userIP)
	if err != nil {
		return nil, err
	}

	if host == "127.0.0.1" {
		return &models.Location{CountryCode: "TEST"}, nil
	}

	db, err := maxminddb.Open("GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	remoteIP := net.ParseIP(userIP)
	logrus.Info("remoteIP:", remoteIP)

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
		} `maxminddb:"country"`
	}

	err = db.Lookup(remoteIP, &record)
	if err != nil {
		return nil, err
	}
	logrus.Info("country from maxmind:", record.Country)

	return &models.Location{
		CountryCode: record.Country.ISOCode,
	}, nil
}
