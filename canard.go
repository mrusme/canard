package main

import (
  "log"
  "time"
  "bytes"
  "strings"
  "io"
  "net/http"
  "mime/multipart"
  "encoding/json"

  "github.com/gdamore/tcell/v2"
  "github.com/rivo/tview"

  "github.com/JohannesKaufmann/html-to-markdown"
  "github.com/grokify/html-strip-tags-go"
  "html"
  // "github.com/charmbracelet/glamour"
)

var VERSION string

type Canard struct {
  App                        *tview.Application
  FeedSwitcher               *tview.DropDown
  ItemsList                  *tview.List
  Grid                       *tview.Grid

  ApiURL                     string
  ApiKey                     string
}

type Item struct {
  ID                         int     `json:"id"`
  FeedID                     int     `json:"feed_id"`
  Title                      string  `json:"title"`
  URL                        string  `json:"url"`
  HTML                       string  `json:"html"`
  Markdown                   string
  PlainText                  string
  CreatedOnTime              int     `json:"created_on_time"`
  IsRead                     int     `json:"is_read"`
  IsSaved                    int     `json:"is_saved"`
}

type Feed struct {
  ID                         int     `json:"id"`
  Title                      string  `json:"title"`
  SiteURL                    string  `json:"site_url"`
  URL                        string  `json:"url"`
  LastUpdatedOnTime          int     `json:"last_updated_on_time"`
}

type ApiResponse struct {
  ApiVersion                 string  `json:"api_version"`
  Auth                       int     `json:"auth"`
  Feeds                      []Feed  `json:"feeds"`
  Items                      []Item  `json:"items"`

}

func main() {
  canard := Canard{}
  canard.ApiURL = LookupStrEnv("CANARD_API_URL", "http://127.0.0.1:8000/fever/")
  canard.ApiKey = LookupStrEnv("CANARD_API_KEY", "9a0f36d70a22b40baa26f3df113cd9eb")

  // tview.Styles = tview.Theme{
  //   PrimitiveBackgroundColor:    tcell.ColorDefault,
  //   ContrastBackgroundColor:     tcell.ColorTeal,
  //   MoreContrastBackgroundColor: tcell.ColorTeal,
  //   BorderColor:                 tcell.ColorWhite,
  //   TitleColor:                  tcell.ColorWhite,
  //   GraphicsColor:               tcell.ColorWhite,
  //   PrimaryTextColor:            tcell.ColorDefault,
  //   SecondaryTextColor:          tcell.ColorBlue,
  //   TertiaryTextColor:           tcell.ColorGreen,
  //   InverseTextColor:            tcell.ColorBlack,
  //   ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
  // }

  canard.App = tview.NewApplication()

  canard.FeedSwitcher = tview.NewDropDown().
    SetLabel("Feed: ").
    SetOptions([]string{"All"}, nil)

  canard.ItemsList = tview.NewList().
    SetWrapAround(true).
    SetHighlightFullLine(true).
    SetSecondaryTextColor(tcell.ColorGrey)

  canard.Grid = tview.NewGrid().
    SetRows(1, 0).
    SetColumns(0).
    SetBorders(true).
    AddItem(canard.FeedSwitcher, 0, 0, 1, 1, 0, 0, true).
    AddItem(canard.ItemsList, 1, 0, 1, 1, 0, 0, false)

  apiResponse, err := call(canard.ApiKey, canard.ApiURL + "?api&feeds&items")
  if err != nil {
    panic(err)
  }

  for _, feed := range apiResponse.Feeds {
    canard.FeedSwitcher.AddOption(feed.Title, nil)
  }

  converter := md.NewConverter("", true, nil)
  for _, item := range apiResponse.Items {
    item.Markdown, err = converter.ConvertString(item.HTML)
    if err != nil {
      log.Fatal(err)
    }
    item.PlainText = strings.TrimSpace(html.UnescapeString(strip.StripTags(item.HTML)))
    canard.ItemsList.AddItem(item.Title, item.PlainText, 0, nil)
  }
  canard.ItemsList.SetCurrentItem(0)

  // canard.App.SetFocus(canard.FeedSwitcher)
  if err := canard.App.SetRoot(canard.Grid, true).Run(); err != nil {
    panic(err)
  }
}

func call(apiKey string, urlPath string) (ApiResponse, error) {
  client := &http.Client{
    Timeout: time.Second * 10,
  }
  body := &bytes.Buffer{}
  writer := multipart.NewWriter(body)
  fw, err := writer.CreateFormField("api_key")
  if err != nil {
  }
  _, err = io.Copy(fw, strings.NewReader(apiKey))
  if err != nil {
      return ApiResponse{}, err
  }
  writer.Close()
  req, err := http.NewRequest("POST", urlPath, bytes.NewReader(body.Bytes()))
  if err != nil {
      return ApiResponse{}, err
  }
  req.Header.Set("Content-Type", writer.FormDataContentType())
  resp, _ := client.Do(req)
  if resp.StatusCode != http.StatusOK {
    log.Printf("Request failed with response code: %d", resp.StatusCode)
  }

  defer resp.Body.Close()

  var response ApiResponse
  if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
    log.Fatal(err)
  }

  return response, nil
}
