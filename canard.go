package main

import (
  "log"
  "strings"

  "github.com/gdamore/tcell/v2"
  "github.com/rivo/tview"

  "github.com/JohannesKaufmann/html-to-markdown"
  "github.com/grokify/html-strip-tags-go"
  "html"
  // "github.com/charmbracelet/glamour"
)

var VERSION string

type CanardItem struct {
  *Item
  Markdown                   string
  PlainText                  string
  FeedID                     int
  FeedTitle                  string
}

type Canard struct {
  App                        *tview.Application
  FeedSwitcher               *tview.DropDown
  ItemsList                  *tview.List
  Grid                       *tview.Grid

  ApiURL                     string
  ApiKey                     string

  Feeds                      []Feed

  Items                      []CanardItem
  ItemsMap                   map[int]int

  CurrentFeedID              int
}

type Item struct {
  ID                         int     `json:"id"`
  FeedID                     int     `json:"feed_id"`
  Title                      string  `json:"title"`
  URL                        string  `json:"url"`
  HTML                       string  `json:"html"`
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
  canard := Canard{
    ItemsMap: make(map[int]int),
    CurrentFeedID: -1,
  }
  canard.ApiURL = LookupStrEnv(
    "CANARD_API_URL",
    "http://127.0.0.1:8000/fever/",
  )
  canard.ApiKey = LookupStrEnv(
    "CANARD_API_KEY",
    "9a0f36d70a22b40baa26f3df113cd9eb",
  )

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
    SetLabel("Feed: ")

  canard.ItemsList = tview.NewList().
    SetWrapAround(true).
    SetHighlightFullLine(true).
    SetSecondaryTextColor(tcell.ColorGrey)

  canard.Grid = tview.NewGrid().
    SetRows(1, 0).
    SetColumns(0).
    SetBorders(true).
    AddItem(canard.FeedSwitcher, 0, 0, 1, 1, 0, 0, false).
    AddItem(canard.ItemsList, 1, 0, 1, 1, 0, 0, true)

  canard.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    switch event.Key() {
    case tcell.KeyCtrlT:
      if canard.FeedSwitcher.HasFocus() == false {
        canard.App.SetFocus(canard.FeedSwitcher)
        canard.App.QueueEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
      } else {
        canard.App.SetFocus(canard.ItemsList)
      }
      return nil
    case tcell.KeyCtrlQ:
      canard.App.Stop()
    }

    return event
  })

  canard.Refresh()
  canard.RefreshUI()

  // canard.App.SetFocus(canard.FeedSwitcher)
  if err := canard.App.SetRoot(canard.Grid, true).Run(); err != nil {
    panic(err)
  }
}

func (canard *Canard) Refresh() () {
  apiResponse, err := call(canard.ApiKey, canard.ApiURL + "?api&feeds&items")
  if err != nil {
    panic(err)
  }

  var feedMap = make(map[int]Feed)

  for _, feed := range apiResponse.Feeds {
    feedMap[feed.ID] = feed
  }

  canard.Feeds = apiResponse.Feeds

  converter := md.NewConverter("", true, nil)
  itemsLen := len(apiResponse.Items)
  for i := (itemsLen - 1); i >= 0; i-- {
    item := apiResponse.Items[i]
    _, exists := canard.ItemsMap[item.ID]
    if exists == true {
      continue
    }

    markdown, err := converter.ConvertString(item.HTML)
    if err != nil {
      log.Fatal(err)
    }

    canardItem := CanardItem{
      &item,
      markdown,
      strings.TrimSpace(html.UnescapeString(strip.StripTags(item.HTML))),
      feedMap[item.FeedID].ID,
      feedMap[item.FeedID].Title,
    }

    canard.ItemsMap[item.ID] = len(canard.Items)
    canard.Items = append(canard.Items, canardItem)
  }
}

func (canard *Canard) SwitchByID(feedID int) (bool) {
  canard.CurrentFeedID = feedID
  return true
}

func (canard *Canard) Switch(feedTitle string) (bool) {
  if feedTitle == "All" {
    canard.CurrentFeedID = -1
    return true
  }

  for _, feed := range canard.Feeds {
    if feedTitle == feed.Title {
      canard.CurrentFeedID = feed.ID
      return true
    }
  }

  return false
}

func (canard *Canard) RefreshUI() (bool) {
  canard.FeedSwitcher.SetOptions(
    []string{"All"},
    func(text string, index int) {
      canard.Switch(text)
      canard.RefreshUI()
      canard.App.SetFocus(canard.ItemsList)
    },
  )

  for i := 0; i < len(canard.Feeds); i++ {
    feed := canard.Feeds[i]
    canard.FeedSwitcher.AddOption(feed.Title, nil)
  }

  canard.ItemsList.Clear()
  for _, item := range canard.Items {
    if item.FeedID == canard.CurrentFeedID ||
       canard.CurrentFeedID == -1 {
      canard.ItemsList.AddItem(item.Title, item.FeedTitle, 0, nil)
    }
  }

  canard.ItemsList.SetCurrentItem(-1)
  return true
}
