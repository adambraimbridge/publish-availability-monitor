package content

import (
	"net/url"
	"strings"
)

type blogAPIEndpointMetadata struct {
	host string
}

func isValidBrand(requestURI string) bool {
	parsedURL, err := url.Parse(requestURI)
	if err != nil || parsedURL.Host == "" {
		warn.Printf("Invalid request URI  [%s].", requestURI)
		return false
	}
	requestHostAndPath := parsedURL.Host + parsedURL.Path
	for _, blogAPIEndpointMetadata := range blogAPIEndpointMetadatas {
		if strings.Contains(requestHostAndPath, blogAPIEndpointMetadata.host) {
			return true
		}
	}
	return false
}

var blogAPIEndpointMetadatas = []blogAPIEndpointMetadata{
	blogAPIEndpointMetadata{
		"blogs.ft.com/the-world",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/brusselsblog",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/businessblog",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/tech-blog",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/westminster",
	},
	blogAPIEndpointMetadata{
		"ftalphaville.ft.com",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/mba-blog",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/beyond-brics",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/gavyndavies",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/material-world",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/ftdata",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/nick-butler",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/photo-diary",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/off-message",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/david-allen-green",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/andrew-smithers",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/lex-live",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/andrew-mcafee",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/the-exchange",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/larry-summers",
	},
	blogAPIEndpointMetadata{
		"www.ft.com/fastft",
	},
	blogAPIEndpointMetadata{
		"blogs.ft.com/fastft",
	},
}
