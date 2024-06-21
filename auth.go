package main

import (
  "context"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"

  "golang.org/x/oauth2"
)

// Recupera um token e retorna o cliente gerado.
func getClient(config *oauth2.Config) *http.Client {
  // O arquivo token.json guarda o token de acesso e 
  // atualização do usuário, e é criado automaticamente
  // quando o fluxo de autorização termina pela primeira vez.
  tokFile := "token.json"
  tok, err := tokenFromFile(tokFile)

  // Se o arquivo não existe, requisita o token da web.
  if err != nil {
    tok = getTokenFromWeb(config)

    // Salva token
    saveToken(tokFile, tok)
  }
  return config.Client(context.Background(), tok)
}

// Requisita um token da web, e retorna o token adquirido.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
  authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

  // Autorização da conta do Google
  fmt.Printf("Go to the following link in your browser then type the "+
    "authorization code: \n%v\n", authURL)

  var authCode string
  if _, err := fmt.Scan(&authCode); err != nil {
    // Código de autorização
    log.Fatalf("Unable to read authorization code: %v", err)
  }

  tok, err := config.Exchange(context.TODO(), authCode)
  if err != nil {
    log.Fatalf("Unable to retrieve token from web: %v", err)
  }
  return tok
}

// Obtem um token de um arquivo local
func tokenFromFile(file string) (*oauth2.Token, error) {
  f, err := os.Open(file)

  if err != nil {
    return nil, err
  }
  defer f.Close()

  tok := &oauth2.Token{}
  err = json.NewDecoder(f).Decode(tok)

  return tok, err
}

// Armazena o token em um arquivo
func saveToken(path string, token *oauth2.Token) {
  fmt.Printf("Saving credential file to: %s\n", path)
  f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
  if err != nil {
    log.Fatalf("Unable to cache oauth token: %v", err)
  }
  defer f.Close()
  json.NewEncoder(f).Encode(token)
}

