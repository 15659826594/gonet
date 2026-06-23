package database

import (
	"fmt"

	"github.com/spf13/viper"
)

type Database struct {
	Type            string `mapstructure:"type"`
	Hostname        string `mapstructure:"hostname"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Hostport        int    `mapstructure:"hostport"`
	Dsn             string `mapstructure:"dsn"`
	Params          string `mapstructure:"params"`
	Charset         string `mapstructure:"charset"`
	Prefix          string `mapstructure:"prefix"`
	Debug           bool   `mapstructure:"debug"`
	Deploy          int    `mapstructure:"deploy"`
	RwSeparate      bool   `mapstructure:"rw_separate"`
	MasterNum       int    `mapstructure:"master_num"`
	SlaveNo         string `mapstructure:"slave_no"`
	FieldsStrict    bool   `mapstructure:"fields_strict"`
	ResultsetType   string `mapstructure:"resultset_type"`
	AutoTimestamp   bool   `mapstructure:"auto_timestamp"`
	DatetimeFormat  bool   `mapstructure:"datetime_format"`
	SqlExplain      bool   `mapstructure:"sql_explain"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

func (d *Database) Viper(v *viper.Viper) *Database {
	_ = v.Unmarshal(d)
	return d
}

func (d *Database) GetDsn() string {
	if d.Dsn != "" {
		return d.Dsn
	}
	d.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&%s", d.Username, d.Password, d.Hostname, d.Hostport, d.Database, d.Charset, d.Params)
	return d.Dsn
}
