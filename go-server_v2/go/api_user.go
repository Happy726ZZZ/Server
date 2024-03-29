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
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"

	//"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	//"github.com/dgrijalva/jwt-go/request"
)

func ByteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func JsonResponse(response interface{}, w http.ResponseWriter, code int) {
	json, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(json)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	articleId := strings.Split(r.URL.Path, "/")[3]
	Id, err := strconv.Atoi(articleId)
	if err != nil {
		response := InlineResponse404{"Wrong ArticleId"}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("articles"))
		if b != nil {
			v := b.Get([]byte(strconv.Itoa(Id)))
			if v == nil {
				return errors.New("Article Not Exists")
			} else {
				return nil
			}
		}
		return errors.New("Article Not Exists")
	})

	if err != nil {
		response := InlineResponse404{err.Error()}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}

	comment := &Comment{
		Date:      time.Now().Format("2006-01-02 15:04:05"),
		Content:   "",
		Author:    "",
		ArticleId: int32(Id),
	}
	err = json.NewDecoder(r.Body).Decode(&comment)

	if err != nil || comment.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err != nil {
			response := InlineResponse404{err.Error()}
			JsonResponse(response, w, http.StatusBadRequest)
		} else {
			response := InlineResponse404{"There is no content in your article"}
			JsonResponse(response, w, http.StatusBadRequest)
		}
		return
	}

	//token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
	//    func(token *jwt.Token) (interface{}, error) {
	//        return []byte(comment.Author), nil
	//    })

	//if err == nil {
	//if token.Valid {
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("comments"))
		if err != nil {
			return err
		}
		id, _ := b.NextSequence()
		encoded, err := json.Marshal(comment)
		return b.Put([]byte(strconv.Itoa(int(id))), encoded)
	})

	if err != nil {
		response := InlineResponse404{err.Error()}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}

	JsonResponse(comment, w, http.StatusOK)
	//    } else {
	//		response := InlineResponse404{"Token is not valid"}
	//		JsonResponse(response, w, http.StatusUnauthorized)
	//    }
	//} else {
	//	response := InlineResponse404{"Unauthorized access to this resource"}
	//	JsonResponse(response, w, http.StatusUnauthorized)
	//}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User

	err = json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)
	if err != nil {
		response := InlineResponse404{err.Error()}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b != nil {
			v := b.Get([]byte(user.Username))
			if ByteSliceEqual(v, []byte(user.Password)) {
				return nil
			} else {
				return errors.New("Wrong Username or Password")
			}
		} else {
			return errors.New("Wrong Username or Password")
		}
	})

	if err != nil {
		response := InlineResponse404{err.Error()}
		JsonResponse(response, w, http.StatusNotFound)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		log.Fatal(err)
	}

	tokenString, err := token.SignedString([]byte(user.Username))
	if err != nil {
		log.Fatal(err)
	}

	response := InlineResponse200{tokenString}
	JsonResponse(response, w, http.StatusOK)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	log.Printf("%s  %s", user.Username, user.Password)
	if err != nil || user.Password == "" || user.Username == "" {
		response := InlineResponse404{"Wrong Username or Password"}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b != nil {
			v := b.Get([]byte(user.Username))
			if v != nil {
				return errors.New("User Exists")
			}
		}
		return nil
	})

	if err != nil {
		response := InlineResponse404{err.Error()}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		return b.Put([]byte(user.Username), []byte(user.Password))
	})

	if err != nil {
		response := InlineResponse404{err.Error()}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With,Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}
