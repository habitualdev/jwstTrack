package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
"github.com/micmonay/keybd_event"
)

const telescope = `
Where
    is      #  /\
      the    #/  \
    JWST?    /    \
            /     /
           /     /
          /     /
         /     /
        /     /
       /     /
      /     /\
     /     /__|
    /     /\_/\
   /     /   \|
  /     /    ||
  \    /     ||
   \  /     /__\
    \/    .''..''.
        .'   ''   '.
       .'    ''    '.
      .'     ''     '.
     .'      ''      '.
    .'       ''       '.
   .'        ''        '.  
  .'         ''         '.
 .'          ''          '.
<>           <>           <>
`

const auConstant = 149598073
const earthFromMars = "https://theskylive.com//objects/mars/chartdata_dg.json"
const jwstFromEarth = "https://theskylive.com/objects/jwst/chartdata_dg.json"

func SwapWorkspace() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}
	kb.HasSuper(true)
	kb.HasSHIFT(true)
	kb.SetKeys(keybd_event.VK_DOWN)
	err = kb.Launching()
	if err != nil {
		panic(err)
	}
	kb.Press()
	time.Sleep(10 * time.Millisecond)
	kb.Release()
}
func CalculateDistance(){
	dateExtract, _ := regexp.Compile(time.Now().Format("2006/01/02")  + ".*\\d.\\d{5}")
	backupDataExtract, _ := regexp.Compile(time.Now().AddDate(0,0,1).Format("2006/01/02")  + ".*\\d.\\d{5}")
	respJwst, err := http.Get(jwstFromEarth)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(respJwst.Body)
	if err != nil {
		log.Fatalln(err)
	}
	respMars, err := http.Get(earthFromMars)
	if err != nil {
		log.Fatalln(err)
	}.
	bodyMars, err := ioutil.ReadAll(respMars.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var splitData []string
	var splitDataMars []string
	if dateExtract.Match(body) {
		splitData = strings.Split(string(dateExtract.Find(body)),",")
	} else {
		splitData = strings.Split(string(backupDataExtract.Find(body)),",")
	}
	if dateExtract.Match(bodyMars) {
		splitDataMars = strings.Split(string(dateExtract.Find(bodyMars)),",")
	} else {
		splitDataMars = strings.Split(string(backupDataExtract.Find(bodyMars)),",")
	}
	jwstDistanceAU, _ := strconv.ParseFloat(splitData[2],8)
	marsDistanceAU, _ := strconv.ParseFloat(splitDataMars[2],8)
	jwstDistanceKM := fmt.Sprintf("%.2f", jwstDistanceAU * auConstant)
	marsDistanceKM :=fmt.Sprintf("%.2f", (marsDistanceAU - jwstDistanceAU) * auConstant)
	entry1 := fmt.Sprintf("Distance to JWST: %s KM\n", jwstDistanceKM)
	entry2 := fmt.Sprintf("Distance from JWST to Mars: %s KM\n", marsDistanceKM)
	ioutil.WriteFile("distances.txt", []byte(telescope + "\n" + entry1 + entry2), 0644)
	fmt.Println(entry1)
	fmt.Println(entry2)

}

func startGedit(){
	pwd, _ := os.Getwd()
	cmd := exec.Command("/usr/bin/gedit",  pwd +"/distances.txt")
	cmd.Start()
}

func main(){
	fmt.Println(telescope)
	CalculateDistance()
	startGedit()
	time.Sleep(500 * time.Millisecond)
	SwapWorkspace()
}