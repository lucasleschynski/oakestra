package main

import (
	"../env"
	"../proxy"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("Create the goProxy tun device")
	myproxy := proxy.New()
	myproxy.Listen()
	errch := myproxy.GetErrCh()
	stopch := myproxy.GetStopCh()
	finishch := myproxy.GetFinishCh()

	config := env.Configuration{
		HostBridgeName:             "goProxyBridge",
		HostBridgeIP:               "172.19.1.1",
		HostBridgeMask:             "/24",
		HostTunName:                "goProxyTun",
		ConnectedInternetInterface: "wlan0",
	}

	time.Sleep(4 * time.Second)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Current Dev IP address for demonstrative purpose: \n")
	host1, _ := reader.ReadString('\n')
	host1 = strings.TrimSuffix(host1, "\n")
	fmt.Println("Current Host ip set to: ", host1)
	fmt.Print("Input Dev2 IP address for demonstrative purpose: \n")
	host2, _ := reader.ReadString('\n')
	host2 = strings.TrimSuffix(host2, "\n")
	fmt.Println("Dev2 Host ip set to: ", host2)

	//Cleanup and create a new environment
	myenv := env.NewCustom(myproxy.GetName(), config)
	myenv.Destroy()
	myenv = env.NewCustom(myproxy.GetName(), config)
	fmt.Println("Initial env: \n ", myenv)
	fmt.Println("Creating service 1 with ip 172.19.1.12 and namespace myapp1")
	ip1 := net.ParseIP("172.19.1.12")
	_, err := myenv.CreateNetworkNamespace("myapp1", ip1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Creating service 2 with ip 172.19.1.15 and namespace myapp2")
	ip2 := net.ParseIP("172.19.1.15")
	_, err = myenv.CreateNetworkNamespace("myapp2", ip2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Final env: \n ", myenv)

	myproxy.SetEnvironment(&myenv)

	//Setup custom name resolution
	myenv.AddTableQueryEntry(env.TableEntry{
		Appname:          "app1",
		Appns:            "default",
		Servicename:      "myapp2",
		Servicenamespace: "default",
		Instancenumber:   0,
		Cluster:          0,
		Nodeip:           net.ParseIP(host2),
		Nodeport:         50011,
		Nsip:             net.ParseIP("172.19.2.15"),
		ServiceIP: []env.ServiceIP{{
			IpType:  env.Closest,
			Address: net.ParseIP("172.30.0.0"),
		}},
	})
	myenv.AddTableQueryEntry(env.TableEntry{
		Appname:          "app1",
		Appns:            "default",
		Servicename:      "myapp2",
		Servicenamespace: "default",
		Instancenumber:   0,
		Cluster:          0,
		Nodeip:           net.ParseIP(host1),
		Nodeport:         50011,
		Nsip:             net.ParseIP("172.19.1.15"),
		ServiceIP: []env.ServiceIP{{
			IpType:  env.Closest,
			Address: net.ParseIP("172.30.0.0"),
		}},
	})
	myenv.AddTableQueryEntry(env.TableEntry{
		Appname:          "app1",
		Appns:            "default",
		Servicename:      "myapp1",
		Servicenamespace: "default",
		Instancenumber:   0,
		Cluster:          0,
		Nodeip:           net.ParseIP(host2),
		Nodeport:         50011,
		Nsip:             net.ParseIP("172.19.2.12"),
		ServiceIP: []env.ServiceIP{{
			IpType:  env.Closest,
			Address: net.ParseIP("172.30.0.1"),
		}},
	})
	myenv.AddTableQueryEntry(env.TableEntry{
		Appname:          "app1",
		Appns:            "default",
		Servicename:      "myapp1",
		Servicenamespace: "default",
		Instancenumber:   0,
		Cluster:          0,
		Nodeip:           net.ParseIP(host1),
		Nodeport:         50011,
		Nsip:             net.ParseIP("172.19.1.12"),
		ServiceIP: []env.ServiceIP{{
			IpType:  env.Closest,
			Address: net.ParseIP("172.30.0.1"),
		}},
	})

	//listen tun device
	cherror := <-errch
	<-finishch
	stopch <- true
	log.Fatal(cherror)
}
