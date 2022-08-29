# Canard

Canard. A command line TUI client for the 
[Journalist](https://マリウス.com/journalist-an-rss-aggregator/) RSS aggregator.

![[Canard](https://player.vimeo.com/video/535676709)](canard.jpg)

***Note: This version of Canard is not yet compatible with Journalist >= v1.0.0!***

## Installation

Download a binary from the [releases][releases] page.

Or build it yourself (requires Go 1.16+):

```bash
make
```

[releases]: https://github.com/mrusme/canard/releases


## User Manual


### Configuration

Export the following environment variables first:

```sh
export CANARD_API_URL="http://YOUR-JOURNALIST-SERVER:8000/fever/"
export CANARD_API_KEY="YOUR-JOURNALIST-API-KEY"
export GLAMOUR_STYLE="dark"
```

`CANARD_API_URL` and `CANARD_API_KEY` are 
[Journalist](https://github.com/mrusme/journalist)-related configuration
parameters. `CANARD_API_KEY` is the Journalist `username:password` string as MD5
hash. If you're using Journalist in it's default configuration (`nobody:nobody`)
then you can simply set `CANARD_API_KEY` to `9a0f36d70a22b40baa26f3df113cd9eb`.
In case you have different users/passwords, run 
`echo -n "username:password" | md5sum` to retrieve the hash.

`GLAMOUR_STYLE` defines how the articles are being rendered, see 
[`glamour`](https://github.com/charmbracelet/glamour) for more info.

Please make sure you're running the latest version of `journalist`!


### Cheatsheet


#### Shortcuts

This is a list of supported keyboard shortcuts.

`ArrowUp` / `k`, `ArrowDn` / `j` \
Scroll list/reader in either direction by one line

`PgUp` / `b`, `PgDn` / `f` \
Scroll list/reader in either direction by one page

`u`, `d` \
Scroll list/reader in either direction by half a page

`g`, `G` \
Scroll list/reader to the very top/bottom

`q` \
While in reader, close reader; While in list, quit Canard

`Ctrl` + `R` \
Refresh feeds

`Ctrl` + `T` \
Open feed switcher

`Ctrl` + `Q` \
Quit Canard
