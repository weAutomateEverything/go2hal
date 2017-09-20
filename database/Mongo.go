package database

import (
	"gopkg.in/mgo.v2"
	"log"
	"github.com/zamedic/go2hal/config"
	"strings"
	"net/url"
	"time"
	"strconv"
	"errors"
	"net"
	"crypto/tls"
)

var database *mgo.Database

func init() {
	log.Println("Starting Database")

	dialinfo, err := parseMongoURL(config.MongoAddress())
	log.Println("Connecting to:  "+dialinfo.Database)
	if err != nil {
		log.Panic(err)
	}


	session, err := mgo.DialWithInfo(dialinfo)
	database = session.DB("hal")
	if err != nil {
		log.Panic(err)
	}
}

func parseMongoURL(rawURL string) (*mgo.DialInfo, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	info := mgo.DialInfo{
		Addrs:    strings.Split(url.Host, ","),
		Database: strings.TrimPrefix(url.Path, "/"),
		Timeout:  10 * time.Second,
	}

	if url.User != nil {
		info.Username = url.User.Username()
		info.Password, _ = url.User.Password()
	}

	query := url.Query()
	for key, values := range query {
		var value string
		if len(values) > 0 {
			value = values[0]
		}

		switch key {
		case "authSource":
			info.Source = value
		case "authMechanism":
			info.Mechanism = value
		case "gssapiServiceName":
			info.Service = value
		case "replicaSet":
			info.ReplicaSetName = value
		case "maxPoolSize":
			poolLimit, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("bad value for maxPoolSize: " + value)
			}
			info.PoolLimit = poolLimit
		case "ssl":
			// Unfortunately, mgo doesn't support the ssl parameter in its MongoDB URI parsing logic, so we have to handle that
			// ourselves. See https://github.com/go-mgo/mgo/issues/84
			ssl, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.New("bad value for ssl: " + value)
			}
			if ssl {
				info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
					return tls.Dial("tcp", addr.String(), &tls.Config{})
				}
			}
		case "connect":
			if value == "direct" {
				info.Direct = true
				break
			}
			if value == "replicaSet" {
				break
			}
			fallthrough
		default:
			return nil, errors.New("unsupported connection URL option: " + key + "=" + value)
		}
	}

	return &info, nil
}
