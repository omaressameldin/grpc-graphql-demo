package db

import (
	"encoding/binary"
	"time"

	"../model"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func AllTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			task, err := model.BufferToTask(v)
			if err != nil {
				return err
			}
			tasks = append(tasks, task)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func ReadTask(id int64) (model.Task, error) {
	var task model.Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		buf := b.Get(itob(id))
		var err error
		task, err = model.BufferToTask(buf)
		return err
	})

	return task, err
}

func CreateTask(task model.Task) (int64, error) {
	var id64 int64
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		idU64, _ := b.NextSequence()
		id64 = int64(idU64)
		return putTask(id64, task, b)
	})
	if err != nil {
		return -1, err
	}

	return id64, nil
}

func UpdateTask(id int64, t model.Task) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return putTask(id, t, b)
	})
}

func DeleteTask(id int64) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(id))
	})
}

func itob(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func putTask(id64 int64, task model.Task, b *bolt.Bucket) error {
	id := itob(id64)
	taskBuffer, err := model.TaskToBuffer(task)
	if err != nil {
		return err
	}
	return b.Put(id, taskBuffer)
}
