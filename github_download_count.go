package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"
)

const GitHubAPIURI = "https://api.github.com"

type Release struct {
    Assets []Asset `json:"assets"`
}

type Asset struct {
    DownloadCount int `json:"download_count"`
}

func getDownloadCounts(releases []Release) (int, int) {
    downloadCount := releases[0].Assets[0].DownloadCount
    totalDownloadsCount := 0
    for _, release := range releases {
        totalDownloadsCount += release.Assets[0].DownloadCount
    }

    return downloadCount, totalDownloadsCount
}

func getReleasesJSON(uri string) []Release {
    _, err := url.Parse(uri)
    if err != nil {
        panic(err)
    }
    netClient := &http.Client {
        Timeout: time.Second * 10,
    }
    resp, err := netClient.Get(uri)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        panic("Error getting JSON from GitHub API.")
    }
    buf, _ := ioutil.ReadAll(resp.Body)
    var releases []Release
    err = json.Unmarshal(buf, &releases)
    if err != nil {
        panic(err)
    }

    return releases
}

func printUsage() {
    fmt.Printf("Usage: %s organization repository\n", os.Args[0])
}

func main() {
    if len(os.Args) != 3 {
        printUsage()
        return
    }
    org := strings.ToLower(os.Args[1])
    repo := strings.ToLower(os.Args[2])
    uri := fmt.Sprintf(
        "%s/repos/%s/%s/releases",
        GitHubAPIURI,
        org,
        repo)
    json := getReleasesJSON(uri)
    dc, tdc := getDownloadCounts(json)
    fmt.Printf("Current Release: %d\n", dc)
    fmt.Printf("   All Releases: %d\n", tdc)
}
