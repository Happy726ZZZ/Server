package swagger

import (
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
			it := &Article {
				Id: 1,
				Name: "【LeetCode】128. Longest Consecutive Sequence",
				Date: "2018-12-16 13:07:57",
				Content: "resource/【LeetCode】128. Longest Consecutive Sequence.md",
			}
			jsons, _ := json.Marshal(it)
			ArticlesTable.Put([]byte(strconv.Itoa(int(it.Id))), jsons)
		
			it = &Article {
				Id: 2,
				Name: "【LeetCode】32. Longest Valid Parentheses",
				Date: "2018-12-16 13:09:01",
				Content: "resource/【LeetCode】32. Longest Valid Parentheses.md",
			}
			jsons, _ = json.Marshal(it)
			ArticlesTable.Put([]byte(strconv.Itoa(int(it.Id))), jsons)
				
			it = &Article {
				Id: 3,
				Name: "【LeetCode】679. 24 Game",
				Date: "2018-12-16 13:13:11",
				Content: "resource/【LeetCode】679. 24 Game.md",
			}
			jsons, _ = json.Marshal(it)
			ArticlesTable.Put([]byte(strconv.Itoa(int(it.Id))), jsons)

			it = &Article {
				Id: 4,
				Name: "【LeetCode】23. Merge k Sorted Lists",
				Date: "2018-12-16 13:14:03",
				Content: "resource/【LeetCode】23. Merge k Sorted Lists.md",
			}
			jsons, _ = json.Marshal(it)
			ArticlesTable.Put([]byte(strconv.Itoa(int(it.Id))), jsons)

			it = &Article {
				Id: 5,
				Name: "【LeetCode】785. Is Graph Bipartite?",
				Date: "2018-12-16 13:06:07",
				Content: "resource/【LeetCode】785. Is Graph Bipartite?.md",
			}
			jsons, _ = json.Marshal(it)
			ArticlesTable.Put([]byte(strconv.Itoa(int(it.Id))), jsons)

			it = &Article {
				Id: 6,
				Name: "【LeetCode】105&106、根据前/后序遍历和中序遍历还原二叉树",
				Date: "2018-12-16 13:15:09",
				Content: "resource/【LeetCode】105&106、根据前/后序遍历和中序遍历还原二叉树.md",
			}
			jsons, _ = json.Marshal(it)
			ArticlesTable.Put([]byte(strconv.Itoa(int(it.Id))), jsons)
			
			it = &Article {
				Id: 7,
				Name: "【LeetCode】210. Course Schedule II",
				Date: "2018-12-16 13:17:26",
				Content: "resource/【LeetCode】685. Redundant Connection II.md",
			}
			jsons, _ = json.Marshal(it)
			ArticlesTable.Put([]byte(strconv.Itoa(int(it.Id))), jsons)

			it = &Article {
				Id: 8,
				Name: "【LeetCode】685. Redundant Connection II",
				Date: "2018-12-16 13:17:26",
				Content: "resource/【LeetCode】215. Kth Largest Element in an Array.md",
			}
			jsons, _ = json.Marshal(it)
			ArticlesTable.Put([]byte(strconv.Itoa(int(it.Id))), jsons)
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