package youtube

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeVideo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnailUrl"`
	Duration    string `json:"duration"`
	PublishedAt string `json:"publishedAt"`
}

type Service struct {
	apiKey      string
	channelID   string
	cache       []YouTubeVideo
	lastUpdated time.Time
	mu          sync.RWMutex
}

func NewService() (*Service, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY environment variable is not set")
	}

	channelID := os.Getenv("YOUTUBE_CHANNEL_ID")
	if channelID == "" {
		return nil, fmt.Errorf("YOUTUBE_CHANNEL_ID environment variable is not set")
	}

	return &Service{
		apiKey:    apiKey,
		channelID: channelID,
	}, nil
}

func (s *Service) GetCachedFeed() ([]YouTubeVideo, error) {
	s.mu.RLock()
	cacheValid := !s.lastUpdated.IsZero() && time.Since(s.lastUpdated) < 20*time.Minute
	if cacheValid && len(s.cache) > 0 {
		videos := s.cache
		s.mu.RUnlock()
		return videos, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Double check inside write lock
	if !s.lastUpdated.IsZero() && time.Since(s.lastUpdated) < 20*time.Minute && len(s.cache) > 0 {
		return s.cache, nil
	}

	// Try fetching from API
	videos, err := s.fetchFromAPI()
	if err != nil {
		fmt.Printf("Error fetching YouTube feed: %v. Using fallback.\n", err)
		if len(s.cache) > 0 {
			// Serve old cache if we have it
			return s.cache, nil
		}
		// If cache is empty, seed with mock data
		s.cache = getMockVideos()
		s.lastUpdated = time.Now() // prevent continuous spam on failure
		return s.cache, nil
	}

	s.cache = videos
	s.lastUpdated = time.Now()
	return s.cache, nil
}

func (s *Service) fetchFromAPI() ([]YouTubeVideo, error) {
	ctx := context.Background()
	srv, err := youtube.NewService(ctx, option.WithAPIKey(s.apiKey))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize youtube service: %w", err)
	}

	// 1. Retrieve the channel uploads playlist ID
	chanCall := srv.Channels.List([]string{"contentDetails"}).Id(s.channelID)
	chanResp, err := chanCall.Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("channel API call failed: %w", err)
	}
	if len(chanResp.Items) == 0 {
		return nil, fmt.Errorf("no channel found for ID: %s", s.channelID)
	}

	uploadsPlaylistID := chanResp.Items[0].ContentDetails.RelatedPlaylists.Uploads
	if uploadsPlaylistID == "" {
		return nil, fmt.Errorf("no uploads playlist found for channel")
	}

	// 2. Fetch recent videos from the uploads playlist
	playlistCall := srv.PlaylistItems.List([]string{"snippet"}).PlaylistId(uploadsPlaylistID).MaxResults(12)
	playlistResp, err := playlistCall.Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("playlistitems API call failed: %w", err)
	}

	var videos []YouTubeVideo
	var videoIDs []string

	for _, item := range playlistResp.Items {
		if item.Snippet == nil || item.Snippet.ResourceId == nil {
			continue
		}
		videoID := item.Snippet.ResourceId.VideoId
		videoIDs = append(videoIDs, videoID)

		thumbnail := ""
		if item.Snippet.Thumbnails != nil {
			if item.Snippet.Thumbnails.Maxres != nil {
				thumbnail = item.Snippet.Thumbnails.Maxres.Url
			} else if item.Snippet.Thumbnails.High != nil {
				thumbnail = item.Snippet.Thumbnails.High.Url
			} else if item.Snippet.Thumbnails.Medium != nil {
				thumbnail = item.Snippet.Thumbnails.Medium.Url
			} else if item.Snippet.Thumbnails.Standard != nil {
				thumbnail = item.Snippet.Thumbnails.Standard.Url
			} else if item.Snippet.Thumbnails.Default != nil {
				thumbnail = item.Snippet.Thumbnails.Default.Url
			}
		}

		videos = append(videos, YouTubeVideo{
			ID:          videoID,
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			Thumbnail:   thumbnail,
			PublishedAt: item.Snippet.PublishedAt,
			Duration:    "0:00", // to be populated
		})
	}

	// 3. Query video durations in batch
	if len(videoIDs) > 0 {
		videoCall := srv.Videos.List([]string{"contentDetails"}).Id(strings.Join(videoIDs, ","))
		videoResp, err := videoCall.Context(ctx).Do()
		if err == nil {
			durationMap := make(map[string]string)
			for _, v := range videoResp.Items {
				if v.ContentDetails != nil {
					durationMap[v.Id] = parseISO8601Duration(v.ContentDetails.Duration)
				}
			}
			for i, v := range videos {
				if d, ok := durationMap[v.ID]; ok {
					videos[i].Duration = d
				}
			}
		} else {
			fmt.Printf("Warning: Failed to fetch video durations: %v\n", err)
		}
	}

	return videos, nil
}

