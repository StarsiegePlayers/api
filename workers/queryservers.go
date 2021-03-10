package workers

import (
	query "github.com/StarsiegePlayers/darkstar-query-go"
	"time"

	"github.com/StarsiegePlayers/darkstar-query-go/master"
	"github.com/StarsiegePlayers/darkstar-query-go/server"
)

type ServerListData struct {
	RequestTime time.Time
	Masters     []*master.Master
	Games       []*server.PingInfo
	Errors      []string
}

var serverListData ServerListData

func GetServerList() ServerListData {
	return serverListData
}

func UpdateServerList() {
	app.Logger.Info("Crontab - Updating server list")
	performServerListUpdate()
}

func performServerListUpdate() {
	errors := make([]string, 0)
	masterQueryOptions := query.Options{
		Search: []string {
			"master1.starsiegeplayers.com:29000",
			"master2.starsiegeplayers.com:29000",
			"master3.starsiegeplayers.com:29000",
			"starsiege1.no-ip.org:29000",
			"starsiege.noip.us:29000",
			"southerjustice.dyndns-server.com:29000",
			"dustersteve.ddns.net:29000",
			"starsiege.from-tx.com:29000",
		},
		Timeout: 5 * time.Second,
	}

	masterServerInfo, gameAddresses, errs := query.Masters(masterQueryOptions)
	if len(errs) >= 0 {
		for _, v := range errs {
			errors = append(errors, v.Error())
		}
	}

	serverQueryOptions := query.Options{
		Search:  gameAddresses,
		Timeout: 5 * time.Second,
	}
	games, errs := query.Servers(serverQueryOptions)
	for _, err := range errs {
		errors = append(errors, err.Error())
	}

	serverListData = ServerListData {
		RequestTime: time.Now(),
		Masters:     masterServerInfo,
		Games:       games,
		Errors:      errors,
	}
}
