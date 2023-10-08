// Copyright (C) 2019-2022 Chrystian Huot <chrystian.huot@saubeo.solutions>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

type Database struct {
	Config         *Config
	DateTimeFormat string
	Sql            *sql.DB
}

func NewDatabase(config *Config) *Database {
	var err error

	database := &Database{Config: config}

	switch config.DbType {
	case DbTypeSqlite:
		database.DateTimeFormat = "2006-01-02 15:04:05.000 -07:00"

		dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout%%3d10000", config.GetDbFilePath())

		if database.Sql, err = sql.Open("sqlite", dsn); err != nil {
			log.Fatal(err)
		}

	case DbTypeMariadb, DbTypeMysql:
		database.DateTimeFormat = "2006-01-02 15:04:05"

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)

		if database.Sql, err = sql.Open("mysql", dsn); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("unknown database type %s\n", config.DbType)
	}

	database.Sql.SetConnMaxLifetime(time.Minute)
	database.Sql.SetMaxIdleConns(25)
	database.Sql.SetMaxOpenConns(25)

	if err = database.migrate(); err != nil {
		log.Fatal(err)
	}

	if err = database.seed(); err != nil {
		log.Fatal(err)
	}

	return database
}

func (db *Database) ParseDateTime(f any) (time.Time, error) {
	switch v := f.(type) {
	case []uint8:
		return time.Parse(db.DateTimeFormat, string(v))
	case string:
		return time.Parse(db.DateTimeFormat, v)
	case time.Time:
		return v, nil
	default:
		return time.Time{}, fmt.Errorf("unknown datetime format %T", v)
	}
}

func (db *Database) migrate() error {
	var (
		err     error
		verbose bool
	)

	verbose, err = db.prepareMigration()

	if err == nil {
		err = db.migration20191028144433(verbose)
	}
	if err == nil {
		err = db.migration20191029092201(verbose)
	}
	if err == nil {
		err = db.migration20191126135515(verbose)
	}
	if err == nil {
		err = db.migration20191220093214(verbose)
	}
	if err == nil {
		err = db.migration20200123094105(verbose)
	}
	if err == nil {
		err = db.migration20200428132918(verbose)
	}
	if err == nil {
		err = db.migration20210115105958(verbose)
	}
	if err == nil {
		err = db.migration20210830092027(verbose)
	}
	if err == nil {
		err = db.migration20211202094819(verbose)
	}
	if err == nil {
		err = db.migration20220101070000(verbose)
	}

	return err
}

func (db *Database) migrateWithSchema(name string, schemas []string, verbose bool) error {
	var (
		count int = 0
		err   error
		query string
		tx    *sql.Tx
	)

	formatError := func(err error, query string) error {
		return fmt.Errorf("%s while doing %s", err.Error(), query)
	}

	query = fmt.Sprintf("select count(*) from `freeScannerMeta` where `name` = '%s'", name)
	if err = db.Sql.QueryRow(query).Scan(&count); err != nil {
		return formatError(err, query)
	}

	if count == 0 {
		if verbose {
			log.Printf("running database migration %s", name)
		}

		if tx, err = db.Sql.Begin(); err == nil {
			for _, query = range schemas {
				if _, err = tx.Exec(query); err != nil {
					tx.Rollback()
					return formatError(err, query)
				}
			}

			query = fmt.Sprintf("insert into `freeScannerMeta` (`name`) values ('%s')", name)
			if _, err = tx.Exec(query); err != nil {
				tx.Rollback()
				return formatError(err, query)
			}

			if err = tx.Commit(); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return nil
}

func (db *Database) migration20191028144433(verbose bool) error {
	var queries []string
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"create table `freeScannerSystems` (`id` integer primary key autoincrement, `createdAt` datetime not null, `updatedAt` datetime not null, `name` varchar(255) not null, `system` integer not null, `talkgroups` json not null)",
			"create unique index `free_scanner_systems_system` on `freeScannerSystems` (`system`)",
		}
	} else {
		queries = []string{
			"create table `freeScannerSystems` (`id` integer primary key auto_increment, `createdAt` datetime not null, `updatedAt` datetime not null, `name` varchar(255) not null, `system` integer not null, `talkgroups` json not null)",
			"create unique index `free_scanner_systems_system` on `freeScannerSystems` (`system`)",
		}
	}
	return db.migrateWithSchema("20191028144433-create-freescanner-system", queries, verbose)
}

