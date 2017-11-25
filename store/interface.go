package store

import (
	"fmt"
	"net/url"

	"github.com/garyburd/redigo/redis"
	"github.com/i-tinerary/cotton/common"
)

type Interface interface {
	SetUser(common.User) error
	GetUser(string) (common.User, error)
	GetUsers() ([]string, error)

	SetPlace(common.Place) error
	GetPlace(int) (common.Place, error)

	SetPlan(common.Plan) error
	GetPlan(string, int) (common.Plan, error)
	GetPlans(string) ([]common.Plan, error)
}

func GetStore(url *url.URL) (Interface, error) {
	switch url.Scheme {
	case "redis":
		conn, err := redis.DialURL(url.String())
		if err != nil {
			return nil, fmt.Errorf("GetStore: %v", err)
		}
		return &impl{conn: conn}, nil
	}
	return nil, fmt.Errorf("unsupported url %q with scheme %q", url.String(), url.Scheme())
}
