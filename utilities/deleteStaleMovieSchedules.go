package utilities

import (
	"context"
	"ticket/api/dbInstance"
	"time"
)

func DeleteStaleMovieSchedules(){
	store := dbInstance.Store;
	
	for{
		time.Sleep(time.Hour*24);
		timeNow := time.Now().Unix();
		store.Queries.DeleteStaleMovieSchedules(context.Background(), timeNow);
	}
}

