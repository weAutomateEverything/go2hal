package database

import (
	"gopkg.in/mgo.v2"
	"log"
	"net"
	"crypto/tls"
	"net/url"
	"strings"
	"time"
	"strconv"
	"errors"
)

var database *mgo.Database

func init() {
	log.Println("Starting Database")

	mongo :=mongoConnectionString()

	var dialinfo *mgo.DialInfo

	if mongo == "" {
		dialinfo = getDialInfoParameters()
	} else {
		var err error
		dialinfo, err = parseMongoURL(mongo)
		if err != nil {
			log.Fatal(err)
		}
	}
	session, err := mgo.DialWithInfo(dialinfo)
	session.SetMode(mgo.Monotonic, true)

	database = session.DB(dialinfo.Database)
	if err != nil {
		log.Panic(err)
	}
}

func getDialInfoParameters() *mgo.DialInfo{
	dialinfo := mgo.DialInfo{}
	dialinfo.Addrs = mongoServers()
	dialinfo.Database = mongoDB()
	dialinfo.Password = mongoPassword()
	dialinfo.Username = mongoUser()
	dialinfo.ReplicaSetName = mongoReplicaSet()
	dialinfo.Source = mongoAuthSource()

	ssl := mongoSSL()

	if ssl {
		dialinfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		}
	}
	return &dialinfo
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
