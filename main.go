package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	unsafeExts := hashset.New(".go", ".mod", ".sum")
	defaultMime := "application/octet-stream"
	mimeMap := map[string]string{
		".html": "text/html",
		".css":  "text/css",
		".js":   "text/javascript",
	}

	app.Get("/*", func(c *fiber.Ctx) error {
		name := c.Params("*")
		lastDot := strings.LastIndex(name, ".")
		fileExt := ""
		if lastDot >= 0 {
			fileExt = strings.ToLower(name[lastDot:])
		}

		if name == "" {
			name = "test.html"
		} else {

			if lastDot < 0 || unsafeExts.Contains(fileExt) {
				name = "nonexist.html"
			}
		}

		file, err := os.ReadFile(name)

		if err != nil {
			c.Set("Content-Type", "text/html")
			fmt.Println(err)
			html404, err404 := os.ReadFile("err404.html")
			if err404 != nil {
				return c.Status(fiber.StatusNotFound).SendString("<html>404 Not Found</html>")
			}
			return c.Status(fiber.StatusNotFound).Send(html404)
		}

		if mimeMap[fileExt] != "" {
			c.Set("Content-Type", mimeMap[fileExt])
		} else {
			c.Set("Content-Type", defaultMime)
		}
		return c.Send(file)
	})
	app.Listen(":3000")
}
