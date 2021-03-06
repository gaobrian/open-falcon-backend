package db

import (
	"database/sql"
	. "gopkg.in/check.v1"
	"github.com/gaobrian/open-falcon-backend/modules/hbs/g"
	commonDb "github.com/gaobrian/open-falcon-backend/common/db"
	commonModel "github.com/gaobrian/open-falcon-backend/common/model"
	hbstesting "github.com/gaobrian/open-falcon-backend/modules/hbs/testing"
)

type TestAgentSuite struct{}

var _ = Suite(&TestAgentSuite{})

type testCaseOfUpdateAgent struct {
	ip string
	agentVersion string
	pluginVersion string
}

// Tests the refresh(insert or update) of agent information
func (suite *TestAgentSuite) TestUpdateAgent(c *C) {
	testCases := []testCaseOfUpdateAgent {
		{ "1.2.3.4", "1.0", "1.0" },
		{ "1.9.3.4", "1.1", "1.1" },
	}

	for _, testCase := range testCases {
		agentInfo := &commonModel.AgentUpdateInfo {
			0,
			&commonModel.AgentReportRequest {
				Hostname: "test-host-1",
				IP: testCase.ip,
				AgentVersion: testCase.agentVersion,
				PluginVersion: testCase.pluginVersion,
			},
		}

		c.Assert(UpdateAgent(agentInfo), IsNil)

		assertUpdateAgent(c, &testCase)
	}
}

func assertUpdateAgent(c *C, testCase *testCaseOfUpdateAgent) {
	dbCtrl := commonDb.NewDbController(DB)

	dbCtrl.QueryForRow(
		commonDb.RowCallbackFunc(func(row *sql.Row) {
			var ip, agentVersion, pluginVersion string

			err := row.Scan(&ip, &agentVersion, &pluginVersion)
			commonDb.DbPanic(err)

			c.Assert(ip, Equals, testCase.ip)
			c.Assert(agentVersion, Equals, testCase.agentVersion)
			c.Assert(pluginVersion, Equals, testCase.pluginVersion)
		}),
		`
		SELECT ip, agent_version, plugin_Version
		FROM host
		WHERE hostname = 'test-host-1'
		`,
	);
}

func (s *TestAgentSuite) SetUpSuite(c *C) {
	(&TestDbSuite{}).SetUpSuite(c)
}

func (s *TestAgentSuite) TearDownSuite(c *C) {
	(&TestDbSuite{}).TearDownSuite(c)
}

func (s *TestAgentSuite) SetUpTest(c *C) {
	if !hbstesting.HasDbEnvForMysqlOrSkip(c) {
		return
	}

	switch c.TestName() {
	case "TestAgentSuite.TestUpdateAgent":
		g.SetConfig(&g.GlobalConfig{
			Hosts: "",
		})
	}
}
func (s *TestAgentSuite) TearDownTest(c *C) {
	dbCtrl := commonDb.NewDbController(DB)

	switch c.TestName() {
	case "TestAgentSuite.TestUpdateAgent":
		g.SetConfig(nil)
		dbCtrl.Exec("DELETE FROM host WHERE hostname = 'test-host-1'")
	}
}
