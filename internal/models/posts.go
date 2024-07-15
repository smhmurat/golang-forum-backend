package models

import (
	"database/sql"
	"github.com/upper/db/v4"
	"strings"
)

type Post struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Categories    []string `json:"categories"`
	UserID        int      `json:"user_id"`
	Likes         int      `json:"likes"`
	Dislikes      int      `json:"dislikes"`
	CommentLength int      `json:"comment_length"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	TotalRecords  int      `json:"-"`
}

type PostModel struct {
	db db.Session
}

var (
	queryTemplate = `
	SELECT COUNT(*) OVER() AS total_records, pq.*, u.name as uname FROM (
	    SELECT p.id, p.title, p.url, p.created_at, p.user_id as uid, COUNT(c.post_id) as comment_count, count(v.post_id) as votes
		FROM posts p 
		LEFT JOIN comments c ON p.id = c.post_id 
	    LEFT JOIN votes v ON p.id = v.post_id
	 	#where#
		GROUP BY p.id
		#orderby#
		) AS pq
	LEFT JOIN users u ON u.id = uid
	#limit#
	`
)

func (m PostModel) GetAll(f Filter) ([]Post, Metadata, error) {
	//db := database.GetDB()
	var posts []Post
	var rows *sql.Rows
	var err error
	meta := Metadata{}

	q := f.applyTemplate(queryTemplate)

	if len(f.Query) > 0 {
		rows, err = m.db.SQL().Query(q, "%"+strings.ToLower(f.Query)+"%", f.limit(), f.offset())
	} else {
		rows, err = m.db.SQL().Query(q, f.limit(), f.offset())
	}

	if err != nil {
		return nil, meta, err
	}

	iter := m.db.SQL().NewIterator(rows)
	err = iter.All(&posts)
	if err != nil {
		return nil, meta, err
	}

	if len(posts) == 0 {
		return nil, meta, nil
	}

	first := posts[0]
	return posts, calculateMetadata(first.TotalRecords, f.Page, f.PageSize), nil

}
