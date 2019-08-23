package database

import (
	"database/sql"
	"db-operator/pkg/apis/db/v1alpha1"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type DB struct {
	*sql.DB
}

var dbCon *sql.DB

// Each database type will have it's own config file
func init() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.SetConfigName("postgres")
	log.Info("Initializing postgresql config")
	// TODO: set defaults

	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err, "Failed to read config")
		os.Exit(1)
	}

	connStr := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable",
		viper.GetString("dbDatabase"),
		viper.GetString("dbUser"),
		viper.GetString("dbPassword"),
		viper.GetString("dbHost"))

	dbCon, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Error(err, "Unable to connect to the database")
	}

	if err = dbCon.Ping(); err != nil {
		log.Error(err, "Unable to access database")
	}
}

func postgresUpdateEvent(db *v1alpha1.Database, usr *user) error {
	users := db.Spec.Users
	if db.Status.Phase == "" {
		//TODO
		log.Info("Create event")
		//reqLogger.Info("Creating a new Database", "Db.Namespace", instance.Namespace, "Db.Name", instance.Name)
		// If we have nothing in the status -> it's create event
		err := postgresCreateDB(db.Name)
		if err != nil {
			return err
		}
		err = postgresCreateUser(db, usr)
		if err != nil {
			return err
		}
		users = append(users, usr.username)
		err = postgresGrantAllWithRole(users, db.Name)
		if err != nil {
			return err
		}

		db.Status.Phase = "Created"
		return err

	} else {
		users = append(users, db.Name)
		err := postgresUpdateGrants(users, db.Name)

		return err
	}
}

func postgresCreateDB(dbName string) error {
	query := fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)
	_, err := dbCon.Exec(query)
	if err != nil {
		log.Error(err, "Unable to create database", "Database:", dbName)
		return err
	}
	log.Info("Database was successfully created!")

	return err
}

func postgresCreateUser(db *v1alpha1.Database, usr *user) error {
	password, err := genPassword()
	if err != nil {
		return err
	}
	usr.password = password
	usr.username = db.Name

	query := fmt.Sprintf(`CREATE USER "%s" WITH ENCRYPTED PASSWORD '%s'`, usr.username, usr.password)
	_, err = dbCon.Exec(query)
	if err != nil {
		log.Error(err, "Unable to create user", "User:", usr.username)
		return err
	}

	log.Info("User was successfully created!")

	return err
}

func postgresRevokeUser(user string, role string) error {
	query := fmt.Sprintf(`REVOKE "%s" FROM "%s"`, role, user)
	_, err := dbCon.Exec(query)
	return err
}

func postgresDelUser(userName string) (string, error) {
	return "nil", nil
}

func postgresDelDB(dbName string) error {
	query := fmt.Sprintf(`DROP DATABASE "%s"`, dbName)
	_, err := dbCon.Exec(query)
	if err != nil {
		log.Error(err, "Unable to drop the database", "Database:", dbName)
	}

	return err
}

func postgresGrantAllWithRole(users []string, database string) error {
	roleName := fmt.Sprintf(`%s_owners`, database)
	query := fmt.Sprintf(`CREATE ROLE "%s"`, roleName)
	_, err := dbCon.Exec(query)
	if err != nil {
		log.Error(err, "Unable to create ROLE", "Role prefix:", database)
		return err
	}

	query = fmt.Sprintf(`GRANT ALL on DATABASE "%s" to "%s"`, database, roleName)
	if _, err := dbCon.Exec(query); err != nil {
		return err
	}

	err = postgresGrantAll(users, roleName)

	return err
}

func postgresGrantAll(users []string, database string) error {
	userlist := strings.Join(users, `", "`)
	query := fmt.Sprintf(`GRANT "%s" to "%s"`, database, userlist)
	_, err := dbCon.Exec(query)
	if err != nil {
		log.Error(err, "Unable to assign permissions", "Database:", database, "Users: ", userlist)
	}

	return err
}

func postgresUpdateGrants(users []string, database string) error {
	roleName := fmt.Sprintf(`%s_owners`, database)
	currentUsers, err := getRoleUsers(roleName)
	if err != nil {
		return err
	}
	// Revoke access if users were removed from the object
	for _, u := range currentUsers {
		if !contains(users, u) {
			err = postgresRevokeUser(u, roleName)
			if err != nil {
				return err
			}
		}
	}
	// Grant access for newly created users
	for _, u := range users {
		if !contains(currentUsers, u) {
			err = postgresGrantAll([]string{u}, roleName)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func getRoleUsers(roleName string) ([]string, error) {
	queryString := `select usename
		from pg_user
		join pg_auth_members on (pg_user.usesysid = pg_auth_members.member)
		join pg_roles on (pg_roles.rolname = '%s' AND pg_roles.oid = pg_auth_members.roleid)`

	query := fmt.Sprintf(queryString, roleName)

	rows, err := dbCon.Query(query)
	if err != nil {
		log.Error(err, err.Error())
	}
	defer rows.Close()

	var users []string

	for rows.Next() {
		var rolname string
		if err := rows.Scan(&rolname); err != nil {
			log.Error(err, "Unable to get user role")
			return users, err
		}
		users = append(users, rolname)
	}

	return users, err
}

func postgresDeleteEvent(db *v1alpha1.Database) error {
	if db.Spec.Drop {
		err := postgresDelDB(db.Name)
		if err != nil {
			return err
		}
	} else {
		log.Info("Database won't be deleted! Protected is set to true.")
	}

	roleName := fmt.Sprintf(`%s_owners`, db.Name)
	query := fmt.Sprintf(`DROP ROLE "%s"`, roleName)
	_, err := dbCon.Exec(query)
	if err != nil {
		log.Error(err, "Unable to drop ROLE", "Role:", roleName)
		return err
	} else {
		log.Info("Role was successfully deleted", "Role:", roleName)
	}
	query = fmt.Sprintf(`DROP USER "%s"`, db.Name)
	_, err = dbCon.Exec(query)
	if err != nil {
		log.Error(err, "Unable to drop User", "User:", db.Name)
		return err
	} else {
		log.Info("User was successfully deleted", "User:", db.Name)
	}

	return err
}
