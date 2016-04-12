package main

import "fmt"
import "os"
import "strconv"
import "encoding/json"
import "github.com/codegangsta/cli"
import "io/ioutil"
import bp "github.com/bitpay/bitpay-go/client"
import ku "github.com/bitpay/bitpay-go/key_utils"

type config struct {
	Pem string
	Sin string
	Tok bp.Token
}

func main() {
	app := cli.NewApp()

	app.Name = "go-bitpay"
	app.Usage = "go bitpay pair|create|get"
	app.Author = "Ekin Koc"
	app.Email = "ekin@eknkc.com"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "pair",
			Usage: "pair <pairCode> - Pair with BitPay API and write bitpay.json to current dir",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 1 {
					fmt.Println("pair code not specified.")
					return
				}

				pem := ku.GeneratePem()
				sin := ku.GenerateSinFromPem(pem)

				bc := bp.Client{
					ApiUri:   "https://bitpay.com",
					Pem:      pem,
					ClientId: sin,
				}

				token, err := bc.PairWithCode(c.Args()[0])

				if err != nil {
					fmt.Println(err)
				} else {
					conf := config{
						Pem: pem,
						Sin: sin,
						Tok: token,
					}

					confByte, err := json.Marshal(conf)

					if err != nil {
						fmt.Println(err)
					} else {
						ioutil.WriteFile("bitpay.json", confByte, 0644)
					}
				}
			},
		},
		cli.Command{
			Name:  "create",
			Usage: "create <price> - Create a new invoice for <price> usd",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 1 {
					fmt.Println("price not specified.")
					return
				}

				price, err := strconv.ParseFloat(c.Args()[0], 64)

				if err != nil {
					fmt.Println(err)
				} else {
					bc, err := loadClient()

					if err != nil {
						fmt.Println(err)
						return
					}

					inv, err := bc.CreateInvoice(price, "USD")

					if err != nil {
						fmt.Println(err)
						return
					}

					js, err := json.MarshalIndent(inv, "", "  ")

					if err != nil {
						fmt.Println(err)
						return
					}

					fmt.Println(string(js))
				}
			},
		},
		cli.Command{
			Name:  "get",
			Usage: "get <invoiceId> - Get the status of an invoice",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 1 {
					fmt.Println("id not specified.")
					return
				}

				bc, err := loadClient()

				if err != nil {
					fmt.Println(err)
					return
				}

				inv, err := bc.GetInvoice(c.Args()[0])

				if err != nil {
					fmt.Println(err)
					return
				}

				js, err := json.MarshalIndent(inv, "", "  ")

				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Println(string(js))
			},
		},
	}

	app.Run(os.Args)
}

func loadClient() (*bp.Client, error) {
	confBytes, err := ioutil.ReadFile("bitpay.json")

	if err != nil {
		return nil, err
	}

	var conf config
	err = json.Unmarshal(confBytes, &conf)

	if err != nil {
		return nil, err
	}

	bc := bp.Client{
		ApiUri:   "https://bitpay.com",
		Pem:      conf.Pem,
		ClientId: conf.Sin,
		Token:    conf.Tok,
	}

	return &bc, nil
}