func (db *Database) migration20191029092201(verbose bool) error {
	var queries []string
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"create table `freeScannerCalls` (`id` integer primary key autoincrement, `createdAt` datetime not null, `updatedAt` datetime not null, `audio` longblob not null, `emergency` tinyint(1) not null, `freq` integer not null, `freqList` json not null, `startTime` datetime not null, `stopTime` datetime not null, `srcList` json not null, `system` integer not null, `talkgroup` integer not null)",
			"create index `free_scanner_calls_start_time` on `freeScannerCalls` (`startTime`)",
			"create index `free_scanner_calls_system` on `freeScannerCalls` (`system`)",
			"create index `free_scanner_calls_talkgroup` on `freeScannerCalls` (`talkgroup`)",
		}
	} else {
		queries = []string{
			"create table `freeScannerCalls` (`id` integer primary key auto_increment, `createdAt` datetime not null, `updatedAt` datetime not null, `audio` longblob not null, `emergency` tinyint(1) not null, `freq` integer not null, `freqList` json not null, `startTime` datetime not null, `stopTime` datetime not null, `srcList` json not null, `system` integer not null, `talkgroup` integer not null)",
			"create index `free_scanner_calls_start_time` on `freeScannerCalls` (`startTime`)",
			"create index `free_scanner_calls_system` on `freeScannerCalls` (`system`)",
			"create index `free_scanner_calls_talkgroup` on `freeScannerCalls` (`talkgroup`)",
		}
	}
	return db.migrateWithSchema("20191029092201-create-freescanner-call", queries, verbose)
}

func (db *Database) migration20191126135515(verbose bool) error {
	var queries []string
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"drop index `free_scanner_calls_system`",
			"drop index `free_scanner_calls_talkgroup`",
		}
	} else {
		queries = []string{
			"drop index `free_scanner_calls_system` on `freeScannerCalls`",
			"drop index `free_scanner_calls_talkgroup` on `freeScannerCalls`",
		}
	}
	return db.migrateWithSchema("20191126135515-optimize-freescanner-calls", queries, verbose)
}

func (db *Database) migration20191220093214(verbose bool) error {
	var queries []string
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"alter table `freeScannerCalls` add column `audioName` varchar(255)",
			"alter table `freeScannerCalls` add column `audioType` varchar(255)",
			"alter table `freeScannerSystems` add column `aliases` json not null",
		}
	} else {
		queries = []string{
			"alter table `freeScannerCalls` add column `audioName` varchar(255)",
			"alter table `freeScannerCalls` add column `audioType` varchar(255)",
			"alter table `freeScannerSystems` add column `aliases` json not null",
		}
	}
	return db.migrateWithSchema("20191220093214-new-v3-tables", queries, verbose)
}

func (db *Database) migration20200123094105(verbose bool) error {
	var queries []string
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"create index `free_scanner_calls_system` on `freeScannerCalls` (`system`)",
			"create index `free_scanner_calls_system_talkgroup` on `freeScannerCalls` (`system`, `talkgroup`)",
		}
	} else {
		queries = []string{
			"create index `free_scanner_calls_system` on `freeScannerCalls` (`system`)",
			"create index `free_scanner_calls_system_talkgroup` on `freeScannerCalls` (`system`, `talkgroup`)",
		}
	}
	return db.migrateWithSchema("20200123094105-optimize-freescanner-calls", queries, verbose)
}

