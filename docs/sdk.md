# Low level call

Here are some example function signatures for your SDK library in Go that provide a one-to-one mapping between message types and functions:

```go
type CloudSDK struct {}

func (sdk *CloudSDK) QueryResources(providerID int) ([]Resource, error) {}
func (sdk *CloudSDK) Bid(providerID int, bidPrice float64) (bool, error) {}
func (sdk *CloudSDK) Reserve(providerID int, vmImage string, ansibleScript string) (ReservationDetails, error) {}
```

# High level call

Here's an example of how you might use this high-level function to reserve cloud resources:

```go
sdk := &CloudSDK{}
client := NewCloudClient(sdk)

req := ReservationRequest{
  ProviderID: <ProviderID>,
  BidPrice: <BidPrice>,
  VMImage: "<VMImage>",
  AnsibleScript: "<AnsibleScript>",
}

result,err := client.ReserveResources(req)
if err !=nil{
//handle error 
}
fmt.Println(result.IPAddress,result.WireGuardConfig)
```

Here's an example of how you might use this high-level function to enter into an SSH session on the provisioned VM:

```go
sdk := &CloudSDK{}
client := NewCloudClient(sdk)

req := ReservationRequest{
  ProviderID: <ProviderID>,
  BidPrice: <BidPrice>,
  VMImage: "<VMImage>",
  AnsibleScript: "<AnsibleScript>",
}

result,err:= client.ReserveResources(req)
if err !=nil{
//handle error 
}

privateKey:= "<PrivateKey>"

err=client.SSHToVM(result,privateKey)
if err!=nil{
//handle error 
}
```
