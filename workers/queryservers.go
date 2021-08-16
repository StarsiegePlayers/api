package workers

import (
	"fmt"
	"os"
	"strings"
	"time"

	query "github.com/StarsiegePlayers/darkstar-query-go"
	"github.com/StarsiegePlayers/darkstar-query-go/master"
	"github.com/StarsiegePlayers/darkstar-query-go/protocol"
	"github.com/StarsiegePlayers/darkstar-query-go/server"
)

const ServerStatsPathFormat = "2006/01/02"

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
	data := performServerListUpdate()
	serverListData = data
	go recordServerListUpdate(data)
}

func performServerListUpdate() ServerListData {
	errors := make([]string, 0)
	masterQueryOptions := protocol.Options{
		Search: []string{
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

	serverQueryOptions := protocol.Options{
		Search:  gameAddresses,
		Timeout: 5 * time.Second,
	}
	games, errs := query.Servers(serverQueryOptions)
	for _, err := range errs {
		errors = append(errors, err.Error())
	}

	return ServerListData{
		RequestTime: time.Now(),
		Masters:     masterServerInfo,
		Games:       games,
		Errors:      errors,
	}
}

func recordServerListUpdate(data ServerListData) {
	var total, errorcount = 0, 0
	currentTime := time.Now()
	directory := fmt.Sprintf("./stats/%s", currentTime.Format(ServerStatsPathFormat))

	for _, s := range data.Games {
		err := os.MkdirAll(directory, 755)
		if err != nil {
			// app.Logger.Error(err)
			errorcount = errorcount + 1
			continue
		}

		fileName := strings.Replace(s.Address, ":", "_", 1)
		filePath := fmt.Sprintf("%s/%s.csv", directory, fileName)
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			// app.Logger.Error(err)
			errorcount = errorcount + 1
			continue
		}

		output := fmt.Sprintf("\"%s\",\"%s\",%s,%s,%d,%d,%s\n", currentTime.Format(time.RFC1123Z), s.Name, s.Ping, s.GameStatus, s.PlayerCount, s.MaxPlayers, s.Address)

		_, err = f.WriteString(output)
		if err != nil {
			// app.Logger.Error(err)
			errorcount = errorcount + 1
			_ = f.Close()
			continue
		}
		total = total + 1
		_ = f.Close()
	}
	app.Logger.Info(fmt.Sprintf("Wrote %d files to %s - %d file errors encountered", total, directory, errorcount))
}
