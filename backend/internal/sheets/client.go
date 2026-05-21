package sheets

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type BlogPost struct {
	Date     string `json:"date"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type Client struct {
	srv         *sheets.Service
	sheetID     string
	cache       []BlogPost
	lastUpdated time.Time
	mu          sync.RWMutex
}

func NewClient() (*Client, error) {
	ctx := context.Background()
	credentialsJSON := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")
	if credentialsJSON == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS_JSON environment variable is not set")
	}

	sheetID := os.Getenv("GOOGLE_SHEET_ID")
	if sheetID == "" {
		return nil, fmt.Errorf("GOOGLE_SHEET_ID environment variable is not set")
	}

	// Clean single quotes if they wrap the JSON env var
	credentialsJSON = strings.TrimSpace(credentialsJSON)
	if len(credentialsJSON) > 1 && credentialsJSON[0] == '\'' && credentialsJSON[len(credentialsJSON)-1] == '\'' {
		credentialsJSON = credentialsJSON[1 : len(credentialsJSON)-1]
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(credentialsJSON)))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %w", err)
	}

	return &Client{
		srv:     srv,
		sheetID: sheetID,
	}, nil
}

func (c *Client) AppendSubmission(sheetName string, rowData []interface{}) error {
	ctx := context.Background()
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{rowData},
	}

	// A1 notation range - e.g., "News Tips!A1"
	writeRange := fmt.Sprintf("%s!A1", sheetName)

	_, err := c.srv.Spreadsheets.Values.Append(c.sheetID, writeRange, valueRange).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("unable to append to sheet %s: %w", sheetName, err)
	}

	return nil
}

func (c *Client) GetCachedBlogPosts() ([]BlogPost, error) {
	c.mu.RLock()
	cacheValid := !c.lastUpdated.IsZero() && time.Since(c.lastUpdated) < 20*time.Minute
	if cacheValid && len(c.cache) > 0 {
		posts := c.cache
		c.mu.RUnlock()
		return posts, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Recheck condition inside lock
	if !c.lastUpdated.IsZero() && time.Since(c.lastUpdated) < 20*time.Minute && len(c.cache) > 0 {
		return c.cache, nil
	}

	posts, err := c.FetchBlogPosts()
	if err != nil {
		fmt.Printf("Warning: Failed to fetch blog posts from Google Sheets: %v. Using mock fallback.\n", err)
		if len(c.cache) > 0 {
			return c.cache, nil
		}
		c.cache = GetMockBlogPosts()
		c.lastUpdated = time.Now()
		return c.cache, nil
	}

	c.cache = posts
	c.lastUpdated = time.Now()
	return c.cache, nil
}

func (c *Client) FetchBlogPosts() ([]BlogPost, error) {
	ctx := context.Background()
	// Fetch "Blog Posts" tab. Range A2:D100 to skip headers
	resp, err := c.srv.Spreadsheets.Values.Get(c.sheetID, "Blog Posts!A2:D100").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed reading Blog Posts sheet: %w", err)
	}

	var posts []BlogPost
	for _, row := range resp.Values {
		if len(row) < 3 {
			continue
		}
		post := BlogPost{
			Date:     getString(row, 0),
			Category: getString(row, 1),
			Title:    getString(row, 2),
		}
		if len(row) >= 4 {
			post.Content = getString(row, 3)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getString(row []interface{}, idx int) string {
	if idx < len(row) && row[idx] != nil {
		if s, ok := row[idx].(string); ok {
			return s
		}
		return fmt.Sprintf("%v", row[idx])
	}
	return ""
}

func GetMockBlogPosts() []BlogPost {
	return []BlogPost{
		{
			Date:     "2026-05-20",
			Category: "Meta Strategy",
			Title:    "Lorcana Shimmering Skies: Tier List & Launch Review",
			Content:  "The skies are shimmering and the cards are shaking up the competitive meta! Here is Liam's launch week breakdown of the top decks to watch.\n\n### 1. Amber/Steel Song (Tier 1)\nThis powerhouse archetype continues to dominate. With the new song additions, steel board sweep control is smoother than ever.\n\n### 2. Ruby/Amethyst Bounce Control (Tier 1.5)\nBounce strategies utilizing Arthur and Madame Mim remain highly resilient. Keep an eye on how they match up against aggressive emerald ink speeds.\n\n### 3. Sapphire/Steel Ramp (Tier 2)\nIf you love drawing cards and filling your inkwell, this is your glimmer stack. Ramp quickly into your giant characters and override control decks before they settle in.",
		},
		{
			Date:     "2026-05-15",
			Category: "Beginner Guide",
			Title:    "Managing Your Inkwell: Core Rules of Lore Summoning",
			Content:  "A common mistake for new Illumiteers is mismanaging the inkwell curve. Master these three principles to secure your victories:\n\n*   **Evaluate Inkability**: Make sure you don't run more than 12-14 'uninkable' cards in your deck structure. Too many will clog your hand early.\n*   **Read the Board**: Do not feel pressured to place a card into your inkwell every single turn. If you have your curve set, hold valuable characters for challenges.\n*   **Ink Order Matters**: Always place your ink card before committing to play items. This keeps your available options clear to read.",
		},
		{
			Date:     "2026-05-08",
			Category: "Card Spotlight",
			Title:    "Unlocking Madam Mim - Elephant Synergy",
			Content:  "Madam Mim - Elephant is quickly becoming a staple in amethyst control stacks. Its ability to return a character back to your hand allows you to re-trigger valuable 'when played' abilities (like Merlin - Goat or Merlin - Rabbit) while dropping a massive challenge body onto the field for very low ink costs.",
		},
	}
}
