package sqlite3

import (
	"database/sql"
	"forum/package/models"
	"log"
)

func AddPostToDB(db *sql.DB, post models.Posts) error {
	stmt, err := db.Prepare("INSERT INTO Posts(UserId, Title,Text,Author)VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(post.UserId, post.Title, post.Text, post.Author)
	if err != nil {
		return err
	}
	postId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	err = AddTagsToDB(db, post.Tags, postId)
	if err != nil {
		return err
	}
	return nil
}

func SelectPostsFromDB(db *sql.DB) ([]models.Posts, error) {
	posts := []models.Posts{}
	row, err := db.Query("SELECT Id,UserId, Title,Text,Author FROM Posts ORDER BY Id DESC")
	if err != nil {
		return nil, err
	}
	for row.Next() {
		post := models.Posts{}
		err := row.Scan(&post.Id, &post.UserId, &post.Title, &post.Text, &post.Author)
		if err != nil {
			log.Println(err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func DeletePost(db *sql.DB, postId int) error {
	_, err := db.Exec("DELETE FROM Posts WHERE Id = ?", postId)
	if err != nil {
		return err
	}
	return nil
}

func GetOnePost(db *sql.DB, postId int) (models.Posts, error) {
	var id, userId, likes, dislikes int
	var title, text, author string
	row := db.QueryRow("SELECT Id, UserId, Title, Text, Author, Likes, Dislikes FROM Posts WHERE Id = ?", postId)
	row.Scan(&id, &userId, &title, &text, &author, &likes, &dislikes)
	comments, err := SelectCommentsFromDB(db, postId)
	if err != nil {
		return models.Posts{}, err
	}
	tags, err := GetAllTags(db, postId)
	if err != nil {
		return models.Posts{}, err
	}
	post := models.Posts{
		Id:       id,
		UserId:   userId,
		Title:    title,
		Text:     text,
		Author:   author,
		Comment:  comments,
		Tag:      tags,
		Likes:    likes,
		Dislikes: dislikes,
	}
	return post, nil
}

func UpdatePost(db *sql.DB, title, text string, postId, userId int) error {
	_, err := db.Exec("UPDATE Posts Set Title = ?,Text = ? WHERE Id = ? AND UserId = ?", title, text, postId, userId)
	if err != nil {
		return err
	}
	return nil
}

func GetPostsByUserId(db *sql.DB, id int) ([]models.Posts, error) {
	posts := []models.Posts{}
	row, err := db.Query("SELECT Id, UserId,Title, Text, Author FROM Posts WHERE UserId = ? ORDER BY Id DESC", id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		post := models.Posts{}
		// МБ УБРАТЬ ЕРРОР ПОДУМАЮ
		err = row.Scan(&post.Id, &post.UserId, &post.Title, &post.Text, &post.Author)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsByPostId(db *sql.DB, postId int) ([]models.Posts, error) {
	posts := []models.Posts{}
	row, err := db.Query("SELECT Id, UserId, Title,Text, Author FROM Posts Where Id = ?", postId)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		post := models.Posts{}
		// TOJE SAMOE
		err = row.Scan(&post.Id, &post.UserId, &post.Title, &post.Text, &post.Author)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func AddRatingToPost(db *sql.DB, like, dislike, postId int) error {
	var allLike, allDislike int
	row := db.QueryRow("SELECT Likes,Dislikes FROM Posts WHERE Id = ?", postId)
	row.Scan(&allLike, &allDislike)
	stmt, err := db.Prepare("UPDATE Posts SET Likes = ?, Dislikes = ? WHERE Id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(allLike+like, allDislike+dislike, postId)
	if err != nil {
		return err
	}
	return nil
}

func AddLikeOrDislikeToRatingPost(db *sql.DB, userId, postId, like, dislike int) error {
	stmt, err := db.Prepare("INSERT INTO RatingPost (PostId,UserId,Like,Dislike)VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(postId, userId, like, dislike)
	if err != nil {
		return err
	}
	return nil
}

func UpdateLikeOrDislikeToRatingPost(db *sql.DB, userId, postId, like, dislike int) error {
	stmt, err := db.Prepare("UPDATE RatingPost SET Like = ?, Dislike = ? WHERE PostId = ? AND UserId = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(like, dislike, postId, userId)
	if err != nil {
		return err
	}
	return nil
}

func CheckRatingPostTable(db *sql.DB, userId, postId int) (bool, models.RatingPost) {
	ratingPost := models.RatingPost{}
	row := db.QueryRow("SELECT PostId,UserId,Like,Dislike FROM RatingPost WHERE PostId = ? AND UserId = ?", postId, userId)
	err := row.Scan(&ratingPost.PostId, &ratingPost.UserId, &ratingPost.Like, &ratingPost.Dislike)
	if err != nil {
		return false, ratingPost
	}
	return true, ratingPost
}

func DeleteRatingByPostId(db *sql.DB, postId int) error {
	_, err := db.Exec("DELETE FROM RatingPost WHERE PostId = ?", postId)
	if err != nil {
		return err
	}
	return nil
}

func GetMyFavoritePosts(db *sql.DB, userId int) ([]models.RatingPost, error) {
	ratingPosts := []models.RatingPost{}
	row, err := db.Query("SELECT PostId, UserId,Like FROM RatingPost WHERE UserId = ?", userId)
	if err != nil {
		// POTESTIT
		row.Next()
	}
	for row.Next() {
		postRating := models.RatingPost{}
		row.Scan(&postRating.PostId, &postRating.UserId, &postRating.Like)
		if postRating.Like != 0 {
			ratingPosts = append(ratingPosts, postRating)
		}
	}
	return ratingPosts, nil
}

func GetMyNotFavoritePost(db *sql.DB, userId int) ([]models.RatingPost, error) {
	ratingPosts := []models.RatingPost{}
	row, err := db.Query("SELECT PostId, UserId,Dislike FROM RatingPost WHERE UserId = ?", userId)
	if err != nil {
		// POTESTIT
		row.Next()
	}
	for row.Next() {
		postRating := models.RatingPost{}
		row.Scan(&postRating.PostId, &postRating.UserId, &postRating.Dislike)
		if postRating.Dislike != 0 {
			ratingPosts = append(ratingPosts, postRating)
		}
	}
	return ratingPosts, nil
}
