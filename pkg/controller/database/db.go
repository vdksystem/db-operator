package database

import (
	dbv1alpha1 "db-operator/pkg/apis/db/v1alpha1"
)

func createDatabase(db *dbv1alpha1.Database) error {
	switch db.Spec.Type {
	case "mysql":
		log.Info("MySQL is not supported yet", "Db.Namespace", db.Namespace, "Db.Name", db.Name, "Db.Type", db.Spec.Type)
		db.Status.Phase = "Error"
		break

	case "postgres":
		err := postgresCreateDB(db.Name)
		if err != nil {
			log.Error(err, "Failed to create database", "Dbname:", db.Name)
			return err
		}
		db.Status.Phase = "Created"
		break

	default:
		log.Info("Database type required")
		db.Status.Phase = "Error"
		break
	}
	return nil
}