func (db *Database) migration20200428132918(verbose bool) error {
	var queries []string
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"drop table `freeScannerSystems`",
			"create table `freeScannerCalls2` (`id` integer primary key autoincrement, `audio` longblob not null, `audioName` varchar(255), `audioType` varchar(255), `dateTime` datetime not null, `frequencies` json not null, `frequency` integer, `source` integer, `sources` json not null, `system` integer not null, `talkgroup` integer not null)",
			"insert into `freeScannerCalls2` select `id`, `audio`, `audioName`, `audioType`, `startTime`, `freqList`, `freq`, null, `srcList`, `system`, `talkgroup` from `freeScannerCalls`",
			"drop table `freeScannerCalls`",
			"alter table `freeScannerCalls2` rename to `freeScannerCalls`",
			"create index `free_scanner_calls_date_time_system_talkgroup` on `freeScannerCalls` (`dateTime`, `system`, `talkgroup`)",
		}
	} else {
		queries = []string{
			"drop table `freeScannerSystems`",
			"create table `freeScannerCalls2` (`id` integer primary key auto_increment, `audio` longblob not null, `audioName` varchar(255), `audioType` varchar(255), `dateTime` datetime not null, `frequencies` json not null, `frequency` integer, `source` integer, `sources` json not null, `system` integer not null, `talkgroup` integer not null)",
			"insert into `freeScannerCalls2` select `id`, `audio`, `audioName`, `audioType`, `startTime`, `freqList`, `freq`, null, `srcList`, `system`, `talkgroup` from `freeScannerCalls`",
			"drop table `freeScannerCalls`",
			"alter table `freeScannerCalls2` rename to `freeScannerCalls`",
			"create index `free_scanner_calls_date_time_system_talkgroup` on `freeScannerCalls` (`dateTime`, `system`, `talkgroup`)",
		}
	}
	return db.migrateWithSchema("20200428132918-new-v4-tables", queries, verbose)
}

func (db *Database) migration20210115105958(verbose bool) error {
	var queries []string
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"create table `freeScannerAccesses` (`_id` integer primary key autoincrement, `code` varchar(255) not null unique, `expiration` datetime, `ident` varchar(255), `limit` integer, `order` integer, `systems` text not null)",
			"create table `freeScannerApiKeys` (`_id` integer primary key autoincrement, `disabled` tinyint(1) default 0, `ident` varchar(255), `key` varchar(255) not null unique, `order` integer, `systems` text not null)",
			"create table `freeScannerCalls2` (`id` integer primary key autoincrement, `audio` longblob not null, `audioName` varchar(255), `audioType` varchar(255), `dateTime` datetime not null, `frequencies` text not null, `frequency` integer, `source` integer, `sources` text not null, `system` integer not null, `talkgroup` integer not null)",
			"create index `free_scanner_calls2_date_time_system_talkgroup` on `freeScannerCalls2` (`dateTime`, `system`, `talkgroup`)",
			"insert into `freeScannerCalls2` select `id`, `audio`, `audioName`, `audioType`, `dateTime`, `frequencies`, `frequency`, `source`, `sources`, `system`, `talkgroup` from `freeScannerCalls`",
			"drop table `freeScannerCalls`",
			"alter table `freeScannerCalls2` rename to `freeScannerCalls`",
			"create table `freeScannerConfigs` (`_id` integer primary key autoincrement, `key` varchar(255) not null unique, `val` text not null)",
			"create index `free_scanner_configs_key` on `freeScannerConfigs` (`key`)",
			"create table `freeScannerDirWatches` (`_id` integer primary key autoincrement, `delay` integer default 0, `deleteAfter` tinyint(1) default 0, `directory` varchar(255) not null unique, `disabled` tinyint(1) default 0, `extension` varchar(255), `frequency` integer, `mask` varchar(255), `order` integer, `systemId` integer, `talkgroupId` integer, `type` varchar(255), `usePolling` tinyint(1) default 0)",
			"create table `freeScannerDownstreams` (`_id` integer primary key autoincrement, `apiKey` varchar(255) not null unique, `disabled` tinyint(1) default 0, `order` integer, `systems` text not null, `url` varchar(255) not null)",
			"create table `freeScannerGroups` (`_id` integer primary key autoincrement, `label` varchar(255) not null)",
			"create table `freeScannerLogs` (`_id` integer primary key autoincrement, `dateTime` datetime not null, `level` varchar(255) not null, `message` varchar(255) not null)",
			"create index `free_scanner_logs_date_time_level` on `freeScannerLogs` (`dateTime`, `level`)",
			"create table `freeScannerSystems` (`_id` integer primary key autoincrement, `autoPopulate` tinyint(1) default 0, `blacklists` text not null, `id` integer not null unique, `label` varchar(255) not null, `led` varchar(255), `order` integer, `talkgroups` text not null, `units` text not null)",
			"create table `freeScannerTags` (`_id` integer primary key autoincrement, `label` varchar(255) not null)",
		}
	} else {
		queries = []string{
			"create table `freeScannerAccesses` (`_id` integer primary key auto_increment, `code` varchar(255) not null unique, `expiration` datetime, `ident` varchar(255), `limit` integer, `order` integer, `systems` text not null)",
			"create table `freeScannerApiKeys` (`_id` integer primary key auto_increment, `disabled` tinyint(1) default 0, `ident` varchar(255), `key` varchar(255) not null unique, `order` integer, `systems` text not null)",
			"create table `freeScannerCalls2` (`id` integer primary key auto_increment, `audio` longblob not null, `audioName` varchar(255), `audioType` varchar(255), `dateTime` datetime not null, `frequencies` text not null, `frequency` integer, `source` integer, `sources` text not null, `system` integer not null, `talkgroup` integer not null)",
			"create index `free_scanner_calls2_date_time_system_talkgroup` on `freeScannerCalls2` (`dateTime`, `system`, `talkgroup`)",
			"insert into `freeScannerCalls2` select `id`, `audio`, `audioName`, `audioType`, `dateTime`, `frequencies`, `frequency`, `source`, `sources`, `system`, `talkgroup` from `freeScannerCalls`",
			"drop table `freeScannerCalls`",
			"alter table `freeScannerCalls2` rename to `freeScannerCalls`",
			"create table `freeScannerConfigs` (`_id` integer primary key auto_increment, `key` varchar(255) not null unique, `val` text not null)",
			"create index `free_scanner_configs_key` on `freeScannerConfigs` (`key`)",
			"create table `freeScannerDirWatches` (`_id` integer primary key auto_increment, `delay` integer default 0, `deleteAfter` tinyint(1) default 0, `directory` varchar(255) not null unique, `disabled` tinyint(1) default 0, `extension` varchar(255), `frequency` integer, `mask` varchar(255), `order` integer, `systemId` integer, `talkgroupId` integer, `type` varchar(255), `usePolling` tinyint(1) default 0)",
			"create table `freeScannerDownstreams` (`_id` integer primary key auto_increment, `apiKey` varchar(255) not null unique, `disabled` tinyint(1) default 0, `order` integer, `systems` text not null, `url` varchar(255) not null)",
			"create table `freeScannerGroups` (`_id` integer primary key auto_increment, `label` varchar(255) not null)",
			"create table `freeScannerLogs` (`_id` integer primary key auto_increment, `dateTime` datetime not null, `level` varchar(255) not null, `message` varchar(255) not null)",
			"create index `free_scanner_logs_date_time_level` on `freeScannerLogs` (`dateTime`, `level`)",
			"create table `freeScannerSystems` (`_id` integer primary key auto_increment, `autoPopulate` tinyint(1) default 0, `blacklists` text not null, `id` integer not null unique, `label` varchar(255) not null, `led` varchar(255), `order` integer, `talkgroups` text not null, `units` text not null)",
			"create table `freeScannerTags` (`_id` integer primary key auto_increment, `label` varchar(255) not null)",
		}
	}
	return db.migrateWithSchema("20210115105958-new-v5.1-tables", queries, verbose)
}

