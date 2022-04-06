package mydb

import (
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type UnitTest struct {
	suite.Suite
	DB
}

var (
	DBConnection *sqlx.DB
	Total        int
)

func (u *UnitTest) SetupTest() {
	Conn, err := u.DB.Connect()

	if err != nil {
		u.Fail(err.Error())
	}

	if err != Conn.Ping() {
		u.Fail(err.Error())
	}

	DBConnection = Conn
}

func (u *UnitTest) BeforeTest() {

}

func (u *UnitTest) AfterTest() {

}

func (u *UnitTest) TestDB_ShouldSeeDatabaseTables() {
	u.Run("It should see tables", func() {
		row := DBConnection.QueryRow("SELECT COUNT(*) as 'Total' FROM information_schema.`TABLES` WHERE TABLE_SCHEMA = '" + os.Getenv("DB_DATABASE") + "'")

		err := row.Scan(&Total)

		if err != nil {
			u.Fail(err.Error())
		}

		assert.Greater(u.T(), Total, 10)
	})
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(UnitTest))
}
