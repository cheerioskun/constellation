package ipmi

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"time"

	expect "github.com/google/goexpect"
)

var (
	SmcIPMIToolJarPath = "/usr/local/bin/smcipmitool"
)

type Client struct {
	session expect.Expecter
}

func SetJarPath(path string) {
	SmcIPMIToolJarPath = path
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(ipAddress, username, password string) error {
	log.Printf("Connecting to %s with username %s", ipAddress, username)
	var err error
	libPath := path.Dir(SmcIPMIToolJarPath)
	cmd := fmt.Sprintf("java -Djava.library.path=%s -jar %s %s %s %s shell",
		libPath,
		SmcIPMIToolJarPath,
		ipAddress,
		username,
		password,
	)
	c.session, _, err = expect.Spawn(cmd, 30*time.Minute)
	if err != nil {
		return err
	}
	// Wait for actual connection to be established
	_, _, err = c.session.Expect(regexp.MustCompile(ipAddress), 10*time.Second)
	return err
}

func (c *Client) Disconnect() {
	c.session.Close()
}

func (c *Client) MountISO(isoPath string) error {
	log.Printf("Waiting for prompt to mount ISO...")
	_, _, err := c.session.Expect(regexp.MustCompile("SIM(WA)"), 5*time.Second)
	if err != nil {
		return err
	}
	mountCmd := fmt.Sprintf("vmwa dev2iso %s", isoPath)
	err = c.session.Send(mountCmd)
	return err
}

func (c *Client) PowerCycle() error {
	log.Printf("Waiting for prompt to power cycle...")
	_, _, err := c.session.Expect(regexp.MustCompile("SIM(WA)"), 5*time.Second)
	if err != nil {
		return err
	}
	powerCycleCmd := "ipmi power reset"
	err = c.session.Send(powerCycleCmd)
	return err
}
