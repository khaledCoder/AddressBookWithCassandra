package models

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/revel/revel"
	"time"
)
var cluster *gocql.ClusterConfig
type Users struct{
	Pk gocql.UUID
	Name string
	Password string
	Email string
	UserName string
}

type Contact struct{
	Pk gocql.UUID
	ContactName string
	ContactPhone string
	UserID gocql.UUID
}
type page struct{
	Contacts[]Contact
}
func init() {
	revel.OnAppStart(func() {
		cluster = gocql.NewCluster("127.0.0.1")
		cluster.Keyspace = "contacts"
		cluster.Consistency = gocql.One
	})
}


func CreateContact(contact Contact){
	var userName string
	Session, _ := cluster.CreateSession()
	defer Session.Close()
	err:=  Session.Query("SELECT user_name FROM users WHERE user_id=?",contact.UserID).Scan(&userName);
	checkErr(err)
	err = Session.Query("INSERT INTO contacts_by_user (user_id,added_time,contact_id,contact_name,contact_phone,user_name) VALUES (?,toTimestamp(now()),uuid(),?,?,? )",contact.UserID,contact.ContactName,contact.ContactPhone,userName).Exec();
	checkErr(err)
}
func FindAllContacts(id gocql.UUID) page{
	p:=page{Contacts:[]Contact{}}
	Session, _ := cluster.CreateSession()
	defer Session.Close()
	 iter := Session.Query("SELECT contact_id,user_id,contact_name, contact_phone FROM contacts_by_user WHERE user_id=? ",id).Iter();
	var b Contact
	for iter.Scan(&b.Pk,&b.UserID,&b.ContactName,&b.ContactPhone) {
		p.Contacts=append(p.Contacts,b)
	}
	err := iter.Close();
	checkErr(err)
	return p
}

func FindContact(contactID gocql.UUID,userID gocql.UUID) Contact{
	session, _ := cluster.CreateSession()
	defer session.Close()
	var contact Contact
	err := session.Query("SELECT contact_id , contact_name , contact_phone FROM contacts_by_user WHERE user_id=? and contact_id=? ALLOW FILTERING",userID,contactID).Scan(&contact.Pk,&contact.ContactName,&contact.ContactPhone);
	checkErr(err)
	return contact
}
func FindAllContactsByContactName(user_id gocql.UUID ,contactName string) page{
	session, _ := cluster.CreateSession()
	p:=page{Contacts:[]Contact{}}
	defer session.Close()
	iter := session.Query("SELECT contact_name, contact_phone FROM contacts_by_name WHERE user_id=? and contact_name=?",user_id,contactName).Iter();
	var b Contact
	for iter.Scan(&b.ContactName,&b.ContactPhone) {
		p.Contacts=append(p.Contacts,b)
	}
	err := iter.Close();
	checkErr(err)
	return p
}
func FindAllContactsByContactPhone(user_id gocql.UUID ,contactPhone string) page{
	// session to manage connection to the cluster
	session, _ := cluster.CreateSession()
	p:=page{Contacts:[]Contact{}}
	defer session.Close()
	iter := session.Query("SELECT contact_name, contact_phone FROM contacts_by_phone WHERE user_id=? and contact_phone=?",user_id,contactPhone).Iter();
	var b Contact
	for iter.Scan(&b.ContactName,&b.ContactPhone) {
		p.Contacts=append(p.Contacts,b)
	}
	 err := iter.Close();
	checkErr(err)
	return p
}
func FindUser(username string , password string) Users{
	// session to manage connection to the cluster
	session, _ := cluster.CreateSession()
	var userData Users
	err := session.Query("SELECT * FROM users WHERE user_name=? and password=? ALLOW FILTERING",username,password).Scan(&userData.Pk, &userData.Name,&userData.Password);
	checkErr(err)
	return userData
}

func UpdateContactInfo(contact Contact){
	// session to manage connection to the cluster
	session, _ := cluster.CreateSession()
	defer session.Close()
	var addedTime time.Time
	err := session.Query("SELECT added_time  FROM contacts_by_user WHERE user_id=? and contact_id=? ALLOW FILTERING",contact.UserID,contact.Pk).Scan(&addedTime);
	checkErr(err)
	fmt.Print("added time in update model:   ")
	fmt.Println(addedTime)
	err = session.Query("UPDATE contacts_by_user SET contact_name=?,contact_phone=? WHERE user_id=? and contact_id=? and added_time=?",contact.ContactName,contact.ContactPhone,contact.UserID,contact.Pk,addedTime).Exec();
	checkErr(err)
}

func DeleteContact(id gocql.UUID,userID gocql.UUID){
	session, _ := cluster.CreateSession()
	defer session.Close()
	 err := session.Query("DELETE FROM contacts_by_user WHERE contact_id=? and user_id=?",id,userID).Exec();
	checkErr(err)
}

func checkErr(err error)  {
	if err!=nil{
		fmt.Println(err.Error())
	}
}