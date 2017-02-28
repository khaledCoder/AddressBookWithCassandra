package controllers

import (
	"github.com/revel/revel"
	"AddressBookWithCassandraRevel/app/models"
	"fmt"
	"github.com/gocql/gocql"
)
type App struct {
	*revel.Controller
}
var selectedUserID gocql.UUID
func (c App) Index(id string) revel.Result {
	u2, _ :=gocql.ParseUUID(id)
 var Contacts =models.FindAllContacts(u2)
	return c.Render(Contacts)
}

func (c App) Login() revel.Result{
	var userData models.Users
	var name , password string
	c.Params.Bind(&name,"name")
	c.Params.Bind(&password,"password")
	if name != "" && password != "" {
		// .. check credentials ..
		userData= models.FindUser(name,password)
		var emptyUUID gocql.UUID
		if userData.Pk!= emptyUUID{
			var id= userData.Pk
			selectedUserID = userData.Pk
			return c.Redirect("/App/Index/%s",id)
		}
	}
	return c.Render()
}

func (c App ) FindContactByContactID(id string)revel.Result{
	var contact models.Contact
	u2, _ :=gocql.ParseUUID(id)
 contact = models.FindContact(u2,selectedUserID)
	return c.RenderJson(contact)
}
func (c App) FindAllContactsByContactPhone() revel.Result{
	var id gocql.UUID
	var phone string
	 var contacts = models.FindAllContactsByContactPhone(id , phone)
	return c.Render(contacts)
}
func (c App) FindAllContactsByContactName() revel.Result{
	var id gocql.UUID
	var phone string
	var contacts = models.FindAllContactsByContactName(id , phone)
	return c.Render(contacts)
}
func (c App) SaveUserInfo()revel.Result{
	var contact  models.Contact
	contact.UserID=selectedUserID
	var hidden string
	c.Params.Bind(&hidden,"idUpdate")
	c.Params.Bind(&contact.ContactName,"name")
	c.Params.Bind(&contact.ContactPhone,"mobile")
	c.Validation.Required(contact.ContactName).Message("name is required!")
	c.Validation.MinSize(contact.ContactPhone, 11).Message("mobile number is not valid!")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect("/App/Index/%s",selectedUserID)
	}
	contact.Pk, _ =gocql.ParseUUID(hidden)
	var emptyUUID gocql.UUID
	if(contact.Pk!=emptyUUID){
		fmt.Println("inside update condition")
		models.UpdateContactInfo(contact)
	}else{
		models.CreateContact(contact)
}
	return c.Redirect("/App/Index/%s",selectedUserID)
}
func (c App) DeleteContact(id string) revel.Result{
	u2, _ :=gocql.ParseUUID(id)
	models.DeleteContact(u2,selectedUserID)
	return c.RenderJson(u2)
}