package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

func (s *Storage) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT 
		posts.id,
		posts.title,
		posts.content,
		authors.id,
		authors.name,
		posts.created_at,
		posts.published_at		
	FROM posts,authors
	WHERE authors.id=posts.author_id;
`)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
			&p.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

func (s *Storage) AddPost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	INSERT INTO authors (id, name) VALUES ($1, '$2');
	INSERT INTO posts (id, author_id, title, content, created_at,published_at) VALUES ($3, $1, '$4', '$5', $6, $7);
`,
		p.AuthorID,
		p.AuthorName,
		p.ID,
		p.Title,
		p.Content,
		p.CreatedAt,
		p.PublishedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	UPDATE posts
	SET 
		title = '$2',
		content = '$3',
		created_at = $4,
		published_at = $5
	WHERE id=$1;
`,
		p.ID,
		p.Title,
		p.Content,
		p.CreatedAt,
		p.PublishedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM posts
	WHERE
		id = $1;
`,
		p.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