func (db *Database) migration20210830092027(verbose bool) error {
	var queries []string
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"create table `freeScannerSystems2` (`_id` integer primary key autoincrement, `autoPopulate` tinyint(1) default 0, `blacklists` text not null, `id` integer not null unique, `label` varchar(255) not null, `led` varchar(255), `order` integer, `talkgroups` longtext not null, `units` longtext not null)",
			"insert into `freeScannerSystems2` select `_id`, `autoPopulate`, `blacklists`, `id`, `label`, `led`, `order`, `talkgroups`, `units` from `freeScannerSystems`",
			"drop table `freeScannerSystems`",
			"alter table `freeScannerSystems2` rename to `freeScannerSystems`",
			"drop index `free_scanner_calls2_date_time_system_talkgroup`",
			"create index `free_scanner_calls_date_time_system_talkgroup` on `freeScannerCalls` (`dateTime`, `system`, `talkgroup`)",
		}

	} else {
		queries = []string{
			"alter table `freeScannerSystems` modify `talkgroups` longtext null not null",
			"alter table `freeScannerSystems` modify `units` longtext null not null",
			"drop index `free_scanner_calls2_date_time_system_talkgroup` on `freeScannerCalls`",
			"create index `free_scanner_calls_date_time_system_talkgroup` on `freeScannerCalls` (`dateTime`, `system`, `talkgroup`)",
		}
	}
	return db.migrateWithSchema("20210830092027-v6.0-rename-index", queries, verbose)
}

