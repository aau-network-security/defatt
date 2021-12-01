package cli

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	pb "github.com/aau-network-security/defatt/app/daemon/proto"
	"github.com/spf13/cobra"
)

var (
	ErrUnableCreateUList = errors.New("failed to create users list")
	ErrPasswordsNoMatch  = errors.New("passwords do not match, so cancelling signup :-(")
)

func (c *Client) CmdUser() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Actions to perform on users",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		c.CmdInviteUser(),
		c.CmdUpdatePasswd(),
		c.CmdSignupUser(),
		c.CmdLoginUser(),
		c.CmdDestroyUser(),
		c.CmdListUsers())

	return cmd
}

func (c *Client) CmdInviteUser() *cobra.Command {
	var superUser bool
	var npUser bool
	cmd := &cobra.Command{
		Use:     "invite",
		Short:   "Create key for inviting other users (only admins)",
		Example: `defat user invite --superuser \ defat user invite --member`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := c.rpcClient.InviteUser(ctx, &pb.InviteUserRequest{SuperUser: superUser, NpUser: npUser})
			if err != nil {
				PrintError(err)
				return
			}

			if r.Error != "" {
				PrintError(fmt.Errorf(r.Error))
				return
			}

			fmt.Println(r.Key)
		},
	}

	cmd.Flags().BoolVarP(&superUser, "super-user", "s", false, "indicates if the sign up key will create a super user")
	cmd.Flags().BoolVarP(&npUser, "member", "m", false, "indicated if the sign up key will create a non-privileged user ")
	return cmd
}

func (c *Client) CmdSignupUser() *cobra.Command {
	return &cobra.Command{
		Use:     "signup",
		Short:   "Signup as user",
		Example: `defat user signup`,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				username  string
				name      string
				surname   string
				email     string
				signupKey string
			)
			// todo: should be improved !
			fmt.Print("Signup key: ")
			fmt.Scanln(&signupKey)

			fmt.Print("Username: ")
			fmt.Scanln(&username)

			fmt.Print("Name: ")
			fmt.Scanln(&name)

			fmt.Print("Surname: ")
			fmt.Scanln(&surname)

			fmt.Print("Email: ")
			fmt.Scanln(&email)

			password, err := ReadSecret("Password: ")
			if err != nil {
				log.Fatal("Unable to read password")
			}

			password2, err := ReadSecret("Password (again): ")
			if err != nil {
				log.Fatal("Unable to read password")
			}

			if password != password2 {
				PrintError(ErrPasswordsNoMatch)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.rpcClient.SignupUser(ctx, &pb.SignupUserRequest{
				Key:      signupKey,
				Username: username,
				Name:     name,
				Surname:  surname,
				Email:    email,
				Password: password,
			})
			if err != nil {
				PrintError(err)
				return
			}

			if r.Error != "" {
				PrintError(fmt.Errorf(r.Error))
				return
			}

			c.Token = r.Token
			if err := c.SaveToken(); err != nil {
				PrintError(err)
			}
		},
	}
}

func (c *Client) CmdLoginUser() *cobra.Command {
	return &cobra.Command{
		Use:     "login",
		Short:   "Login as user",
		Example: `defat user login`,
		Run: func(cmd *cobra.Command, args []string) {
			var username string
			fmt.Print("Username: ")
			fmt.Scanln(&username)

			password, err := ReadSecret("Password: ")
			if err != nil {
				log.Fatal("Unable to read password")
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.rpcClient.LoginUser(ctx, &pb.LoginUserRequest{
				Username: username,
				Password: password,
			})

			if err != nil {
				fmt.Println(err)
				return
			}

			if r.Error != "" {
				PrintError(fmt.Errorf(r.Error))
				return
			}

			c.Token = r.Token

			if err := c.SaveToken(); err != nil {
				PrintError(err)
			}
		},
	}
}

func (c *Client) CmdListUsers() *cobra.Command {

	return &cobra.Command{
		Use:     "list",
		Short:   "Lists available users [only admins]",
		Example: "defat user list",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := c.rpcClient.ListUsers(ctx, &pb.Empty{})

			if err != nil {
				fmt.Println(err)
				return
			}

			f := formatter{
				header: []string{"Username", "Name", "Surname", "Email", "Admin", "Created At"},
				fields: []string{"Username", "Name", "Surname", "Email", "IsSuperUser", "CreatedAt"},
			}
			var elements []formatElement
			for _, u := range r.Users {
				elements = append(elements, u)
			}

			table, err := f.AsTable(elements)
			if err != nil {
				PrintError(ErrUnableCreateUList)
				return
			}
			fmt.Print(table)

		},
	}
}

func (c *Client) CmdUpdatePasswd() *cobra.Command {

	var username string
	var password string
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Update password",
		Example: "defat user update --username <username> --password <password>",
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, e := c.rpcClient.ChangeUserPasswd(ctx, &pb.UpdatePasswdRequest{
				Username: username,
				Password: password,
			})
			if e != nil {
				PrintError(e)
				return
			}
			fmt.Println(r.Message)
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Supply password from the command line flag")

	return cmd

}

func (c *Client) CmdDestroyUser() *cobra.Command {
	var username string

	cmd := &cobra.Command{
		Use:     "dl",
		Short:   "Destroys the user information [only admins]",
		Example: "defat user dl --username <username>",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.rpcClient.DestroyUser(ctx, &pb.DestroyUserRequest{Username: username})
			if err != nil {
				PrintError(err)
				return
			}
			fmt.Println(r.Message)
		},
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "Used to destroy users ")

	return cmd
}
