package backup

import "github.com/greenplum-db/gpbackup/utils"

/*
 * This file contains global variables and setter functions for those variables
 * used in testing.
 */

/*
 * Non-flag variables
 */
var (
	backupReport  *utils.Report
	connection    *utils.DBConn
	globalCluster utils.Cluster
	globalTOC     *utils.TOC
	logger        *utils.Logger
	objectCounts  map[string]int
	version       string
)

/*
 * Command-line flags
 */
var (
	backupDir         *string
	backupGlobals     *bool
	dataOnly          *bool
	dbname            *string
	debug             *bool
	excludeSchemas    utils.ArrayFlags
	excludeTableFile  *string
	excludeTables     utils.ArrayFlags
	includeSchemas    utils.ArrayFlags
	includeTableFile  *string
	includeTables     utils.ArrayFlags
	leafPartitionData *bool
	metadataOnly      *bool
	noCompression     *bool
	printVersion      *bool
	quiet             *bool
	verbose           *bool
	withStats         *bool
)

/*
 * Setter functions
 */

func SetConnection(conn *utils.DBConn) {
	connection = conn
}

func SetCluster(cluster utils.Cluster) {
	globalCluster = cluster
}

func SetExcludeSchemas(schemas []string) {
	excludeSchemas = schemas
}

func SetIncludeSchemas(schemas []string) {
	includeSchemas = schemas
}

func SetExcludeTables(tables []string) {
	excludeTables = tables
}

func SetIncludeTables(tables []string) {
	includeTables = tables
}

func SetLeafPartitionData(which bool) {
	leafPartitionData = &which
}

func SetLogger(log *utils.Logger) {
	logger = log
}

func SetTOC(toc *utils.TOC) {
	globalTOC = toc
}

func SetVersion(v string) {
	version = v
}
