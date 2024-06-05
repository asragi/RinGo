package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ConnectionSettings struct {
	UserName string
	Password string
	Port     string
	Protocol string
	Host     string
	Database string
}

func ConnectDB(settings *ConnectionSettings) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s",
		settings.UserName,
		settings.Password,
		settings.Protocol,
		settings.Host,
		settings.Port,
		settings.Database,
	)
	dataSource = fmt.Sprintf("%s?parseTime=True", dataSource)
	db, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type (
	DBAccessor struct {
		writer *sqlx.DB
		reader *sqlx.DB
	}

	abstractSqlxDB interface {
		BindNamed(query string, arg interface{}) (q string, args []interface{}, err error)
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	}

	ctxKeyTx struct{}
)

func NewDBAccessor(
	writer *sqlx.DB,
	reader *sqlx.DB,
) *DBAccessor {
	return &DBAccessor{
		writer: writer,
		reader: reader,
	}
}

func withTxDB(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, &ctxKeyTx{}, tx)
}

func getTxFromContext(ctx context.Context) (*sqlx.Tx, error) {
	if v := ctx.Value(&ctxKeyTx{}); v != nil {
		tx, ok := v.(*sqlx.Tx)
		if !ok {
			return nil, fmt.Errorf("failed to convert to *sqlx.Tx: %v", v)
		}

		return tx, nil
	}

	return nil, nil
}

func (dba *DBAccessor) Close() error {
	writerErr := dba.writer.Close()
	if dba.writer == dba.reader {
		return writerErr
	}
	readerErr := dba.reader.Close()
	if writerErr != nil {
		if readerErr != nil {
			return fmt.Errorf("writerErr: %w, readerErr: %w", writerErr, readerErr)
		}
		return writerErr
	}
	if readerErr != nil {
		return readerErr
	}
	return nil
}

func (dba *DBAccessor) Transaction(ctx context.Context, txFunc func(context.Context) error) (err error) {
	tx, err := dba.writer.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			if err := tx.Rollback(); err != nil {
				// ログ出力
			}
			err = fmt.Errorf("panicked on execution txFunc: %v", r)
		}
	}()

	txCtx := withTxDB(ctx, tx)

	if txErr := txFunc(txCtx); txErr != nil {
		if err := tx.Rollback(); err != nil {
			// ログ出力
		}

		return txErr
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type ExecFunc func(
	ctx context.Context,
	query string,
	arg interface{},
) (sql.Result, error)

func (dba *DBAccessor) Exec(
	ctx context.Context,
	query string,
	arg interface{},
) (sql.Result, error) {
	tx, err := getTxFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to getTxFromContext: %w", err)
	}

	var da abstractSqlxDB
	if tx != nil {
		da = tx
	} else {
		da = dba.writer
	}

	if arg == nil {
		result, err := da.ExecContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("failed to run exec: %w", err)
		}
		return result, nil
	}

	// prepare
	namedQuery, namedArgs, err := da.BindNamed(query, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to BindNamed: %w", err)
	}
	namedQuery, namedArgs, err = sqlx.In(namedQuery, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to sqlx.In: %w", err)
	}

	result, err := da.ExecContext(ctx, namedQuery, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to run exec: %w", err)
	}

	return result, nil
}

type QueryFunc func(
	ctx context.Context,
	query string,
	arg interface{},
) (*sqlx.Rows, error)

func (dba *DBAccessor) Query(
	ctx context.Context,
	query string,
	arg interface{},
) (*sqlx.Rows, error) {
	tx, err := getTxFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to getTxFromContext: %w", err)
	}

	var da abstractSqlxDB
	if tx != nil {
		da = tx
	} else {
		da = dba.reader
	}

	if arg == nil {
		rows, err := da.QueryxContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("failed to run query: %w", err)
		}
		return rows, nil
	}

	// prepare
	namedQuery, namedArgs, err := da.BindNamed(query, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to BindNamed: %w", err)
	}
	// namedQuery, namedArgs, err = sqlx.In(namedQuery, namedArgs...)
	namedQuery, namedArgs, err = sqlx.In(namedQuery, namedArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to sqlx.In: %w", err)
	}

	rows, err := da.QueryxContext(ctx, namedQuery, namedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}

	return rows, nil
}
