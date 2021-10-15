package controllers

import (
	"ambassador/src/database"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)


func Rankings(c *fiber.Ctx) error {
	rankings, err := database.Cache.ZRevRangeByScoreWithScores(context.Background(), "rankings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

	if err != nil {
		return err
	}

	result := make(map[string]float64)

	for _, ranking := range rankings {
		result[ranking.Member.(string)] = ranking.Score
	}

	return c.JSON(result)
}
