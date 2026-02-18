package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"test-lbc/http"

	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
)

var httpCmd = &cobra.Command{
	Use:   "http-server",
	Short: "Start http server",
	Long:  "Start a http-server for LBC technical test",
	Run:   startHttpServer,
}

var (
	bindAddr           string
	prometheusBindAddr string
	sqlHost            string
	sqlDB              string
)

func init() {
	httpCmd.Flags().StringVarP(&bindAddr, "bind-addr", "b", ":8080", "Http port")
	httpCmd.Flags().StringVarP(&prometheusBindAddr, "prometheus-bind-addr", "p", ":2112", "prometheus metrics port")
	httpCmd.PersistentFlags().StringVarP(&sqlHost, "mysql-host", "h", "localhost", "MySQL host")
	httpCmd.PersistentFlags().StringVarP(&sqlDB, "mysql-db", "d", "", "MySQL database")

	httpCmd.MarkPersistentFlagRequired("mysql-db")
}

func startHttpServer(cmd *cobra.Command, args []string) {
	db, err := getDB(sqlHost, sqlDB)
	if err != nil {
		log.Fatal(err)
	}

	err = http.New(db, bindAddr, prometheusBindAddr).Start()
	if err != nil {
		log.Fatal(err)
	}
}

func getDB(sqlHost, sqlDB string) (*sql.DB, error) {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, sqlHost, sqlDB)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql db : %s", err.Error())
	}

	return db, nil
}
