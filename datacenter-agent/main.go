package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func findIPAdresses() (string, string){

	// Private IP
	var nodeIP net.IP
	ifaces, err := net.Interfaces()
	for _, iface := range ifaces {
		addresses, _ := iface.Addrs()
		key := iface.Name
		for _, address := range addresses {
			if strings.Contains(key, "en") || strings.Contains(key, "eth") || strings.Contains(key, "wl") {
				addIP, _, _ := net.ParseCIDR(address.String())
				if addIP.To4() != nil {
					nodeIP = addIP
				}
			}
		}
	}
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	privateIP := nodeIP.String()

	//2. Find Public IP

	resp, err := http.Get("http://ipinfo.io/ip")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	publicIP := responseBody


	return privateIP, publicIP
}

func main() {
	log.Println("INFO: Bootstrapping Local ENEDI Master")

	//1. Find Private and Public IP Addresses of the (Virtual) Machine
	privateIP, publicIP := findIPAdresses()
	log.Println("INFO: Private IP of the machine: "+privateIP)
	log.Println("INFO: Public IP of the machine: "+publicIP)

	//2. Check if netdata is installed
	log.Println("INFO: Checking if netdata is installed")
	var out bytes.Buffer
	if _, errexists := os.Stat("/opt/netdata/etc/netdata/netdata.conf"); errexists != nil {
		log.Println("INFO: Netdata is not installed. Installation will begin shortly...")

		//2a. Download the official script and run it
		cmd := exec.Command("bash", "-c", "curl -Ss 'https://my-netdata.io/kickstart-static64.sh' > /tmp/kickstart.sh")
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Could not execute curl")
		}
		syscall.Chmod("/tmp/kickstart.sh", 0777)
		cmd = exec.Command("/tmp/kickstart.sh", "--dont-wait")
		if err = cmd.Run(); err != nil {
			log.Printf("Error %v\n", "Could not start process")
		}
		/*ch := make(chan error)
		go func() {
			ch <- cmd.Wait()
		}()

		select {
		case <-time.After(10 * time.Second):
			if err = cmd.Process.Kill(); err != nil {
				log.Printf("Error %v\n", "Could not kill process")
			}
		case err = <-ch:
			if err != nil {
				log.Printf("Error 2 %v\n", "Could not kill process")
			}
		}
		if err != nil {
			log.Fatalf("Could not run script")
		}*/
		cmd = exec.Command("rm", "/tmp/kickstart.sh", "-rf")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Could not remove script")
		}

		log.Println("INFO: Getting netdata configuration")

		time.Sleep(2 * time.Second)
		//2b. Get Netdata Configuration and update the [backend] prefix
		cmd = exec.Command("bash", "-c", "sudo wget -O /opt/netdata/etc/netdata/netdata.conf 'http://localhost:19999/netdata.conf'")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Fatalf("ERROR: Could not load netdata configuration ")
		}

		//2c. Edit Netdata Configuration
		input, err := ioutil.ReadFile("/opt/netdata/etc/netdata/netdata.conf")
		if err != nil {
			log.Fatalf("ERROR: Reading netdata config: %v\n", err)
		}

		lines := strings.Split(string(input), "\n")

		for i, line := range lines {
			if strings.Contains(line, "# prefix = netdata") {
				lines[i] = "\tprefix = crete_dc1_master"
			}
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile("/opt/netdata/etc/netdata/netdata.conf", []byte(output), 0644)
		if err != nil {
			log.Fatalf("ERROR: Could not write netdata configuration: %v\n", err)
		}

		//2d. sudo service netdata restart
		cmd = exec.Command("bash","-c", "service netdata restart")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			log.Fatalf("ERROR: Restart failed :%v\n", err)
		}
		log.Println("INFO: Netdata is successfully installed")


	} else {
		log.Println("INFO: Netdata is installed")
	}

	//3. Check for Local Consul Master

	//4. Start Prometheus

	//5. Configure Prometheus

	//6. Insert Record to DB

	//7. Connect to remote consul with some tags


}