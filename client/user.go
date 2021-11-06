package client

import (
	"fmt"
	"log"
)

type User struct {
	Id       string `mikrotik:".id"`
	Address  string `mikrotik:"address"`
	Comment  string `mikrotik:"comment"`
	Disabled bool   `mikrotik:"disabled"`
	Expired  bool   `mikrotik:"expired"`
	Group    string `mikrotik:"group"`
	Name     string `mikrotik:"name"`
}

func (client Mikrotik) AddUser(user *User) (*User, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := Marshal("/user/add", user)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] ip address creation response: `%v`", r)

	if err != nil {
		return nil, err
	}

	id := r.Done.Map["ret"]

	return client.FindUser(id)
}

func (client Mikrotik) ListUser() ([]User, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}
	cmd := []string{"/user/print"}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] found users: %v", r)

	usr := []User{}

	err = Unmarshal(*r, &usr)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (client Mikrotik) FindUser(id string) (*User, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := []string{"/user/print", "?.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	r, err := c.RunArgs(cmd)

	log.Printf("[DEBUG] ip address response: %v", r)

	if err != nil {
		return nil, err
	}

	usr := User{}
	err = Unmarshal(*r, &usr)

	if err != nil {
		return nil, err
	}

	if usr.Id == "" {
		return nil, NewNotFound(fmt.Sprintf("User `%s` not found", id))
	}

	return &usr, nil
}

func (client Mikrotik) UpdateUser(user *User) (*User, error) {
	c, err := client.getMikrotikClient()

	if err != nil {
		return nil, err
	}

	cmd := Marshal("/user/set", user)
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)

	if err != nil {
		return nil, err
	}

	return client.FindUser(user.Id)
}

func (client Mikrotik) DeleteUser(id string) error {
	c, err := client.getMikrotikClient()

	if err != nil {
		return err
	}

	cmd := []string{"/user/remove", "=.id=" + id}
	log.Printf("[INFO] Running the mikrotik command: `%s`", cmd)
	_, err = c.RunArgs(cmd)
	return err
}
