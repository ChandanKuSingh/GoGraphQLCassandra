package repository

import (
	"github.com/ChandanKuSingh/GoGraphQLCassandra/graph/model"
	"github.com/gocql/gocql"
	"log"
)

type DBConnection struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

var connection DBConnection

func SetupDBConnection() {
	connection.cluster = gocql.NewCluster("127.0.0.1")
	connection.cluster.Consistency = gocql.Quorum
	connection.cluster.Keyspace = "first_app"
	connection.session, _ = connection.cluster.CreateSession()
}

func CreateAlert(query string, values ...interface{}) {
	if err := connection.session.Query(query).Bind(values...).Exec(); err != nil {
		log.Fatal("Error while inserting alerts in Cassandra")
		log.Fatal(err)
	}
}

func AddAlert(alert *model.Alert) {
	query := `INSERT INTO alerts(id, system_wwn, status, severity, type, count, description, last_time_occured) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	CreateAlert(query, alert.ID, alert.SystemWwn, alert.Status, alert.Severity, alert.Type, alert.Count, alert.Description, alert.LastOccuredTime)
}

func FindAll() []*model.Alert {
	var result []*model.Alert
	m := map[string]interface{}{}

	iter := connection.session.Query("SELECT * FROM alerts").Iter()
	for iter.MapScan(m) {
		result = append(result, &model.Alert{
			ID:              m["id"].(string),
			Count:           m["count"].(int),
			Description:     m["description"].(string),
			LastOccuredTime: m["last_time_occured"].(string),
			Severity:        m["severity"].(string),
			Status:          m["status"].(string),
			SystemWwn:       m["system_wwn"].(string),
			Type:            m["type"].(string),
		})
		m = map[string]interface{}{}
	}

	return result
}

//func (db *DBConnection) getAllAlerts(query string) {
//	return connection.session.Query("SELECT * FROM alerts").Iter()
//}
