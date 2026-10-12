package casbin

import (
	"backend-rms/core/config"
	"backend-rms/core/database"
	"encoding/json"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	CASBIN_P_PREFIX = "role_"
	CASBIN_G_PREFIX = "manager_"
)

var enforcer *casbin.Enforcer
var watcher persist.Watcher

func InitCasbin() {
	watcher = InitCasbinWatcher()

	a, newAdapterErr := gormadapter.NewAdapterByDB(database.RMSDB)
	if newAdapterErr != nil {
		panic("NewAdapterErr, error=" + newAdapterErr.Error())
	}

	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", `g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act || r.sub == "role_1"`)

	e, newEnforcerErr := casbin.NewEnforcer(m, a)
	if newEnforcerErr != nil {
		panic("NewAdapterByDB, error=" + newEnforcerErr.Error())
	}

	err := e.SetWatcher(watcher)
	if err != nil {
		panic(err.Error())
	}

	err = watcher.SetUpdateCallback(rediswatcher.DefaultUpdateCallback(e))
	if err != nil {
		panic(err.Error())
	}

	loadPolicyErr := e.LoadPolicy()
	if loadPolicyErr != nil {
		panic("loadPolicy, error=" + loadPolicyErr.Error())
	}

	enforcer = e
}

func InitCasbinWatcher() persist.Watcher {
	w, err := rediswatcher.NewWatcher(config.Conf.Redis.RMS.Host, rediswatcher.WatcherOptions{
		Options: redis.Options{
			Network:  "tcp",
			Password: config.Conf.Redis.RMS.Password,
			DB:       config.Conf.Redis.RMS.DB,
		},
		Channel: "/casbin",
		// Only exists in test, generally be true
		IgnoreSelf: true,
	})
	if err != nil {
		panic("NewWatcher, error=" + err.Error())
	}
	return w
}

func EnforceByRole(role, uri, act string) (bool, error) {
	return enforcer.Enforce(role, uri, act)
}

func EnforceByManager(manager, uri, act string) (bool, error) {
	return enforcer.Enforce(manager, uri, act)
}

func AddPolicy(pType string, rule []string) error {
	if pType == "p" {
		ok, err := enforcer.AddPolicy(rule)
		if err != nil {
			return errors.Wrap(err, "casbin: add p type policy fail")
		}

		if !ok {
			return errors.New("add p type casbin policy fail")
		}

		err = watcher.Update()
		if err != nil {
			return errors.Wrap(err, "casbin: watcher update failed")
		}
		return nil
	}

	ok, err := enforcer.AddGroupingPolicy(rule)
	if err != nil {
		return errors.Wrap(err, "casbin: add g type policy fail")
	}

	if !ok {
		return errors.New("add g type casbin policy fail")
	}

	err = watcher.Update()
	if err != nil {
		return errors.Wrap(err, "casbin: watcher update failed")
	}
	return nil
}

func AddPolicies(pType string, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}
	if pType == "p" {
		ok, err := enforcer.AddPolicies(rules)
		if err != nil {
			return errors.Wrap(err, "casbin: add p type policies fail")
		}

		if !ok {
			return errors.New("add p type casbin policies fail")
		}

		err = watcher.Update()
		if err != nil {
			return errors.Wrap(err, "casbin: watcher update failed")
		}
		return nil
	}

	ok, err := enforcer.AddGroupingPolicies(rules)
	if err != nil {
		return errors.Wrap(err, "casbin: add g type policies fail")
	}

	if !ok {
		return errors.New("add g type casbin policies fail")
	}

	err = watcher.Update()
	if err != nil {
		return errors.Wrap(err, "casbin: watcher update failed")
	}
	return nil
}

func RemovePolicy(pType string, rule []string) error {
	if pType == "p" {
		ok, err := enforcer.RemovePolicy(rule)
		if err != nil {
			return errors.Wrap(err, "casbin: remove p type policy fail")
		}

		if !ok {
			return errors.New("remove p type casbin policy fail")
		}

		err = watcher.Update()
		if err != nil {
			return errors.Wrap(err, "casbin: watcher update failed")
		}
		return nil
	}

	ok, err := enforcer.RemoveGroupingPolicy(rule)
	if err != nil {
		return errors.Wrap(err, "casbin: remove g type policy fail")
	}

	if !ok {
		return errors.New("remove g type casbin policy fail")
	}

	err = watcher.Update()
	if err != nil {
		return errors.Wrap(err, "casbin: watcher update failed")
	}
	return nil
}

func RemovePolicies(ptype string, rules [][]string) error {
	if ptype == "p" {
		ok, err := enforcer.RemovePolicies(rules)
		if err != nil {
			return errors.Wrap(err, "casbin: remove p type policies fail")
		}

		if !ok {
			return errors.New("remove p type casbin policies fail")
		}

		err = watcher.Update()
		if err != nil {
			return errors.Wrap(err, "casbin: watcher update failed")
		}
		return nil
	}

	ok, err := enforcer.RemoveGroupingPolicies(rules)
	if err != nil {
		return errors.Wrap(err, "casbin: remove g type policies fail")
	}

	if !ok {
		return errors.New("remove g type casbin policies fail")
	}

	err = watcher.Update()
	if err != nil {
		return errors.Wrap(err, "casbin: watcher update failed")
	}
	return nil
}

func RemoveMenuPolicy(menuPath string) error {
	var oldPath []map[string]string
	err := json.Unmarshal([]byte(menuPath), &oldPath)
	if err != nil {
		return errors.Wrap(err, "json: json unmarshal fail")
	}
	for _, v := range oldPath {
		ok, err := enforcer.RemoveFilteredPolicy(1, v["api"], v["method"])
		if err != nil {
			return errors.Wrap(err, "casbin: remove type menu filtered policy fail")
		}

		if !ok {
			return errors.New("remove casbin filtered policies fail")
		}

		err = watcher.Update()
		if err != nil {
			return errors.Wrap(err, "casbin: watcher update failed")
		}
	}
	return nil
}

func RemoveRolePolicy(roleId int64) error {
	ok, err := enforcer.RemoveFilteredPolicy(0, CASBIN_P_PREFIX+strconv.FormatInt(roleId, 10))
	if err != nil {
		return errors.Wrap(err, "casbin: remove type role filtered policy fail")
	}

	if !ok {
		return nil
	}

	err = watcher.Update()
	if err != nil {
		return errors.Wrap(err, "casbin: watcher update failed")
	}
	return nil
}

// SavePolicy 当出现更换配置的数据库时，更新内存policy规则至新的数据库
func SavePolicy() (bool, error) {
	err := enforcer.SavePolicy()
	if err != nil {
		return false, errors.Wrap(err, "casbin: save policy fail")
	}
	return true, nil
}
