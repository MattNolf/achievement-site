package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mehanizm/airtable"
)

var site = fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Achievement.dev</title>
	<!--Contact: hello@achievement.dev-->
	<style>
		.astext {
			background:none;
			border:none;
			margin:0;
			padding:0;
			cursor: pointer;
			color: blue;
			font-family:Inconsolata,monospace;
			display: inline;
		}
		p { display: inline; }
		button { display: inline; }
	</style>
</head>
<body onload="hasRegistered()">
	<img src="assets/logo.png" alt="Achievement.dev" width="100%%">
	<div style="text-align:center;font-family:Inconsolata,monospace;">	
		<p id="cs">Coming soon...</p>	
		<button id="registerButton" onclick="register()" class="astext">(Register interest)</button>
	</div>
	<script>
		function register() {
			let registeredEmail = prompt("Enter your email:", "you@email.com");
			if (registeredEmail == null || registeredEmail == "" || registeredEmail == "you@email.com") {
				return
			}
			
			document.getElementById("cs").innerHTML = "Thanks for registering interest " + registeredEmail;
		
			var x = document.getElementById("registerButton");
			x.style.display = "none";

			var now = new Date();
			var expireTime = now.getTime() + (3600 * 1000 * 24 * 365);
			now.setTime(expireTime);

			document.cookie = "registeredEmail" + "=" + registeredEmail + ";expires=" +now.toUTCString()

			var xhr = new XMLHttpRequest();
			xhr.open("POST", "%s/register_interest", true);
			xhr.setRequestHeader('Content-Type', 'application/json');
			xhr.send(JSON.stringify({
				email: registeredEmail
			}));
		}

		function hasRegistered() {
            var registeredEmail = document.cookie;

            if(registeredEmail!="") {
                document.getElementById("cs").innerHTML = "Thanks for registering interest " + registeredEmail.split('=')[1];
                var x = document.getElementById("registerButton");
				x.style.display = "none";
            }
        }

	</script>
</body>
</html>
`, os.Getenv("API_HOSTNAME"))

var client *airtable.Client
var table *airtable.Table

func main() {
	client = airtable.NewClient(os.Getenv("AIRTABLE_API_KEY"))
	table = client.GetTable(os.Getenv("AIRTABLE_DATABASE_ID"), os.Getenv("AIRTABLE_TABLE_ID"))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, site)
	}))
	http.Handle("/register_interest", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			Email string `json:"email"`
		}

		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			fmt.Println("Error on register: ", err)
		}

		fmt.Println("Register: ", body)

		go submitRegisterToAirtable(req.Context(), body.Email)
	}))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	http.ListenAndServe(":8080", nil)
}

func submitRegisterToAirtable(ctx context.Context, email string) error {
	rs := &airtable.Records{
		Records: []*airtable.Record{
			{
				Fields: map[string]interface{}{
					"email": email,
				},
			},
		},
	}

	_, err := table.AddRecords(rs)
	if err != nil {
		fmt.Println("Error persisting to Airtable: ", err)
	}
	return err
}
