type CloudClient struct {
	sdk *CloudSDK
  }
  
  func NewCloudClient(sdk *CloudSDK) *CloudClient {
	return &CloudClient{sdk: sdk}
  }
  
  type ReservationRequest struct {
	ProviderID int
	BidPrice float64
	VMImage string
	AnsibleScript string
  }
  
  type ReservationResult struct {
	IPAddress string
	WireGuardConfig string
  }
  
  func (c *CloudClient) ReserveResources(req ReservationRequest) (*ReservationResult,error){
	resources,err := c.sdk.QueryResources(req.ProviderID)
	if err != nil || len(resources)==0{
		return nil,err 
	 }
	 bidAccepted,err := c.sdk.Bid(req.ProviderID,req.BidPrice)
	 if err != nil || !bidAccepted{
		return nil,err 
	 }
	 reservationDetails,err := c.sdk.Reserve(req.ProviderID,req.VMImage,req.AnsibleScript)
	 if err != nil{
		return nil,err 
	 }
	 result := &ReservationResult{
		 IPAddress: reservationDetails.IPAddress,
		 WireGuardConfig: reservationDetails.WireGuardConfig,
	 }
	 return result,nil 
  }

  func (c *CloudClient) SSHToVM(result *ReservationResult,privateKey string) error{
	wgConfig,err := wgtypes.ParseConfig(result.WireGuardConfig)
	if err !=nil{
	   return err 
	}
	wgClient,err := wgctrl.New()
	if err !=nil{
	   return err 
	}
	ifaceName := "<InterfaceName>"
	_,err = net.InterfaceByName(ifaceName)
	if os.IsNotExist(err){
		cmd := exec.Command("ip","link","add",ifaceName,"type","wireguard")
		cmd.Run()
		cmd = exec.Command("ip","address","add",wgConfig.Addresses[0].String(),"dev",ifaceName)
		cmd.Run()
		cmd = exec.Command("ip","link","set",ifaceName,"up")
		cmd.Run()
	}
	
	wgClient.ConfigureDevice(ifaceName,*wgConfig)
 
	 signer,err := ssh.ParsePrivateKey([]byte(privateKey))
	 if err!=nil{
		 return 	err 
	 }
	 config := &ssh.ClientConfig{
		 User: "<Username>",
		 Auth: []ssh.AuthMethod{
			 ssh.PublicKeys(signer),
		 },
	 }
 
	 client,err:= ssh.Dial("tcp",result.IPAddress+":22",config)
	 if err!=nil{
		return 	err 
	 }
	 defer client.Close()
 
	 session,err:= client.NewSession()
	 if err!=nil{
		return 	err 
	 }
	 defer session.Close()
 
	 session.Stdout = os.Stdout
	 session.Stderr = os.Stderr
	 session.Stdin = os.Stdin
 
	 modes := ssh.TerminalModes{ssh.ECHO: 1}
 
	 err=session.RequestPty("xterm-256color",80,40,modes)
	 if err!=nil{
		return 	err 
	 }
 
	 err=session.Shell()
	 if err!=nil{
		return 	err 
	 }
 
	 session.Wait()
 
	 return nil 
 
 }
 