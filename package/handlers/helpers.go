package handlers

import (
	"errors"
)

type Errors struct {
	ErrorMessage error
	CheckErr     bool
	Cng          string
}

var parse Errors

func CheckLogin(str string) error {
	UppercaseLetters := false
	LowercaseLetters := false
	Numbers := false
	Only := false
	for _, i := range str {
		if i == 32 {
			return errors.New("Please dont use space for login")
		}
		if i >= 65 && i <= 90 {
			UppercaseLetters = true
		} else if i >= 97 && i <= 122 {
			LowercaseLetters = true
		} else if i >= 48 && i <= 57 {
			Numbers = true
		} else {
			Only = true
		}

	}
	if !UppercaseLetters {
		return errors.New("Please use Uppercase letters for login")
	}
	if !LowercaseLetters {
		return errors.New("Please use Lowercase letters for login")
	}
	if !Numbers {
		return errors.New("Please use Numbers for login")
	}
	if Only {
		return errors.New("Please use only Uppercase,Lowercase letters and Numbers for login")
	}
	return nil
}

func CheckPass(str string) error {
	UppercaseLetters := false
	LowercaseLetters := false
	Numbers := false
	Only := false
	for _, i := range str {
		if i == 32 {
			return errors.New("Please dont use space for password")
		}
		if i >= 65 && i <= 90 {
			UppercaseLetters = true
		} else if i >= 97 && i <= 122 {
			LowercaseLetters = true
		} else if i >= 48 && i <= 57 {
			Numbers = true
		} else {
			Only = true
		}

	}
	if !UppercaseLetters {
		return errors.New("Please use Uppercase letters for password")
	}
	if !LowercaseLetters {
		return errors.New("Please use Lowercase letters for password")
	}
	if !Numbers {
		return errors.New("Please use Numbers for password")
	}
	if Only {
		return errors.New("Please use only Uppercase,Lowercase letters and Numbers for password")
	}
	return nil
}

func CheckPost(title, text string, tags []string) error {
	mytags := []string{"Twitch", "YouTube", "GitHub", "Anime", "IT"}
	if len(title) < 1 {
		return errors.New("title length is less than 1")
	} else if len(title) > 40 {
		return errors.New("title length more than 40 characters")
	}
	if len(text) < 1 {
		return errors.New("text length is less than 1")
	} else if len(text) > 250 {
		return errors.New("text length more than 250 characters")
	}
	var space float64
	for _, i := range title {
		if i == 32 {
			space++
		}
	}
	if space > float64(len(title))/2 {
		space = 0
		return errors.New("too many spaces in title")
	}
	space = 0
	for _, j := range text {
		if j == 32 {
			space++
		}
	}
	if space > float64(len(text))/2 {
		space = 0
		return errors.New("too many spaces in post content")
	}
	if len(tags) == 0 {
		return errors.New("please select at least 1 tag ")
	}
	check := 0
	for _, i := range tags {
		if len(i) == 0 {
			return errors.New("please select at least 1 tag ")
		}
	}
	for k := range tags {
		for j := range mytags {
			if tags[k] == mytags[j] {
				check++
			}
		}
		if check == 0 {
			return errors.New("unknown tag")
		}
		check = 0
	}
	for q := 0; q < len(tags); q++ {
		for w := q + 1; w < len(tags); w++ {
			if tags[q] == tags[w] {
				return errors.New("tags are duplicated ")
			}
		}
	}
	return nil
}

func CheckComment(text string) error {
	if len(text) < 1 {
		return errors.New("comment length is less than 1")
	} else if len(text) > 100 {
		return errors.New("comment length more than 100 characters")
	}
	var space float64
	for _, j := range text {
		if j == 32 {
			space++
		}
	}
	if space > float64(len(text))/2 {
		space = 0
		return errors.New("too many spaces in comment")
	}
	return nil
}

func CheckEmptyLogin(login, password string) error {
	if len(login) == 0 && len(password) == 0 {
		return errors.New("empty login and password")
	} else if len(login) == 0 {
		return errors.New("empty login")
	} else if len(password) == 0 {
		return errors.New("empry password")
	}
	return nil
}

func CheckRate(numbL, numbD int) string {
	if numbL == 1 && numbD == 0 {
		return "like"
	} else if numbL == 0 && numbD == 1 {
		return "dislike"
	} else if numbL == -1 && numbD == 0 {
		return "deleted rating"
	} else if numbL == 0 && numbD == -1 {
		return "deleted rating"
	} else if numbL == 1 && numbD == -1 {
		return "changed rating to like "
	} else if numbL == -1 && numbD == 1 {
		return "changed rating to dislike "
	}
	return ""
}
