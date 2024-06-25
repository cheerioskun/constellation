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
	Verbose            = false
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
	log.Printf("Running command: %s", cmd)
	c.session, _, err = expect.Spawn(cmd,
		30*time.Minute,
		expect.Verbose(Verbose),
		expect.CheckDuration(1*time.Second),
	)
	if err != nil {
		return err
	}
	// Wait for actual connection to be established
	var out string
	out, _, err = c.session.Expect(regexp.MustCompile(">"), 10*time.Second)
	log.Println(out)
	return err
}

func (c *Client) Disconnect() {
	c.session.Send("exit\n")
	c.session.Close()
}

func (c *Client) MountISO(isoPath string) error {

	mountCmd := fmt.Sprintf("vmwa dev2iso %s\n", isoPath)
	err := c.session.Send(mountCmd)
	if err != nil {
		return err
	}
	// Wait for the iso to be mounted: Plug-In OK!! message
	var out string
	out, _, err = c.session.Expect(regexp.MustCompile("OK"), 30*time.Second)
	log.Println(out)
	return err
}

func (c *Client) PowerCycle() error {
	powerCycleCmd := "ipmi power reset\n"
	err := c.session.Send(powerCycleCmd)
	if err != nil {
		return err
	}
	// Wait for the power to be reset: Done message
	var out string
	out, _, err = c.session.Expect(regexp.MustCompile("Done"), 10*time.Second)
	log.Println(out)
	return err
}
