package database

import (
	"db-operator/pkg/apis/db/v1alpha1"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/viper"
)

func updateSecret(db *v1alpha1.Database, usr *user) (*corev1.Secret, error) {
	var err error
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-db-secret", db.Name),
			Namespace: db.Namespace,
		},
		Data: map[string][]byte{
			"database-host":     []byte(viper.GetString("dbHost")),
			"database-port":     []byte("5432"),
			"database-name":     []byte(db.Name),
			"database-user":     []byte(usr.username),
			"database-password": []byte(usr.password),
		},
	}

	return secret, err
}
