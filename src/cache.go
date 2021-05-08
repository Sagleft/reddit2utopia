package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

const (
	cacheFilename      = "cache.json"
	cachedElementValue = 1
)

// CacheHandler ..
type CacheHandler struct {
	Data      CachedData
	CachePath string
}

type cacheElement int

// CachedData ..
type CachedData struct {
	Posts map[string]map[string]cacheElement `json:"posts"` //chatID (key): postID (key): 1
}

// NewCacheHandler - creates cache handler
func NewCacheHandler(cachePath string) (*CacheHandler, error) {
	if cachePath == "" {
		return nil, errors.New("cache path is not set")
	}

	// check cache folder exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		mkdirErr := os.Mkdir(cachePath, 0755)
		if mkdirErr != nil {
			return nil, mkdirErr
		}
	}

	// check cache file exists
	newCachedData := CachedData{
		Posts: make(map[string]map[string]cacheElement),
	}
	cacheJSONPath := cachePath + "/" + cacheFilename
	if _, err := os.Stat(cacheJSONPath); os.IsNotExist(err) {
		// cache file does not exists
	} else {
		dataJSONBytes, err := ioutil.ReadFile(cacheJSONPath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dataJSONBytes, &newCachedData)
		if err != nil {
			return nil, err
		}
	}
	return &CacheHandler{
		Data:      newCachedData,
		CachePath: cachePath,
	}, nil
}

// IsPostUsed - check post already used
func (cache *CacheHandler) IsPostUsed(chatID, postID string) bool {
	chatPostIDs, chatIDused := cache.Data.Posts[chatID]
	if !chatIDused {
		return false
	}
	_, postUsed := chatPostIDs[postID]
	return postUsed
}

// MarkPostUsed in cache
func (cache *CacheHandler) MarkPostUsed(chatID, postID string) error {
	_, chatIDused := cache.Data.Posts[chatID]
	if !chatIDused {
		cache.Data.Posts[chatID] = map[string]cacheElement{
			postID: cachedElementValue,
		}
	} else {
		cache.Data.Posts[chatID][postID] = cachedElementValue
	}

	cacheJSONPath := cache.CachePath + "/" + cacheFilename
	jsonBytes, err := json.Marshal(cache.Data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cacheJSONPath, jsonBytes, 0777)
	return err
}
