package e2e_test

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/calvinchengx/gin-go-pg/e2e"
	"github.com/calvinchengx/gin-go-pg/manager"
	"github.com/calvinchengx/gin-go-pg/model"
	"github.com/calvinchengx/gin-go-pg/repository"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/go-pg/pg/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var superUser *model.User

type E2ETestSuite struct {
	suite.Suite
	db       *pg.DB
	postgres *embeddedpostgres.EmbeddedPostgres
	m        *manager.Manager
}

// SetupSuite runs before all tests in this test suite
func (suite *E2ETestSuite) SetupSuite() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	projectRoot := filepath.Dir(d)
	tmpDir := path.Join(projectRoot, "tmp2")
	testConfig := embeddedpostgres.DefaultConfig().
		Username("db_test_user").
		Password("db_test_password").
		Database("db_test_database").
		Version(embeddedpostgres.V12).
		RuntimePath(tmpDir).
		Port(9877)

	suite.postgres = embeddedpostgres.NewDatabase(testConfig)
	_ = suite.postgres.Start()

	suite.db = pg.Connect(&pg.Options{
		Addr:     "localhost:9877",
		User:     "db_test_user",
		Password: "db_test_password",
		Database: "db_test_database",
	})

	log, _ := zap.NewDevelopment()
	accountRepo := repository.NewAccountRepo(suite.db, log)
	roleRepo := repository.NewRoleRepo(suite.db, log)
	suite.m = manager.NewManager(accountRepo, roleRepo, suite.db)

	superUser, _ = e2e.SetupDatabase(suite.m)
}

// TearDownSuite runs after all tests in this test suite
func (suite *E2ETestSuite) TearDownSuite() {
	suite.postgres.Stop()
}

func (suite *E2ETestSuite) TestGetModels() {
	models := manager.GetModels()
	sql := `SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public';`
	var count int
	res, err := suite.db.Query(pg.Scan(&count), sql, nil)

	assert.NotNil(suite.T(), res)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(models), count)

	sql = `SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';`
	var names pg.Strings
	res, err = suite.db.Query(&names, sql, nil)

	assert.NotNil(suite.T(), res)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(models), len(names))
}

func (suite *E2ETestSuite) TestSuperUser() {
	assert.NotNil(suite.T(), superUser)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