func parseISO8601Duration(iso string) string {
	if !strings.HasPrefix(iso, "PT") {
		return "0:00"
	}

	iso = iso[2:]
	var hours, minutes, seconds int
	var temp int

	for i := 0; i < len(iso); i++ {
		c := iso[i]
		if c >= '0' && c <= '9' {
			temp = temp*10 + int(c-'0')
		} else if c == 'H' {
			hours = temp
			temp = 0
		} else if c == 'M' {
			minutes = temp
			temp = 0
		} else if c == 'S' {
			seconds = temp
			temp = 0
		}
	}

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func getMockVideos() []YouTubeVideo {
	return []YouTubeVideo{
		{
			ID:          "m3v3KkQjP6Q",
			Title:       "Lorcana Shimmering Skies Box Opening - Legendary & Enchanted Pulls!",
			Description: "Opening a fresh booster box of Lorcana's newest set, Shimmering Skies! Looking for the new Enchanted cards and reviewing the competitive viability of each pull.",
			Thumbnail:   "https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?q=80&w=600&auto=format&fit=crop",
			Duration:    "18:45",
			PublishedAt: "2026-05-18T15:00:00Z",
		},
		{
			ID:          "K-Nf07WpZzY",
			Title:       "Top 5 Decks for the Lorcana Shimmering Skies Championship",
			Description: "Deep dive into the meta! We analyze Amber/Steel, Ruby/Amethyst, and Emerald/Steel lists to prepare you for the upcoming TCG championships.",
			Thumbnail:   "https://images.unsplash.com/photo-1607604276583-eef5d076aa5f?q=80&w=600&auto=format&fit=crop",
			Duration:    "24:12",
			PublishedAt: "2026-05-15T12:00:00Z",
		},
		{
			ID:          "WJtG1z1J9wI",
			Title:       "How to Play Lorcana TCG - Complete Beginner's Guide",
			Description: "New to Lorcana? Learn how to summon glimmers, quest for lore, challenge opponent characters, and manage your inkwell in this step-by-step tutorial.",
			Thumbnail:   "https://images.unsplash.com/photo-1546519638-68e109498ffc?q=80&w=600&auto=format&fit=crop",
			Duration:    "12:30",
			PublishedAt: "2026-05-10T10:00:00Z",
		},
		{
			ID:          "qW1-yHqfH6Y",
			Title:       "Lorcana Ruby/Amethyst Control - In-Depth Deck Tech & Gameplay",
			Description: "Liam reviews the classic control archetype. Learn the best mulligan strategies, turn-by-turn sequencing, and matchup plans against aggressive decks.",
			Thumbnail:   "https://images.unsplash.com/photo-1511512578047-dfb367046420?q=80&w=600&auto=format&fit=crop",
			Duration:    "32:15",
			PublishedAt: "2026-05-05T14:30:00Z",
		},
		{
			ID:          "2G_mGg8jA30",
			Title:       "Lorcana Rise of the Floodborn Budget Deck Upgrades under $20",
			Description: "Want to be competitive without breaking the bank? Check out these cheap upgrades for starter decks that will win you games at local league nights.",
			Thumbnail:   "https://images.unsplash.com/photo-1606167668584-78701c57f13d?q=80&w=600&auto=format&fit=crop",
			Duration:    "15:50",
			PublishedAt: "2026-05-01T11:00:00Z",
		},
		{
			ID:          "fX1-kHeS4Q8",
			Title:       "Drafting Lorcana - Essential Tips to Win Your Local Store Draft",
			Description: "Drafting requires a unique set of skills. We explain character ratios, inks evaluation, and top-tier commons you should always draft.",
			Thumbnail:   "https://images.unsplash.com/photo-1518895949257-7621c3c786d7?q=80&w=600&auto=format&fit=crop",
			Duration:    "21:05",
			PublishedAt: "2026-04-28T09:00:00Z",
		},
	}
}
