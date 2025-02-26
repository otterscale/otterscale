package or

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

func getOutputLines(pool *sql.DB, statement string) (string, error) {
	var state int
	var output string

	ctx, cancel := context.WithTimeout(context.Background(), 23*time.Second)
	defer cancel()

	sql := `BEGIN DBMS_OUTPUT.ENABLE(:1); END;`
	if _, err := pool.ExecContext(ctx, sql, 0x7FFF); err != nil {
		return "", err
	}

	if _, err := pool.ExecContext(ctx, statement); err != nil {
		return "", err
	}

	sql = `
		DECLARE
			l	VARCHAR2(255);
			s	NUMBER;
			b	LONG;
		BEGIN
			LOOP EXIT
				WHEN LENGTH(b) + 255 > :MAXBYTES OR s = 1;
				dbms_output.get_line( l, s );
				IF LENGTH(l) > 0 THEN
					b := b || l || CHR(10);
				END IF;
			END LOOP;
			:DONE := s;
			:BUFFER := b;
		END;
	`
	if _, err := pool.ExecContext(ctx, sql, 0x7FFF, go_ora.Out{Dest: &state}, go_ora.Out{Dest: &output, Size: 0x7FFF}); err != nil {
		return "", err
	}

	return output, nil
}

func getTypeOIDs(pool *sql.DB, tableName string) (map[string]uint32, error) {
	oids := make(map[string]uint32)

	statement := fmt.Sprintf(`
		DECLARE
			c				NUMBER;
			c_count			INTEGER;
			t_descriptor	DBMS_SQL.DESC_TAB;
		BEGIN
			c := DBMS_SQL.OPEN_CURSOR;
				DBMS_SQL.PARSE(c, 'SELECT * FROM %s WHERE ROWNUM = 1', DBMS_SQL.NATIVE);
				DBMS_SQL.DESCRIBE_COLUMNS(c, c_count, t_descriptor);
				FOR i IN 1..c_count LOOP
					DBMS_OUTPUT.PUT_LINE(i || ',' || t_descriptor(i).col_type);
				END LOOP;
			DBMS_SQL.CLOSE_CURSOR(c);
			EXCEPTION
				WHEN OTHERS THEN
					DBMS_SQL.CLOSE_CURSOR(c);
					RAISE;
		END;
	`, sanitize(tableName))

	output, err := getOutputLines(pool, statement)
	if err != nil {
		return nil, err
	}

	str_rdr := strings.NewReader(output)
	csv_rdr := csv.NewReader(str_rdr)

	recs, err := csv_rdr.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, rec := range recs {
		key := rec[0]

		val, err := strconv.ParseUint(rec[1], 10, 64)
		if err != nil {
			return nil, err
		}

		oids[key] = uint32(val)
	}

	return oids, nil
}
