package cmd

import (
	"database/sql"
	"fmt"
	"log"
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
	bindAddr string
	mysqlDSN string
)

func init() {
	httpCmd.Flags().StringVarP(&bindAddr, "bind-addr", "b", ":8080", "Http port")
	httpCmd.PersistentFlags().StringVar(&mysqlDSN, "mysql-dsn", "", "MySQL DSN to connect to DB")
}

func startHttpServer(cmd *cobra.Command, args []string) {
	db, err := getDB(mysqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	err = http.New(db, bindAddr).Start()
	if err != nil {
		log.Fatal(err)
	}
}

func getDB(dsn string) (*sql.DB, error) {
	// return nil, nil
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql db : %s", err.Error())
	}
	return db, nil
}
