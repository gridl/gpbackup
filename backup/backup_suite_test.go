package backup_test

/*
 * This file contains integration tests for gpbackup as a whole, rather than
 * tests relating to functions in any particular file.
 */

import (
	"os/exec"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/greenplum-db/gpbackup/testutils"
	"github.com/greenplum-db/gpbackup/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	connection   *utils.DBConn
	mock         sqlmock.Sqlmock
	logger       *utils.Logger
	stdout       *gbytes.Buffer
	stderr       *gbytes.Buffer
	logfile      *gbytes.Buffer
	buffer       = gbytes.NewBuffer()
	gpbackupPath = ""
	toc          *utils.TOC
	backupfile   *utils.FileWithByteCount
	baseVersion  utils.GPDBVersion
)

/* This function is a helper function to execute gpbackup and return a session
 * to allow checking its output.
 */
func gpbackup() *gexec.Session {
	command := exec.Command(gpbackupPath, "-dbname", "testdb")
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ShouldNot(HaveOccurred())
	<-session.Exited
	return session
}

func TestBackup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "backup tests")
}

var _ = BeforeSuite(func() {
	connection, mock, logger, stdout, stderr, logfile = testutils.SetupTestEnvironment()
	baseVersion = connection.Version
})

var _ = BeforeEach(func() {
	buffer = gbytes.NewBuffer()
	connection.Version = baseVersion
})