func (db *Database) migration20211202094819(verbose bool) error {
	var queries []string
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"alter table `freeScannerDownstreams` rename to `freeScannerDownstreams2`",
			"create table `freeScannerDownstreams` (`_id` integer primary key autoincrement, `apiKey` varchar(255) not null, `disabled` tinyint(1) default 0, `order` integer, `systems` text not null, `url` varchar(255) not null)",
			"insert into `freeScannerDownstreams` select * from `freeScannerDownstreams2`",
			"drop table `freeScannerDownstreams2`",
		}
	} else {
		queries = []string{
			"alter table `freeScannerDownstreams` rename to `freeScannerDownstreams2`",
			"create table `freeScannerDownstreams` (`_id` integer primary key auto_increment, `apiKey` varchar(255) not null, `disabled` tinyint(1) default 0, `order` integer, `systems` text not null, `url` varchar(255) not null)",
			"insert into `freeScannerDownstreams` select * from `freeScannerDownstreams2`",
			"drop table `freeScannerDownstreams2`",
		}
	}
	return db.migrateWithSchema("20211202094819-v6.0.2-alter-table", queries, verbose)
}

func (db *Database) migration20220101070000(verbose bool) error {
	var (
		err        error
		frequency  any
		id         uint
		label      string
		led        any
		name       string
		queries    []string
		rows       *sql.Rows
		stra       string
		strb       string
		talkgroups []*Talkgroup
		units      []*Unit
	)
	if db.Config.DbType == DbTypeSqlite {
		queries = []string{
			"create table `freeScannerCalls2` (`id` integer primary key autoincrement, `audio` longblob not null, `audioName` varchar(255), `audioType` varchar(255), `dateTime` datetime not null, `frequencies` text not null, `frequency` integer, `patches` text not null, `source` integer, `sources` text not null, `system` integer not null, `talkgroup` integer not null)",
			"insert into `freeScannerCalls2` select `id`, `audio`, `audioName`, `audioType`, `dateTime`, `frequencies`, `frequency`, '[]', `source`, `sources`, `system`, `talkgroup` from `freeScannerCalls`",
			"drop table `freeScannerCalls`",
			"alter table `freeScannerCalls2` rename to `freeScannerCalls`",
			"create index `free_scanner_calls_date_time_system_talkgroup` on `freeScannerCalls` (`dateTime`, `system`, `talkgroup`)",
			"create table `freeScannerSystems2` (`_id` integer primary key autoincrement, `autoPopulate` tinyint(1) default 0, `blacklists` text not null, `id` integer not null unique, `label` varchar(255) not null, `led` varchar(255), `order` integer)",
			"insert into `freeScannerSystems2` select `_id`, `autoPopulate`, `blacklists`, `id`, `label`, `led`, `order` from `freeScannerSystems`",
			"drop table `freeScannerSystems`",
			"alter table `freeScannerSystems2` rename to `freeScannerSystems`",
			"create table `freeScannerTalkgroups` (`_id` integer primary key autoincrement, `frequency` integer, `groupId` integer not null, `id` integer not null, `label` varchar(255) not null, `led` varchar(255), `name` varchar(255) not null, `order` integer, `systemId` integer not null, `tagId` integer not null)",
			"create unique index `free_scanner_talkgroups_system_id_id` on `freeScannerTalkgroups` (`systemId`, `id`)",
			"create table `freeScannerUnits` (`_id` integer primary key autoincrement, `id` integer not null, `label` varchar(255) not null, `order` integer, `systemId` integer not null)",
			"create unique index `free_scanner_units_system_id_id` on `freeScannerUnits` (`systemId`, `id`)",
		}
	} else {
		queries = []string{
			"alter table `freeScannerCalls` add column `patches` text not null",
			"alter table `freeScannerSystems` drop column `talkgroups`",
			"alter table `freeScannerSystems` drop column `units`",
			"create table `freeScannerTalkgroups` (`_id` integer primary key auto_increment, `frequency` integer, `groupId` integer not null, `id` integer not null, `label` varchar(255) not null, `led` varchar(255), `name` varchar(255) not null, `order` integer, `systemId` integer not null, `tagId` integer not null)",
			"create unique index `free_scanner_talkgroups_system_id_id` on `freeScannerTalkgroups` (`systemId`, `id`)",
			"create table `freeScannerUnits` (`_id` integer primary key auto_increment, `id` integer not null, `label` varchar(255) not null, `order` integer, `systemId` integer not null)",
			"create unique index `free_scanner_units_system_id_id` on `freeScannerUnits` (`systemId`, `id`)",
		}
	}
	if rows, err = db.Sql.Query("select `id`, `talkgroups`, `units` from `freeScannerSystems`"); err == nil {
		for rows.Next() {
			if err = rows.Scan(&id, &stra, &strb); err != nil {
				break
			}
			if err = json.Unmarshal([]byte(stra), &talkgroups); err != nil {
				break
			}
			if err = json.Unmarshal([]byte(strb), &units); err != nil {
				break
			}
			for i, tg := range talkgroups {
				switch v := tg.Frequency.(type) {
				case uint:
					frequency = v
				default:
					frequency = "null"
				}
				label = strings.ReplaceAll(tg.Label, "'", "''")
				switch v := tg.Led.(type) {
				case string:
					led = fmt.Sprintf("'%v'", strings.ReplaceAll(v, "'", "''"))
				default:
					led = "null"
				}
				name = strings.ReplaceAll(tg.Name, "'", "''")
				tg.Order = uint(i + 1)
				queries = append(queries, fmt.Sprintf("insert into `freeScannerTalkgroups` (`frequency`, `groupId`, `id`, `label`, `led`, `name`, `order`, `systemId`, `tagId`) values (%v, %v, %v, '%v', %v, '%v', %v, %v, %v)", frequency, tg.GroupId, tg.Id, label, led, name, tg.Order, id, tg.TagId))
			}
			for i, unit := range units {
				label = strings.ReplaceAll(unit.Label, "'", "''")
				unit.Order = uint(i + 1)
				queries = append(queries, fmt.Sprintf("insert into `freeScannerUnits` (`id`, `label`, `order`, `systemId`) values (%v, '%v', %v, %v)", unit.Id, label, unit.Order, id))
			}
		}
		rows.Close()
		if err != nil {
			return err
		}
	}
	return db.migrateWithSchema("20220101070000-v6.1.0", queries, verbose)
}

