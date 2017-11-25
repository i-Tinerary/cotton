package store

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/i-tinerary/cotton/common"
)

type impl struct {
	conn redis.Conn
}

func (i *impl) SetUser(u common.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("SetUser: marshal: %v", err)
	}
	_, err = i.conn.Do("SADD", "users", u.Name)
	if err != nil {
		return fmt.Errorf("SetUser: redis1: %v", err)
	}
	_, err = i.conn.Do("SET", "user:"+u.Name, data)
	if err != nil {
		return fmt.Errorf("SetUser: redis2: %v", err)
	}

	return nil
}

func (i *impl) GetUser(name string) (common.User, error) {
	data, err := redis.Bytes(i.conn.Do("GET", "user:"+name))
	if err != nil {
		return common.User{}, fmt.Errorf("GetUser: redis: %v", err)
	}
	var user common.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		return common.User{}, fmt.Errorf("GetUser: unmarshal: %v", err)
	}
	return user, nil
}

func (i *impl) GetUsers() ([]string, error) {
	users, err := redis.Strings(i.conn.Do("SMEMBERS", "users"))
	if err != nil {
		return nil, fmt.Errorf("GetUsers: redis: %v", err)
	}
	return users, nil
}

func (i *impl) SetPlace(place common.Place) error {
	data, err := json.Marshal(place)
	if err != nil {
		return fmt.Errorf("SetPlace: marshal: %v", err)
	}
	id, err := redis.Int(i.conn.Do("INCR", "place_id"))
	if err != nil {
		return fmt.Errorf("SetPlace: create id: %v", err)
	}
	_, err = i.conn.Do("SET", "place:"+strconv.Itoa(id), data)
	if err != nil {
		return fmt.Errorf("SetPlace: redis: %v", err)
	}
	return nil
}

func (i *impl) GetPlace(id int) (common.Place, error) {
	data, err := redis.Bytes(i.conn.Do("GET", "place:"+strconv.Itoa(id)))
	if err != nil {
		return common.Place{}, fmt.Errorf("GetPlace: redis: %v", err)
	}
	var place common.Place
	err = json.Unmarshal(data, &place)
	if err != nil {
		return common.Place{}, fmt.Errorf("GetPlace: unmarshal: %v", err)
	}
	return place, nil
}

func (i *impl) SetPlan(common.Plan) error                { return nil }
func (i *impl) GetPlan(string, int) (common.Plan, error) { return common.Plan{}, nil }
func (i *impl) GetPlans(string) ([]common.Plan, error)   { return nil, nil }
