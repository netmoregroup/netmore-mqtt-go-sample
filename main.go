package main

import (
  "crypto/tls"
  "crypto/x509"
  "fmt"
  "io/ioutil"
  "log"
  "os"

  mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
  fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func main() {
  if len(os.Args) != 2 {
    log.Fatalf("Id from cert not given! use ca name as is. openssl x509 -in ./certs/<id>/client.crt -noout -text\nNote that the id need to match the cert common name, if not you will not receive any data.")
  }
  id := os.Args[1]

  url := "ssl://mqtt.netmoregroup.com:8883"

  caPath := fmt.Sprintf("./certs/%s/ca.crt", id)
  certPath := fmt.Sprintf("./certs/%s/client.crt", id)
  keyPath := fmt.Sprintf("./certs/%s/client.key", id)

  caCert, err := ioutil.ReadFile(caPath)
  if err != nil {
    log.Fatalf("Error reading CA certificate: %v", err)
  }

  clientCert, err := tls.LoadX509KeyPair(certPath, keyPath)
  if err != nil {
    log.Fatalf("Error loading client certificate: %v", err)
  }

  caCertPool := x509.NewCertPool()
  caCertPool.AppendCertsFromPEM(caCert)

  tlsConfig := &tls.Config{
    Certificates:       []tls.Certificate{clientCert},
    RootCAs:            caCertPool,
    InsecureSkipVerify: true,
    ClientAuth: tls.NoClientCert,
    ClientCAs: nil,
  }

  opts := mqtt.NewClientOptions()
  opts.AddBroker(url)
  opts.SetClientID("ssl-sample").SetTLSConfig(tlsConfig)
  opts.SetDefaultPublishHandler(messagePubHandler)

	// Start the connection
  c := mqtt.NewClient(opts)
  if token := c.Connect(); token.Wait() && token.Error() != nil {
    panic(token.Error())
  }
  topic := fmt.Sprintf("client/%s/#", id)
  c.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("TOPIC: %s\n", msg.Topic())
    fmt.Printf("MSG: %s\n", msg.Payload())
  })

  fmt.Printf("D: Waiting for MQTT at client/%s/#\n", id)

  select {}
}
