package search

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	_ "github.com/mark3labs/mcp-go/mcp"
	"github.com/zmb3/spotify/v2"
	"log"
	"spotify-mcp/internal/client"
	"spotify-mcp/internal/server/tools"
)

const playlsitOrAlbumNameParameter = "Playlist Name"

func PlayListSearchTools() []tools.ToolEntry {
	return []tools.ToolEntry{
		simplePlaylistSearch(),
	}
}

func simplePlaylistSearch() tools.ToolEntry {
	toolDefinition := mcp.NewTool(
		"simple_playlist_and_album_search",
		mcp.WithDescription("Search for a playlist or album by name"),
		mcp.WithString(playlsitOrAlbumNameParameter,
			mcp.Required(),
			mcp.Description("Name of the playlist or album to search for. Extra information: "+SearchQueryInformation),
		),
	)

	toolBehaviour := playlistSearchBehaviour

	return tools.ToolEntry{
		ToolDefinition: toolDefinition,
		ToolBehaviour:  toolBehaviour,
	}
}

func playlistSearchBehaviour(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	playlistName, err := tools.GetParamFromRequest(request, playlsitOrAlbumNameParameter)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist name: %w", err)
	}

	results, err := client.SpotifyClient.Search(ctx, playlistName, spotify.SearchTypePlaylist|spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	jsonPlaylists := ""
	if results.Playlists != nil {
		bytePlaylists, err := json.Marshal(results.Playlists)
		if err != nil {
			log.Fatal(err)
		}

		jsonPlaylists = string(bytePlaylists)
	}

	mcpResult := mcp.NewToolResultText(jsonPlaylists)

	return mcpResult, nil
}
