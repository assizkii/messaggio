package storages

import (
	"github.com/assizkii/messaggio/internal/domain/interfaces"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type PgStorage struct {
	connection *sqlx.DB
}

func New(dsn string) interfaces.MessageStorage {

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	err = Migrate(db)
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	return &PgStorage{db}
}

func (pg *PgStorage) Add(m *interfaces.Message) (int, error) {

	query := `insert into messages(phone, text)
				 values($1, $2) RETURNING id`

	var id int
	err := pg.connection.QueryRow(query, m.Phone, m.Text).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// sql migrations from migrations folder
func Migrate(db *sqlx.DB) error {

	err := filepath.Walk("./migrations", func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if !info.IsDir() {
			log.Println(path)
			fileData, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			fileValue := string(fileData)
			_, err = db.Exec(fileValue)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
