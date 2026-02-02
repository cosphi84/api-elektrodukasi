package repositories

import (
	"api-elektrodukasi/models"
	"database/sql"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}
func (r *ArticleRepository) GetAll() ([]models.Article, error) {
	query := "SELECT a.*, b.name AS Author, c.name as Category FROM articles a LEFT JOIN authors b ON b.id = a.author_id LEFT JOIN categories c ON c.id = a.category_id"
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]models.Article, 0)
	for rows.Next() {
		var article models.Article
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Slug,
			&article.Author,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}
