package core

import (
	"net/http"
	"strings"
)

var NextWeeksMenu string = `
# Bee Boop, ich habe essen bestellt

## Bestellung

### 12.12.23
* Bestellt  
  * Käsespätzle: *Nudelgericht mit Käse*
* Alternativen:
  * Wurst: *Fleischerzeugnis aus Schwein*
  * Salatteller: *Mit Gemüse*

## Info

* Kontostand: 23,23€
`

func SendNotification() {
	http.Post("https://ntfy.sh/thinkjd_munch_o_matic", "text/markdown",
		strings.NewReader(NextWeeksMenu))

}
