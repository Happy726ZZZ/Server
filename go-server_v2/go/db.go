package swagger

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"log"
	"github.com/boltdb/bolt"
	"strconv"
	
)

func createdb() {
	//1.create database
	db, err := bolt.Open("test.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
	}
	defer db.Close()
	//2.create table
    err = db.Update(func(tx *bolt.Tx) error {
		//if the table exists
		ArticlesTable := tx.Bucket([]byte("articles"))
		//create article table
		_, err := tx.CreateBucket([]byte("articles"))
		if err != nil {
			//insert data
			file_name := [8]string{
				"【LeetCode】128. Longest Consecutive Sequence",
				"【LeetCode】32. Longest Valid Parentheses",
				"【LeetCode】679. 24 Game",
				"【LeetCode】23. Merge k Sorted Lists",
				"【LeetCode】785. Is Graph Bipartite?",
				"【LeetCode】105&106、根据前/后序遍历和中序遍历还原二叉树",
				"【LeetCode】210. Course Schedule II",
				"【LeetCode】685. Redundant Connection II",
			}
			file_date := [8]string{
				"2018-12-16 13:07:57",
				"2018-12-16 13:09:01",
				"2018-12-16 13:13:11",
				"2018-12-16 13:14:03",
				"2018-12-16 13:06:07",
				"2018-12-16 13:15:09",
				"2018-12-16 13:16:10",
				"2018-12-16 13:17:26",
			}

			for n := 0; n < 8; n++ {
				file_path := "Server/go-server-v2/resource/" + file_name[n] +".md"
				file, _ := os.Open(file_path)
				defer file.Close()
				content, _ := ioutil.ReadAll(file)
				it := &Article {
					Id: int32(n+1),
					Name: file_name[n],
					Date: file_date[n],
					Content: string(content),
				}
				jsons, _ := json.Marshal(it)
				ArticlesTable.Put([]byte(strconv.Itoa(int(it.Id))), jsons)
			}
		}

		//if the table exists
		CommentsTable := tx.Bucket([]byte("comments"))
		if CommentsTable == nil {
		//create users table
			_, err := tx.CreateBucket([]byte("comments"))
			if err != nil {
				//insert data
				log.Fatal(err)
			}
		}

		//if the table exists
		UsersTable := tx.Bucket([]byte("users"))
		if UsersTable == nil {
		//create users table
			_, err := tx.CreateBucket([]byte("users"))
			if err != nil {
				//insert data
				log.Fatal(err)
			}
		}
		return nil
	})

	//create table faile
	if err != nil {
		log.Fatal(err)
	}
}
