package backup

/*
 * This file contains structs and functions related to backing up query planner
 * statistics on the master.
 */

import (
	"fmt"
	"strings"

	"github.com/greenplum-db/gpbackup/utils"
	"github.com/lib/pq"
)

func PrintStatisticsStatements(statisticsFile *utils.FileWithByteCount, toc *utils.TOC, tables []Relation, attStats map[uint32][]AttributeStatistic, tupleStats map[uint32]TupleStatistic) {
	start := statisticsFile.ByteCount
	statisticsFile.MustPrintf(`SET allow_system_table_mods="DML";`)
	toc.AddMetadataEntry("", "", "STATISTICS GUC", start, statisticsFile)
	for _, table := range tables {
		PrintStatisticsStatementsForTable(statisticsFile, toc, table, attStats[table.Oid], tupleStats[table.Oid])
	}
}

func PrintStatisticsStatementsForTable(statisticsFile *utils.FileWithByteCount, toc *utils.TOC, table Relation, attStats []AttributeStatistic, tupleStat TupleStatistic) {
	start := statisticsFile.ByteCount
	tupleQuery := GenerateTupleStatisticsQuery(table, tupleStat)
	statisticsFile.MustPrintf("\n\n%s\n", tupleQuery)
	for _, attStat := range attStats {
		attributeQuery := GenerateAttributeStatisticsQuery(table, attStat)
		statisticsFile.MustPrintf("\n\n%s\n", attributeQuery)
	}
	toc.AddMetadataEntry(table.Schema, table.Name, "STATISTICS", start, statisticsFile)
}

func GenerateTupleStatisticsQuery(table Relation, tupleStat TupleStatistic) string {
	tupleQuery := `UPDATE pg_class
SET
	relpages = %d::int,
	reltuples = %f::real
WHERE relname = '%s'
AND relnamespace = %d;`
	return fmt.Sprintf(
		tupleQuery,
		tupleStat.RelPages,
		tupleStat.RelTuples,
		strings.Replace(tupleStat.Table, "'", "''", -1),
		table.SchemaOid)
}

func GenerateAttributeStatisticsQuery(table Relation, attStat AttributeStatistic) string {
	/*
	 * When restoring statistics to a new database, we cannot determine what the
	 * new OID for a given object will be, so we need to perform an explicit cast
	 * from the name to an OID in the below statements, rather than backing up the
	 * OID in the source database.
	 */
	starelidStr := fmt.Sprintf("'%s'::regclass::oid", strings.Replace(table.ToString(), "'", "''", -1))
	// The entry may or may not already exist, so we can't either just UPDATE or just INSERT without a DELETE.
	attributeQuery := fmt.Sprintf(`DELETE FROM pg_statistic WHERE starelid = %s AND staattnum = %d;

INSERT INTO pg_statistic VALUES (
	%s,
	%d::smallint,
	%f::real,
	%d::integer,
	%f::real,`, starelidStr, attStat.AttNumber, starelidStr, attStat.AttNumber, attStat.NullFraction, attStat.Width, attStat.Distinct)

	/*
	 * If a type name starts with exactly one underscore, it describes an array
	 * type.  We can't restore statistics of array columns, so we'll zero and
	 * NULL everything out.
	 */
	if len(attStat.Type) > 1 && attStat.Type[0] == '_' && attStat.Type[1] != '_' {
		attributeQuery += `
	0::smallint,
	0::smallint,
	0::smallint,
	0::smallint,
	0::oid,
	0::oid,
	0::oid,
	0::oid,
	NULL::real[],
	NULL::real[],
	NULL::real[],
	NULL::real[],
	NULL,
	NULL,
	NULL,
	NULL`
	} else {
		attributeQuery += fmt.Sprintf(`
	%d::smallint,
	%d::smallint,
	%d::smallint,
	%d::smallint,
	%d::oid,
	%d::oid,
	%d::oid,
	%d::oid,
	%s::real[],
	%s::real[],
	%s::real[],
	%s::real[],
	%s,
	%s,
	%s,
	%s`, attStat.Kind1,
			attStat.Kind2,
			attStat.Kind3,
			attStat.Kind4,
			attStat.Operator1,
			attStat.Operator2,
			attStat.Operator3,
			attStat.Operator4,
			RealValues(attStat.Numbers1),
			RealValues(attStat.Numbers2),
			RealValues(attStat.Numbers3),
			RealValues(attStat.Numbers4),
			AnyValues(attStat.Values1, attStat.Type),
			AnyValues(attStat.Values2, attStat.Type),
			AnyValues(attStat.Values3, attStat.Type),
			AnyValues(attStat.Values4, attStat.Type))
	}
	attributeQuery += `
);`
	return attributeQuery
}

func SliceToPostgresArray(slice []string) string {
	quotedStrings := make([]string, len(slice))
	for i, str := range slice {
		escapedStr := strings.Replace(str, "'", "''", -1)
		escapedStr = strings.Replace(escapedStr, `"`, `\"`, -1)
		quotedStrings[i] = fmt.Sprintf(`"%s"`, escapedStr)
	}
	return fmt.Sprintf("'{%s}'", strings.Join(quotedStrings, ","))
}

func RealValues(reals pq.StringArray) string {
	if len(reals) > 0 {
		return SliceToPostgresArray(reals)
	}
	return "NULL"
}

/*
 * A given type is not guaranteed to have a corresponding array type, so we need
 * to use array_in() instead of casting to an array.
 */
func AnyValues(any pq.StringArray, typ string) string {
	if len(any) > 0 {
		return fmt.Sprintf("array_in(%s, '%s'::regtype::oid, -1)", SliceToPostgresArray(any), typ)
	}
	return fmt.Sprintf("NULL")
}
