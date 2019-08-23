package database

import (
	dbv1alpha1 "db-operator/pkg/apis/db/v1alpha1"
)

func updateEvent(db *dbv1alpha1.Database, usr *user) error {
	switch db.Spec.Type {
	case "mysql":
		log.Info("MySQL is not supported yet", "Db.Namespace", db.Namespace, "Db.Name", db.Name, "Db.Type", db.Spec.Type)
		db.Status.Phase = "Unsupported"
		break

	case "postgres":
		err := postgresUpdateEvent(db, usr)
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

func deleteEvent(db *dbv1alpha1.Database) error {
	if db.Spec.Drop {
		err := postgresDeleteEvent(db)
		if err != nil {
			return err
		}
	} else {
		log.Info("Database won't be deleted! Drop is set to false.")
	}
	return nil
}
