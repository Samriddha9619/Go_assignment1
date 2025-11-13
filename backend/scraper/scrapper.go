package scraper
import (
    "context"
    "fmt"
    "log"
    "math/rand"
    "strings"
    "time"
    "github.com/Samriddha9619/Go_assignment1/backend/models"
    "github.com/chromedp/chromedp"
	"github.com/chromedp/cdproto/cdp"  
)
var mockHotels = map[string][]models.Hotel{
	"Delhi": {
		{Name: "The Oberoi New Delhi", Location: "Dr Zakir Hussain Marg", Price: 15000, Rating: 4.7, Source: "mock"},
		{Name: "The Leela Palace New Delhi", Location: "Diplomatic Enclave", Price: 18000, Rating: 4.8, Source: "mock"},
		{Name: "ITC Maurya", Location: "Sardar Patel Marg", Price: 12000, Rating: 4.6, Source: "mock"},
		{Name: "The Taj Mahal Hotel", Location: "Man Singh Road", Price: 14000, Rating: 4.7, Source: "mock"},
		{Name: "The Imperial", Location: "Janpath", Price: 16000, Rating: 4.8, Source: "mock"},
		{Name: "Radisson Blu Plaza Delhi", Location: "National Highway 8", Price: 8000, Rating: 4.3, Source: "mock"},
		{Name: "Le Meridien New Delhi", Location: "Windsor Place", Price: 11000, Rating: 4.5, Source: "mock"},
		{Name: "JW Marriott Hotel New Delhi", Location: "Asset Area 2", Price: 13000, Rating: 4.6, Source: "mock"},
		{Name: "Hyatt Regency Delhi", Location: "Bhikaji Cama Place", Price: 9500, Rating: 4.4, Source: "mock"},
		{Name: "The Lodhi", Location: "Lodhi Road", Price: 20000, Rating: 4.9, Source: "mock"},
	},
	"Mumbai": {
		{Name: "Taj Mahal Palace Mumbai", Location: "Apollo Bunder", Price: 22000, Rating: 4.8, Source: "mock"},
		{Name: "The Oberoi Mumbai", Location: "Nariman Point", Price: 18000, Rating: 4.7, Source: "mock"},
		{Name: "Trident Nariman Point", Location: "Nariman Point", Price: 12000, Rating: 4.5, Source: "mock"},
		{Name: "ITC Grand Central", Location: "Parel", Price: 11000, Rating: 4.6, Source: "mock"},
		{Name: "JW Marriott Mumbai Juhu", Location: "Juhu Beach", Price: 15000, Rating: 4.7, Source: "mock"},
		{Name: "The St. Regis Mumbai", Location: "Lower Parel", Price: 25000, Rating: 4.9, Source: "mock"},
		{Name: "Four Seasons Mumbai", Location: "Worli", Price: 24000, Rating: 4.8, Source: "mock"},
		{Name: "Sofitel Mumbai BKC", Location: "Bandra Kurla Complex", Price: 13000, Rating: 4.6, Source: "mock"},
		{Name: "Hyatt Regency Mumbai", Location: "Sahar Airport Road", Price: 10000, Rating: 4.4, Source: "mock"},
		{Name: "Renaissance Mumbai", Location: "Powai", Price: 9500, Rating: 4.5, Source: "mock"},
	},
	"Bangalore": {
		{Name: "The Leela Palace Bengaluru", Location: "Old Airport Road", Price: 16000, Rating: 4.8, Source: "mock"},
		{Name: "ITC Gardenia", Location: "Residency Road", Price: 12000, Rating: 4.6, Source: "mock"},
		{Name: "Taj West End", Location: "Race Course Road", Price: 14000, Rating: 4.7, Source: "mock"},
		{Name: "The Oberoi Bengaluru", Location: "MG Road", Price: 15000, Rating: 4.7, Source: "mock"},
		{Name: "JW Marriott Hotel Bengaluru", Location: "Vittal Mallya Road", Price: 13000, Rating: 4.6, Source: "mock"},
		{Name: "Sheraton Grand Bangalore", Location: "Brigade Gateway", Price: 11000, Rating: 4.5, Source: "mock"},
		{Name: "The Ritz-Carlton Bangalore", Location: "Residency Road", Price: 17000, Rating: 4.8, Source: "mock"},
		{Name: "Hilton Bangalore Embassy", Location: "Domlur", Price: 10000, Rating: 4.5, Source: "mock"},
		{Name: "Vivanta Bengaluru", Location: "Whitefield", Price: 9000, Rating: 4.4, Source: "mock"},
		{Name: "Radisson Blu Atria", Location: "Palace Road", Price: 8500, Rating: 4.3, Source: "mock"},
	},
}
type ScraperConfig struct {
	Headless bool
	Timeout  time.Duration
}
func ScrapeCity(city string) []models.Hotel {
	log.Printf("🔍 Starting scrape for %s...", city)
	hotels := []models.Hotel{}
	mmtHotels := scrapeMakeMyTrip(city)
	if len(mmtHotels) > 0 {
		log.Printf("Scraped %d hotels from MakeMyTrip for %s", len(mmtHotels), city)
		hotels = append(hotels, mmtHotels...)
	}
	goibiboHotels := scrapeGoibibo(city)
	if len(goibiboHotels) > 0 {
		log.Printf("Scraped %d hotels from Goibibo for %s", len(goibiboHotels), city)
		hotels = append(hotels, goibiboHotels...)
	}
	if len(hotels) == 0 {
		log.Printf("All scrapers failed for %s, using mock data", city)
		hotels = getMockDataWithVariation(city)
	}
	log.Printf("Total %d hotels for %s", len(hotels), city)
	return hotels
}
func scrapeMakeMyTrip(city string) []models.Hotel {
	hotels := []models.Hotel{}
	cityURLs := map[string]string{
		"Delhi":     "https://www.makemytrip.com/hotels/hotel-listing/?checkin=11162025&checkout=11192025&locusId=CTDEL&locusType=city&city=CTDEL&country=IN&searchText=Delhi&roomStayQualifier=2e0e&_uCurrency=INR&reference=hotel&type=city&rsc=1e2e0e",
		"Mumbai":    "https://www.makemytrip.com/hotels/hotel-listing/?checkin=11162025&checkout=11192025&locusId=CTBOM&locusType=city&city=CTBOM&country=IN&searchText=Mumbai&roomStayQualifier=2e0e&_uCurrency=INR&reference=hotel&type=city&rsc=1e2e0e",
		"Bangalore": "https://www.makemytrip.com/hotels/hotel-listing/?checkin=11162025&city=CTBLR&checkout=11192025&roomStayQualifier=2e0e&locusId=CTBLR&country=IN&locusType=city&searchText=Bengaluru&regionNearByExp=3&rsc=1e2eundefinede",
	}
	url, exists := cityURLs[city]
	if !exists {
		return hotels
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.ScrollIntoView(`.listingRowOuter`, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		chromedp.OuterHTML(`html`, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		log.Printf("ChromeDP error for MakeMyTrip: %v", err)
		return hotels
	}
	hotels = parseMakeMyTripHTML(htmlContent, city)
	return hotels
}
func scrapeGoibibo(city string) []models.Hotel {
	hotels := []models.Hotel{}
	cityURLs := map[string]string{
		"Delhi":     "https://www.goibibo.com/hotels/hotel-listing/?checkin=20251118&checkout=20251119&roomString=1-2-0&searchText=Delhi&locusId=CTDEL&locusType=city&cityCode=CTDEL&cc=IN&_uCurrency=INR&vcid=2820046943342890302&sType=city",
		"Mumbai":    "https://www.goibibo.com/hotels/hotel-listing/?checkin=20251118&checkout=20251119&roomString=1-2-0&searchText=Mumbai&locusId=CTBOM&locusType=city&cityCode=CTBOM&cc=IN&_uCurrency=INR&vcid=4213513766539949483&sType=city",
		"Bangalore": "https://www.goibibo.com/hotels/hotel-listing/?checkin=20251118&checkout=20251119&roomString=1-2-0&searchText=Bengaluru&locusId=CTBLR&locusType=city&cityCode=CTBLR&cc=IN&_uCurrency=INR&vcid=6771549831164675055",
	}
	url, exists := cityURLs[city]
	if !exists {
		return hotels
	}
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Nodes(`.HotelCardstyles__Card-sc-gchsmo-0`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Printf("ChromeDP error for Goibibo: %v", err)
		return hotels
	}
	if len(nodes) == 0 {
		log.Printf("No hotel cards found on Goibibo for %s", city)
		return hotels
	}
	for i, node := range nodes {
		if i >= 20 {
			break
		}
		var name, priceText, location, ratingText string
		chromedp.Run(ctx,
			chromedp.Text(`.HotelCardstyles__HotelName-sc-gchsmo-4`, &name, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`.HotelCardstyles__HotelPrice-sc-gchsmo-9`, &priceText, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`.HotelCardstyles__HotelLocation-sc-gchsmo-5`, &location, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`.HotelCardstyles__HotelRating-sc-gchsmo-6`, &ratingText, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		price := parsePrice(priceText)
		rating := parseRating(ratingText)
		
		if name != "" && price > 0 {
			hotel := models.Hotel{
				Name:     cleanString(name),
				City:     city,
				Location: cleanString(location),
				Price:    price,
				Rating:   rating,
				ImageURL: fmt.Sprintf("https://picsum.photos/seed/%s/400/300", name),
				Source:   "goibibo",
			}
			hotels = append(hotels, hotel)
		}
	}
	return hotels
}
func parseMakeMyTripHTML(html, city string) []models.Hotel {
	hotels := []models.Hotel{}
	lines := strings.Split(html, "\n")
	for _, line := range lines {
		if strings.Contains(line, "makeFlex") && strings.Contains(line, "Hotel") {
			hotel := models.Hotel{
				City:     city,
				Source:   "makemytrip",
				Rating:   4.0 + rand.Float64(),
				ImageURL: "https://picsum.photos/400/300",
			}
			hotels = append(hotels, hotel)
		}
	}
	return hotels
}
func parsePrice(priceStr string) float64 {
	if priceStr == "" {
		return 0
	}
	cleaned := strings.NewReplacer("₹", "", "$", "", ",", "", " ", "").Replace(priceStr)
	cleaned = strings.TrimSpace(cleaned)
	var price float64
	fmt.Sscanf(cleaned, "%f", &price)
	if strings.Contains(priceStr, "$") {
		price *= 83
	}
	if price > 0 && price < 1000 {
		price *= 100
	}
	return price
}
func parseRating(ratingStr string) float64 {
	if ratingStr == "" {
		return 4.0 + rand.Float64()
	}
	cleaned := strings.TrimSpace(ratingStr)
	cleaned = strings.Split(cleaned, "/")[0]
	var rating float64
	fmt.Sscanf(cleaned, "%f", &rating)
	if rating == 0 || rating > 5 {
		return 4.0 + rand.Float64()
	}
	return rating
}
func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	return strings.Join(strings.Fields(s), " ")
}
func getMockDataWithVariation(city string) []models.Hotel {
	mockData, exists := mockHotels[city]
	if !exists {
		return []models.Hotel{}
	}
	hotels := make([]models.Hotel, len(mockData))
	for i, hotel := range mockData {
		variation := 1.0 + (rand.Float64()-0.5)*0.3
		
		hotels[i] = models.Hotel{
			Name:     hotel.Name,
			City:     city,
			Location: hotel.Location,
			Price:    hotel.Price * variation,
			Rating:   hotel.Rating,
			ImageURL: fmt.Sprintf("https://picsum.photos/seed/%s/400/300", hotel.Name),
			Source:   "mock",
		}
	}
	return hotels
}

