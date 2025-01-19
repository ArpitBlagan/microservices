package redis

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)



type contextKey string
const userKey contextKey = "user"
var redisClient *redis.Client
var ctx =context.Background()
func InitRedis(){
	redisClient=redis.NewClient(
		&redis.Options{
			Addr:"",
			Password:"",
			DB:0,
		},
	)
	_, err:=redisClient.Ping(ctx).Result()
	if err!=nil{
		log.Fatal("Error while creating redis client :(")
	}

}

func RateLimitRequest(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,request *http.Request){
		userID, ok := request.Context().Value(userKey).(string)
		if !ok {
			http.Error(w, "User not found in context", http.StatusUnauthorized)
			return
		}
		key := fmt.Sprintf("rate_limit:%s", userID)
		maxRequest:=10
		duration:=time.Minute
		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			log.Printf("Error updating rate limit: %v", err)
			 http.Error(w,"Rate limit exceed please try again later :(",http.StatusTooManyRequests)
			 return
		}
		// Set expiration time if it's the first request
		if count == 1 {
			redisClient.Expire(ctx, key, duration)
		}

		 if count > int64(maxRequest){
		 	http.Error(w,"Rate limit exceed please try again later :(",http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, request)
	})
	
}


func SetCache(key string,value string, expirationTime time.Duration) error{
	
	return redis.Client.Set(*redisClient,ctx,key,value,expirationTime).Err();
}

func GetCache(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		userID, ok := request.Context().Value(userKey).(string)
		if !ok {
			http.Error(w, "User not found in context", http.StatusUnauthorized)
			return
		}
		val,err := redisClient.Get(ctx, userID).Result()
		if err== redis.Nil {
			next.ServeHTTP(w, request)
			return
		}
		if err != nil {
			next.ServeHTTP(w, request)
			return
		}
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(val))
	})
}