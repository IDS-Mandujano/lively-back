package mysql

import (
	"context"
	"database/sql"
	"lively-backend/src/radio/Rooms/domain/entity"
)

type MySQLRoomRepository struct {
	db *sql.DB
}

func NewMySQLRoomRepository(db *sql.DB) *MySQLRoomRepository {
	return &MySQLRoomRepository{
		db: db,
	}
}

func (repo *MySQLRoomRepository) Create(ctx context.Context, room *entity.Room) error {
	query := `
		INSERT INTO rooms (id, name, description, created_by) 
		VALUES (?, ?, ?, ?)`

	_, err := repo.db.ExecContext(ctx, query, room.ID, room.Name, room.Description, room.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MySQLRoomRepository) GetAll(ctx context.Context) ([]entity.Room, error) {
	query := "SELECT id, name, description, created_by, created_at FROM rooms ORDER BY created_at DESC"

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []entity.Room
	for rows.Next() {
		var r entity.Room
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedBy, &r.CreatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}

	return rooms, nil
}