func (db *Database) prepareMigration() (bool, error) {
	var (
		err     error
		verbose bool = true
		query   string
	)

	query = "select count(*) as count from `freeScannerMeta`"
	if _, err = db.Sql.Exec(query); err != nil {
		query = "select count(*) as count from `SequelizeMeta`"
		if _, err = db.Sql.Exec(query); err == nil {
			log.Println("Preparing for database migration")
			query = "alter table `SequelizeMeta` rename to `freeScannerMeta`"
			_, err = db.Sql.Exec(query)
		} else {
			verbose = false
			query = "create table `freeScannerMeta` (name varchar(255) not null unique primary key)"
			_, err = db.Sql.Exec(query)
		}
	}

	return verbose, err
}

func (db *Database) seed() error {
	if err := db.seedGroups(); err != nil {
		return err
	}

	if err := db.seedTags(); err != nil {
		return err
	}

	return nil
}

func (db *Database) seedGroups() error {
	var count uint

	formatError := func(err error) error {
		return fmt.Errorf("database.seedgroups: %s", err.Error())
	}

	if err := db.Sql.QueryRow("select count(*) from `freeScannerGroups`").Scan(&count); err != nil {
		return formatError(err)
	}

	if count == 0 {
		if tx, err := db.Sql.Begin(); err == nil {
			for _, group := range defaults.groups {
				if _, err := tx.Exec("insert into `freeScannerGroups` (`label`) values (?)", group); err != nil {
					tx.Rollback()
					return formatError(err)
				}
			}

			if err := tx.Commit(); err != nil {
				return formatError(err)
			}

		} else {
			return formatError(err)
		}
	}

	return nil
}

func (db *Database) seedTags() error {
	var count uint

	formatError := func(err error) error {
		return fmt.Errorf("database.seedtags: %s", err.Error())
	}

	if err := db.Sql.QueryRow("select count(*) from `freeScannerTags`").Scan(&count); err != nil {
		return formatError(err)
	}

	if count == 0 {
		if tx, err := db.Sql.Begin(); err == nil {
			for _, group := range defaults.tags {
				if _, err := tx.Exec("insert into `freeScannerTags` (`label`) values (?)", group); err != nil {
					tx.Rollback()
					return formatError(err)
				}
			}

			if err := tx.Commit(); err != nil {
				return formatError(err)
			}

		} else {
			return formatError(err)
		}
	}

	return nil
}
