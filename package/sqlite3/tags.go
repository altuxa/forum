package sqlite3

import (
	"database/sql"
	"forum/package/models"
)

func AddTagsToDB(db *sql.DB, tags []string, postId int64) error {
	stmt, err := db.Prepare("INSERT INTO Tags (PostId, Value)VALUES(?,?)")
	if err != nil {
		return err
	}
	for i := 0; i < len(tags); i++ {
		_, err = stmt.Exec(postId, tags[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllTags(db *sql.DB, postId int) ([]models.Tags, error) {
	tags := []models.Tags{}
	row, err := db.Query("SELECT Id,PostId,Value FROM Tags WHERE PostId = ?", postId)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		tag := models.Tags{}
		err := row.Scan(&tag.Id, &tag.PostId, &tag.Value)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func DeleteTags(db *sql.DB, postId int) error {
	_, err := db.Exec("DELETE FROM Tags WHERE PostId = ?", postId)
	if err != nil {
		return err
	}
	return nil
}

func GetTagsByTagsName(db *sql.DB, tag string) ([]models.Tags, error) {
	tags := []models.Tags{}
	row, err := db.Query("SELECT Id,PostId, Value FROM Tags WHERE Value = ?", tag)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		tag := models.Tags{}
		err = row.Scan(&tag.Id, &tag.PostId, &tag.Value)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
