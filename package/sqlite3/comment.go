package sqlite3

import (
	"database/sql"
	"forum/package/models"
)

func AddCommentToDB(db *sql.DB, com models.Comments) (error, int64) {
	stmt, err := db.Prepare("INSERT INTO Comments(UserId, PostId, Text, Author)VALUES(?,?,?,?)")
	if err != nil {
		return err, 0
	}
	res, err := stmt.Exec(com.UserId, com.PostId, com.Text, com.Author)
	if err != nil {
		return err, 0
	}
	comId, err := res.LastInsertId()
	return nil, comId
}

func SelectCommentsFromDB(db *sql.DB, postId int) ([]models.Comments, error) {
	comment := []models.Comments{}
	row, err := db.Query("SELECT Id, UserId,PostId,Text, Author, Likes, Dislikes FROM Comments WHERE PostId = ? ORDER BY Id DESC", postId)
	if err != nil {
		return comment, nil
	}
	for row.Next() {
		com := models.Comments{}
		row.Scan(&com.Id, &com.UserId, &com.PostId, &com.Text, &com.Author, &com.Likes, &com.Dislikes)
		comment = append(comment, com)
	}
	return comment, nil
}

func DeleteCommentFromDB(db *sql.DB, Id, UserId int) error {
	_, err := db.Exec("DELETE FROM Comments WHERE Id = ? AND UserId = ?", Id, UserId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateComment(db *sql.DB, text string, comId, userId int) error {
	_, err := db.Exec("UPDATE Comments Set Text = ? WHERE Id = ? AND UserId = ?", text, comId, userId)
	if err != nil {
		return err
	}
	return nil
}

func GetOneCommentfromDB(db *sql.DB, id, UserId int) models.Comments {
	comment := models.Comments{}
	row := db.QueryRow("SELECT Id, UserId, PostId, Text,Author FROM Comments WHERE Id = ? AND UserId = ?", id, UserId)
	row.Scan(&comment.Id, &comment.UserId, &comment.PostId, &comment.Text, &comment.Author)
	return comment
}

func DeleteAllCommentFromPost(db *sql.DB, postId int) error {
	_, err := db.Exec("DELETE FROM Comments WHERE PostId = ?", postId)
	if err != nil {
		return err
	}
	return nil
}

func GetPostIdByCommentId(db *sql.DB, id int) (int, error) {
	var postId int
	row := db.QueryRow("SELECT PostId FROM Comments WHERE Id = ?", id)
	err := row.Scan(&postId)
	if err != nil {
		return 0, err
	}
	return postId, nil
}

func AddRatingToComment(db *sql.DB, like, dislike, comId int) error {
	var allLike, allDislike int
	row := db.QueryRow("SELECT Likes,Dislikes FROM Comments WHERE Id = ?", comId)
	row.Scan(&allLike, &allDislike)
	stmt, err := db.Prepare("UPDATE Comments SET Likes = ?, Dislikes = ? WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(allLike+like, allDislike+dislike, comId)
	if err != nil {
		return err
	}
	return nil
}

func AddLikeOrDislikeToRatingComment(db *sql.DB, userId, comId, like, dislike int) error {
	stmt, err := db.Prepare("INSERT INTO RatingComment (CommentId,UserId,Like,Dislike)VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(comId, userId, like, dislike)
	if err != nil {
		return err
	}
	return nil
}

func UpdateLikeOrDislikeToRatingComment(db *sql.DB, userId, comId, like, dislike int) error {
	stmt, err := db.Prepare("UPDATE RatingComment SET Like = ?, Dislike = ? WHERE CommentId = ? AND UserId = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(like, dislike, comId, userId)
	if err != nil {
		return err
	}
	return nil
}

func CheckRatingCommentTable(db *sql.DB, userId, comId int) (bool, models.RatingComment) {
	ratingComment := models.RatingComment{}
	row := db.QueryRow("SELECT CommentId,UserId,Like,Dislike FROM RatingComment WHERE CommentId = ? AND UserId = ?", comId, userId)
	err := row.Scan(&ratingComment.CommentId, &ratingComment.UserId, &ratingComment.Like, &ratingComment.Dislike)
	if err != nil {
		return false, ratingComment
	}
	return true, ratingComment
}

func DeleteRatingCommentByCommentId(db *sql.DB, comId int) error {
	_, err := db.Exec("DELETE FROM RatingComment WHERE CommentId = ?", comId)
	if err != nil {
		return err
	}
	return nil
}

func GetMyComments(db *sql.DB, userId int) ([]models.Comments, error) {
	allCom := []models.Comments{}
	row, err := db.Query("SELECT Id,UserId,PostId,Text,Author FROM Comments WHERE UserId = ? ORDER BY Id DESC", userId)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		oneCom := models.Comments{}
		row.Scan(&oneCom.Id, &oneCom.UserId, &oneCom.PostId, &oneCom.Text, &oneCom.Author)
		allCom = append(allCom, oneCom)
	}
	return allCom, nil
}
