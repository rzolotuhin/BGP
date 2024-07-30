package main

import (
	"bgp/bird"
	"bgp/tools"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"time"
)

type ResourceList struct {
	Data struct {
		Resources struct {
			IPv4 []string
			IPv6 []string
		}
	}
}

var listPath = tools.GetEnvDefault("list_path", "ru-subnet.list")

func main() {
	bird.UnixSocketPath = tools.GetEnvDefault("bird_unix_socket_path", "/run/bird/bird.ctl")

	sleep := time.NewTicker(24 * time.Hour)
	c := make(chan bool, 1)

	c <- tools.GetEnvBool("force")
	for {
		select {
		case force := <-c:
			if !force {
				continue
			}
		case <-sleep.C:
		}

		if info, err := os.Stat(listPath); err != nil {
			log.Println(err)
			log.Println("generate new list")
		} else {
			if time.Since(info.ModTime()) < time.Hour {
				log.Printf("the %s file has not aged yet\n", listPath)
				continue
			}
		}

		if err := UpdateList(); err != nil {
			log.Println(err)
			continue
		}

		log.Println("wait 30s")
		time.Sleep(30 * time.Second)

		b, err := bird.Command("configure")
		if err != nil {
			log.Println(err)
		} else {
			log.Println(b)
		}
	}
}

func UpdateList() error {
	res, err := http.Get("https://stat.ripe.net/data/country-resource-list/data.json?resource=ru")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	list := ResourceList{}
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}

	log.Printf("\n - IPv4: %d\n",
		len(list.Data.Resources.IPv4),
	)

	f, err := os.Create(listPath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, subnet := range list.Data.Resources.IPv4 {
		_, err := netip.ParsePrefix(subnet)
		if err == nil {
			_, err = fmt.Fprintf(f, "route %s unreachable;\n", subnet)
			if err != nil {
				return err
			}
			continue
		}

		addr := strings.Split(subnet, "-")
		if len(addr) == 2 {
			a1 := net.ParseIP(addr[0])
			a2 := net.ParseIP(addr[1])
			if a1 != nil && a2 != nil {
				list, err := tools.IpRangeToCIDR(a1.String(), a2.String())
				if err == nil {
					for _, subnet := range list {
						_, err := fmt.Fprintf(f, "route %s unreachable;\n", subnet)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
