package main

import "github.com/JannoTjarks/azure-dyndns2/cmd"
import _ "github.com/JannoTjarks/azure-dyndns2/docs"

// @title		azure-dyndns2
// @version		0.1.1
// @description	Simple dyndns2-compatible web api for Azure DNS.

// @contact.name	Janno Tjarks
// @contact.url	https://tjarks.dev
// @contact.email	janno.tjarks@mailbox.org

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	cmd.Execute()
}
