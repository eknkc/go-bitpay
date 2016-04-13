package main

import "fmt"
import "os"
import "strconv"
import "encoding/json"
import "github.com/codegangsta/cli"
import "io/ioutil"
import "net/http"
import "errors"
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
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ipn",
					Usage: "Notificateion URL",
				},
				cli.StringFlag{
					Name:  "redirect",
					Usage: "Redirect URL",
				},
			},
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
						panic(err)
					}

					payload := make(map[string]string)
					payload["price"] = strconv.FormatFloat(price, 'f', 2, 64)
					payload["currency"] = "USD"
					payload["token"] = bc.Token.Token
					payload["id"] = bc.ClientId

					if len(c.String("ipn")) > 0 {
						payload["notificationURL"] = c.String("ipn")
						payload["fullNotifications"] = "true"
					}

					if len(c.String("redirect")) > 0 {
						payload["redirectURL"] = c.String("redirect")
					}

					if resp, err := post(bc, "invoices", payload); err != nil {
						panic(err)
					} else {
						js, _ := json.MarshalIndent(resp, "", "  ")
						fmt.Println(string(js))
					}
				}
			},
		},
	}

	app.Run(os.Args)
}

func post(bc *bp.Client, path string, payload map[string]string) (interface{}, error) {
	process := func(response *http.Response) (interface{}, error) {
		defer response.Body.Close()

		if contents, err := ioutil.ReadAll(response.Body); err != nil {
			return nil, err
		} else {
			var jsonContents map[string]interface{}
			json.Unmarshal(contents, &jsonContents)

			if response.StatusCode/100 != 2 {
				fmt.Println(string(contents))
				return nil, errors.New("Unable to create invoice")
			} else {
				return jsonContents["data"], nil
			}
		}
	}

	if response, err := bc.Post(path, payload); err != nil {
		return nil, err
	} else if data, err := process(response); err != nil {
		return nil, err
	} else {
		return data, nil
	}
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
