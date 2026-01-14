package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func main() {
	today := time.Now()
	pageUrl := ""
	if today.Day()%2 == 0 {
		pageUrl = "https://quotes.toscrape.com/page/2/"
	} else {
		pageUrl = "https://quotes.toscrape.com/"
	}

	fmt.Printf("📅 %s | 📄 %s\n", today.Format("Mon, 02 Jan 2006"), pageUrl)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var data []Quote
	err := chromedp.Run(ctx,
		chromedp.Navigate(pageUrl),
		chromedp.WaitVisible(".quote"),
		chromedp.EvaluateAsDevTools(`
			(() => {
				return [...document.querySelectorAll('.quote')]
					.slice(0, 10)
					.map(q => ({
						text: q.querySelector('.text').innerText,
						author: q.querySelector('.author').innerText
					}));
			})()
		`, &data),
	)
	if err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}

	outputDir := filepath.Join(".", "public")
	os.MkdirAll(outputDir, os.ModePerm)
	outputPath := filepath.Join(outputDir, "data.json")

	file, _ := os.Create(outputPath)
	defer file.Close()
	json.NewEncoder(file).Encode(data)

	fmt.Printf("✅ Data scraped and saved to %s\n", outputPath)
}