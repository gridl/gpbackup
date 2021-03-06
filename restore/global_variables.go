package restore

import "github.com/greenplum-db/gpbackup/utils"

/*
 * This file contains global variables and setter functions for those variables
 * used in testing.
 */

/*
 * Non-flag variables
 */

var (
	backupConfig  *utils.BackupConfig
	connection    *utils.DBConn
	globalCluster utils.Cluster
	globalTOC     *utils.TOC
	logger        *utils.Logger
	version       string
)

/*
 * Command-line flags
 */

var (
	backupDir       *string
	createdb        *bool
	debug           *bool
	numJobs         *int
	onErrorContinue *bool
	printVersion    *bool
	quiet           *bool
	redirect        *string
	restoreGlobals  *bool
	timestamp       *string
	verbose         *bool
	withStats       *bool
)

/*
 * Setter functions
 */

func SetConnection(conn *utils.DBConn) {
	connection = conn
}

func SetLogger(log *utils.Logger) {
	logger = log
}

func SetNumJobs(jobs int) {
	numJobs = &jobs
}

func SetTOC(toc *utils.TOC) {
	globalTOC = toc
}
