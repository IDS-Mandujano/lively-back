package mysql

import (
	"context"
	"database/sql"
	"lively-backend/src/radio/Rooms/domain/entity"
	"log"
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
		INSERT INTO rooms (name, description, created_by) 
		VALUES (?, ?, ?)`

	result, err := repo.db.ExecContext(ctx, query, room.Name, room.Description, room.CreatedBy)
	if err != nil {
		log.Printf("[rooms:repo] create failed created_by=%d name=%q err=%v", room.CreatedBy, room.Name, err)
		return err
	}

	insertedID, err := result.LastInsertId()
	if err == nil {
		room.ID = int(insertedID)
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
