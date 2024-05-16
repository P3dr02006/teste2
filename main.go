package main

import (
	"fmt"
	
	 // Importe o pacote LDAP adequado
	
)

// Configurações do servidor LDAP
var (
	ldapServer   = "LDAP://10.0.9.56:389"
	ldapPort     = 389
	baseDN       = "dc=example,dc=com"
	bindUsername = "cn=admin,dc=example,dc=com"
	bindPassword = "password"
)

func login() bool {
	var username, password string
	fmt.Println("=============Login=============")
	fmt.Print("Username: ")
	fmt.Scan(&username)
	fmt.Print("Password: ")
	fmt.Scan(&password)

	// Conectando ao servidor LDAP
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		fmt.Printf("Erro ao conectar ao servidor LDAP: %v\n", err)
		return false
	}
	defer l.Close()

	// Bind com credenciais de admin
	err = l.Bind(bindUsername, bindPassword)
	if err != nil {
		fmt.Printf("Erro ao fazer bind no servidor LDAP: %v\n", err)
		return false
	}

	// Buscar o usuário no LDAP
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(uid=%s)", username),
		[]string{"dn"},
		nil,
	)

	searchResult, err := l.Search(searchRequest)
	if err != nil {
		fmt.Printf("Erro ao buscar usuário no LDAP: %v\n", err)
		return false
	}

	if len(searchResult.Entries) != 1 {
		fmt.Println("Usuário não encontrado ou múltiplos usuários encontrados")
		return false
	}

	userDN := searchResult.Entries[0].DN

	// Tentar autenticar o usuário com a senha fornecida
	err = l.Bind(userDN, password)
	if err != nil {
		fmt.Printf("Falha na autenticação: %v\n", err)
		return false
	}

	fmt.Println("==========Login Succes==========")
	return true
}

func main() {
	login()
}
