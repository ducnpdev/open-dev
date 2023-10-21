package main

import (
	"fmt"
	"net/http"

	"github.com/sclevine/agouti"
)

// /opt/homebrew/bin/chromedriver
func main() {

	// Start a new WebDriver session
	driver := agouti.ChromeDriver(
		// agouti.ChromeOptions("args", []string{
		// 	"--remote-debugging-port=9222", // Use the same port as the CDP address
		// }),
		agouti.ChromeOptions("args",
			[]string{"--remote-debugging-port=9222"},
		),
		// agouti.ChromeOptions("", ""),
		agouti.Browser("safari"),
	)

	// agouti.EdgeDriver()

	if err := driver.Start(); err != nil {
		fmt.Printf("Failed to start driver: %v\n", err)
		return
	}
	defer driver.Stop()

	// Open a new page
	page, err := driver.NewPage()
	if err != nil {
		fmt.Printf("Failed to open page: %v\n", err)
		return
	}
	// Navigate to a website
	err = page.Navigate("https://betterchoice.vn/chi-tiet-de-cu/ngan-hang-hdbank-137.htm")
	if err != nil {
		fmt.Printf("Failed to navigate to the website: %v\n", err)
		return
	}

	errC := page.SetCookie(&http.Cookie{

		Name:   "jXt17L2CndBMMbxFSLr7708a9231dfb566",
		Value:  "EAAAAOmeRpf48rBDeJ8QLKlCwif2lPlbahDDrNMr3/qRwXVdaVmgWrafeB6wV+x9dR+tHbwwNj2HBEVqJ4EaSTsM3SiUd1uZNP+ZreC/sMSCLNQmJ9xzwI7NskQoMvOc0xF5Xs70+sdAeOduVwLCAo8Dq5sYoOg5Wru2zgwKlpWVsJt5O0r8F3saI/lYeQo2R9W70U8ezoTm8VRFpWr43jtTLI3NMVh1bJ4DJl71B+Z/rdiAgLpxz8EJoi76Y5rmHzFzLC3v23xdZtruji/4y18krNwzf1xtqM0pFER68H2gfkqievCd0+Z/wp5LQ0OivKLdtcSR9z7ox6DbikeK005E95LGaLdOwMRQeXJyi4hwU4knJsssZAOeEIfcrLHzDd78+w==",
		Path:   "/",
		Domain: "api.betterchoice.vn",
	})
	if errC != nil {
		fmt.Printf("Failed to navigate to the website: %v\n", errC)
		return
	}
	errC1 := page.SetCookie(&http.Cookie{
		Name:  "SIDCC",
		Value: "ACA-OxMWoMKjW9ZAJLaIosA_Kp0FY3HEUQ3xIicKeMI8zRMcfyOmUWZAw29Opf2DsSSFz1Ed",
	})
	if errC1 != nil {
		fmt.Printf("Failed to navigate to the website: %v\n", errC1)
		return
	}

	errC2 := page.SetCookie(&http.Cookie{
		Name:  "ASP.NET_SessionId",
		Value: "2da4vfhhjzwfg1driz32l1eg",
	})
	if errC2 != nil {
		fmt.Printf("Failed to navigate to the website: %v\n", errC2)
		return
	}

	errC3 := page.SetCookie(&http.Cookie{
		Name:  "ASP.NET_SessionId",
		Value: "mkxadwyxukyvz0n5sqxo4zgx",
	})
	if errC3 != nil {
		fmt.Printf("Failed to navigate to the website: %v\n", errC3)
		return
	}

	fmt.Println("Webpage opened successfully")

	// selected := page.AllByClass("btn-vote js-vote-action")
	// ac, err := selected.Active()
	// fmt.Println("1111:", ac, err)
	login(page)
	// for i := 0; i < 1; i++ {
	// 	click(page)
	// }

	// popupPage, err := driver.NewPage(agouti.Browser("chrome"), agouti.Page)
	// if err != nil {
	// 	fmt.Printf("Failed to switch to the popup window: %v\n", err)
	// 	return
	// }

	// Now you can interact with elements within the popup
	// For example, you can click a button in the popup
	// popupButton := popupPage.FindByClass("buttonCSS")
	// if err := popupButton.Click(); err != nil {
	// 	fmt.Printf("Failed to click the popup button: %v\n", err)
	// 	return
	// }

	// clErr := page.FindByClass("swal2-height-auto").FindByClass("swal2-actions")
	// if clErr != nil {
	// 	fmt.Println("213:", clErr)
	// }
	// Locate the button element by its CSS selector
	ch := make(chan bool)
	<-ch
	// Perform further actions on the webpage here
}
func login(page *agouti.Page) {
	ele := "btn-login"
	button := page.FindByClass(ele)
	if button == nil {
		fmt.Printf("Button element not found\n")
		return
	}
	err := button.Click()
	if err != nil {
		fmt.Println("login err", err)
		panic(err)
	}

	page1 := page.FindByClass("popup-login-content")
	if page1 == nil {
		fmt.Println("page1 nil")
	} else {
		fmt.Println("page1 ok")

		clik1 := page1.Find("#action-link")
		fmt.Println(clik1)

		clik := page1.FindByClass("action-link")

		fmt.Println(clik)

		jsCode := `document.getElementsByClassName("action-link")[1].click();`
		if err := page.RunScript(jsCode, nil, nil); err != nil {
			fmt.Printf("Failed to execute the JavaScript: %v\n", err)
			return
		}

		text := page.FindByName("identifier")
		e := text.Fill("ducios123123@gmail.com")
		fmt.Println("Fill", e)
		e = page.FindByID("identifierNext").Click()
		fmt.Println("identifierNext", e)

	}

	// e := page.Navigate("https://accounts.google.com/v3/signin/identifier?opparams=%253F&dsh=S-1676679792%3A1697558810060291&client_id=899054189562-jvqpb80b624vhrm96lsbiabjiqu723m2.apps.googleusercontent.com&o2v=1&redirect_uri=https%3A%2F%2Fapi.betterchoice.vn%2FLogin%2FGg%2FExtractRedirect.aspx&response_type=token&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&service=lso&theme=glif&flowName=GeneralOAuthFlow&continue=https%3A%2F%2Faccounts.google.com%2Fsignin%2Foauth%2Fconsent%3Fauthuser%3Dunknown%26part%3DAJi8hAN-a2rQbjy4vXSZ3Ug-gHD4cpv4na2vhuUVdt84P3ADUaKv3UqeepN0byyftUj1D_7KIJAM__1e9RM1H_ttUqUaoqvyuXawNTLRfZLXf9_ES9t4-yotqBgKAdy5TmKKhNTaLAVwe0F_8Oc8l9t-tfDiq6Pxy0IRXHqXO_rR6OnqT1TdgBfEdxkZW_gDuqiVv6fduRPZH1ixRLATCUESyIJJ-U9AzuM9MHsH61z3RHDeCbvoDA-CRDJ2CA74leEhbUf6tohtYQH1qtbb2JVdXiGsR83kWz_izr7pBVZQuMPOAnzwyGNapa_7Td4jm7dQ6oyl_PvYScysfbXWGpB9FX7HfJ0eaniijR7eamSWnIU7sNtm1zkXmg7QGPvdyexmeZnoLcbP5znVx-IVusXqnDRKDp9wRuTQwI_kX790mgtaHBN4dD36D-W9-S4WsxUQL2AVx8kaD6BLjD9x5T8tl1ZLO57ZJ1c8divVr8RCG2Opu5Ot-h0%26as%3DS-1676679792%253A1697558810060291%26client_id%3D899054189562-jvqpb80b624vhrm96lsbiabjiqu723m2.apps.googleusercontent.com%26theme%3Dglif%23&app_domain=https%3A%2F%2Fapi.betterchoice.vn&rart=ANgoxcdm95kFcwQDZ8DIE-tYKWpraD0LRmDJR9YjF_TB8vVqslYbPI2MPe9mTUZS-8xlQ__FQLNOAjddQ6gLQp_gBnlfgmVQ9Q")
	// fmt.Println("Navigate", e)
	// text := page.FindByName("identifier")
	// e = text.Fill("ducios123123@gmail.com")
	// fmt.Println("Fill", e)

	// e = page.FindByID("identifierNext").Click()
	// fmt.Println("identifierNext", e)
	fmt.Println("login")
}
func click(page *agouti.Page) {
	if page == nil {
		return
	}
	buttonSelector := "detail-header" // Replace with the actual CSS selector for your button
	// page.FindByLink()
	button := page.FindByClass(buttonSelector)
	if button == nil {
		fmt.Printf("Button element not found\n")
		return
	}
	buttonSelector1 := "btn-vote" // Replace with the actual CSS selector for your button

	mul := button.FindByClass(buttonSelector1)
	errC := mul.Click()
	if errC != nil {
		fmt.Printf("errC: %v\n", errC)
		panic(errC)
	}
	// Click the button
	err := button.Click()
	if err != nil {
		fmt.Printf("Failed to click the button: %v\n", err)
		panic(err)
	}

	// a := "swal2-confirm"
	// mul1 := button.FindByClass(a)
	// errC1 := mul1.Click()
	// if errC1 != nil {
	// 	fmt.Println(errC1)
	// 	panic(errC1)
	// }
	fmt.Println("Button clicked successfully")
}
