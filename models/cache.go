pacakge models

var ctx = context.Background()

func init()  {
	rdb := redis.NewClient(
		&redis.Options{
			Addr :"localhost:6379",
			Password :"",
			DB : 0

		}
	)
	pong, err := redis.Ping(ctx).Result()
	fmt.Println(pong)
}