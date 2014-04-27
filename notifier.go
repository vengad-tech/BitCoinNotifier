package main
import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "os/exec"
import "math"
import "log"
import "flag"
//This method shows mac-os specific notification
//you need terminal-notifier command line installed to use this
//read the README to install this utility
func showNotification(title , message string){
_ ,err := exec.Command("/usr/bin/terminal-notifier","-title",title,"-message",message).Output()
if err != nil {
	log.Fatal("Error in showing notification ,is terminal-notifier installed ",err)
}

}
//Fetches data by reading the given URI
func fetchURLData(addr string) ([]byte ,error) {
	resp, err := http.Get(addr)
	if err != nil {
		log.Println("Unable to connect to block chain")
		return nil,err
	}
	defer resp.Body.Close()
	data,err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		log.Println("unable to read stream from network")
		return nil,err1
	}
	return data , nil
}
//Fetches current bitcoin exchange price by reading blockchain ticket information
func fetchTicker(rate chan float64, err chan error) {
	type Price struct {
		Last float64
	}
	type Ticker struct {
		USD Price
		
	}
	var obj Ticker
	addr := "https://blockchain.info/ticker"
	data,err1 := fetchURLData(addr)
	if err1 != nil {
		log.Println("Unable to fetch url data")
		err<-err1
		return
	}

	err2 := json.Unmarshal(data,&obj)
	if err2 != nil{
		log.Println("Unable to unmarshal json object ",err2)
		err<-err2
		return
	}
	//fmt.Println("Decoded JSON is ",obj)
	err <- nil
	rate<-obj.USD.Last 
}
//Fetches balance for a given bitcoin address
func fetchBalance(address string,balance chan float64, err chan error){
	type Wallet struct {
		Address string
		Final_balance float64
		
	}
	var obj Wallet
	addr := fmt.Sprintf("https://blockchain.info/address/%s?format=json",address)
	data,err1 := fetchURLData(addr)
	if err1 != nil {
		fmt.Println("Unable to fetch url data")
		err <-err1
		return
	}

	err2 := json.Unmarshal(data,&obj)
	if err2 != nil{
		fmt.Println("Unable to unmarshal json object ",err2)
		err <- err2
		return
	}
	//fmt.Println("Decoded JSON is ",obj)
	err <- nil
	balance <- obj.Final_balance / 100000000
	
}

func main(){
	//64.36
	//Getting the necessary command line options
	buying_price := flag.Float64("buyingprice",-1 ,"Required , Amount in USD that you invested in bitcoins")
	bitcoin_address :=flag.String("btaddress","Required","Bitcoin address to check")
	flag.Parse()
	if *bitcoin_address == "Required" {
		log.Fatal("Bitcoin address is a required field use -h option to see options present")
	}
	if *buying_price == -1 {
		log.Fatal("Buying price is a required field use -h option to see options present")
	}
	balance := make(chan float64)
	newrate :=make(chan float64)
	bal_err := make(chan error)
	rate_err := make(chan error)
	//fetcing thebalance and ticker ||ly
	go fetchBalance(*bitcoin_address,balance,bal_err)
	go fetchTicker(newrate,rate_err)
	if <-bal_err != nil || <-rate_err != nil  {
		//If error occured
		showNotification("Unable to fetch Bitcoin details","Connected to internet ??")
		log.Fatal("Got error in fetching data")
		return
	}
	bal , rate :=  <-balance , <-newrate
	current_price := bal * rate
	earning := current_price - *buying_price
	msg := ""
	if earning < 0 {
		msg = fmt.Sprintf("Under loss of $%f",math.Abs(earning))
	} else {
		msg = fmt.Sprintf("Under profit of $%f %s",earning,"Is it good time to Sell !!??")
	}
	//Showing the notification
	showNotification("Bitcoin investment details ",msg)
	
}
