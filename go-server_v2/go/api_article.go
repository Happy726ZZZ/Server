/*
 * Swagger Blog
 *
 * A Simple Blog
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"encoding/json"
	//"fmt"
    "log"
	"net/http"
	"net/url"
    "strings"
	"errors"
	"strconv"
	"github.com/boltdb/bolt"
)

func GetArticles(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("test.db", 0600, nil);
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	u, err := url.Parse(r.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	m, _ := url.ParseQuery(u.RawQuery)
	page := m["page"][0]
	IdIndex, err:= strconv.Atoi(page)
	IdIndex = (IdIndex - 1) * 5 +1
	
	var articles ArticlesResponse
	var article ArticleResponse

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("articles"))
		v := b.Get([]byte(strconv.Itoa(IdIndex)))
		if v == nil {
			return errors.New("Page is out of index")
		}

		count := 0
	    for v != nil && count < 5 {
		    err = json.Unmarshal(v, &article)
		    if err != nil {
			   return err
		    }
		    articles.Articles = append(articles.Articles, article)
			count = count + 1
			IdIndex = IdIndex + 1
			v = b.Get([]byte(strconv.Itoa(IdIndex)))
	    }
	    return nil
	})
	if err != nil {
		reponse := InlineResponse404{err.Error()}
		JsonResponse(reponse, w, http.StatusNotFound)
		return
	}
	JsonResponse(articles, w, http.StatusOK)
}

func GetArticleById(w http.ResponseWriter, r *http.Request) {

	db, err := bolt.Open("test.db", 0600, nil);
	if err != nil {
			log.Fatal(err)
	}
	defer db.Close()

	
	articleId := strings.Split(r.URL.Path, "/")[3]
	_, err = strconv.Atoi(articleId)
	if err != nil {
		reponse := InlineResponse404{"Wrong ArticleId"}
		JsonResponse(reponse, w, http.StatusBadRequest)
		return
	}
	var article Article

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("articles"))
		v := b.Get([]byte(articleId))

		if v == nil {
			reponse := InlineResponse404{"Article Not Exists"}
			JsonResponse(reponse, w, http.StatusNotFound)
			return nil
		} else {
			_ = json.Unmarshal(v, &article)
		}

		return nil
	})

	JsonResponse(article, w, http.StatusOK)
}

func GetCommentsOfArticle(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("test.db", 0600, nil);
	if err != nil {
			log.Fatal(err)
	}
	defer db.Close()

	articleId := strings.Split(r.URL.Path, "/")[3]
	Id, err := strconv.Atoi(articleId)
	if err != nil {
		reponse := InlineResponse404{"Wrong ArticleId"}
		JsonResponse(reponse, w, http.StatusBadRequest)
		return
	}

	var article []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("articles"))
		if b != nil {
			v := b.Get([]byte(strconv.Itoa(Id)))
			if v == nil {
				return errors.New("Article Not Exists1")
			} else {
				article = v
				return nil
			}
		} else {
			return errors.New("Article Not Exists")
		}
	})

	if err != nil {
		reponse := InlineResponse404{err.Error()}
		JsonResponse(reponse, w, http.StatusNotFound)
		return
	}
	var comments Comments
	var comment Comment
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("comments"))
		if b != nil {
			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				err = json.Unmarshal(v, &comment)
				if err != nil {
					return err
				}
				if comment.ArticleId == int32(Id) {
					comments.Contents = append(comments.Contents, comment)
				}
			}

			return nil
		} else {
			return errors.New("Comment Not Exists")
		}
	})

	if err != nil {
		reponse := InlineResponse404{err.Error()}
		JsonResponse(reponse, w, http.StatusNotFound)
		return
	}

	JsonResponse(comments, w, http.StatusOK)
}
