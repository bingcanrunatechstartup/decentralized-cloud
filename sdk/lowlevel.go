package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	host "github.com/libp2p/go-libp2p-core/host"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	swarm "github.com/libp2p/go-libp2pd/swarm"

	ma "github.com/multiformats/go-multiaddr"
)

const bidProtocolID = protocol.ID("/bid/1.0.0")

type BidMessage struct {
	ProviderID int     `json:"provider_id"`
	BidPrice   float64 `json:"bid_price"`
}

type BidResponse struct {
	BidAccepted bool `json:"bid_accepted"`
}

func (sdk *CloudSDK) Bid(providerID int, bidPrice float64) (bool, error) {
	ctx := context.Background()

	privKey, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return false,err
	}

	hst, err := libp2pp.New(ctx,
		libP2P.Identity(privKey),
	)
	if err != nil {
		return false,err
	}
	defer hst.Close()

	data,err:= ioutil.ReadFile("<BootstrapListPath>")
	if err!=nil{
	   return false,err 
    }
	var addrs []string 
	err = json.Unmarshal(data,&addrs)
	if err!=nil{
	   return false,err 
    }

	for _,addrStr:= range addrs{
	    addr,err:= ma.NewMultiaddr(addrStr)
	    if err!=nil{
	       continue  
	    }
	    pi,err:= pstore.InfoFromP2PAddr(addr)
	    if err!=nil{
	       continue  
	    }
	    hst.Peerstore().AddAddrs(pi.ID,pi.Addrs,pstore.PermanentAddrTTL)
	    
	    s,err:= hst.NewStream(ctx,bidProtocolID)
	    if err!=nil{
	       continue  
	    }
	    
        bidMsg := &BidMessage{ProviderID: providerID,BidPrice: bidPrice}
        data,_= json.Marshal(bidMsg)

        _,err = s.Write(data)
        if err!=nil{
            s.Reset()
            continue 
        }

        data,_= ioutil.ReadAll(s)

        var resp BidResponse
        json.Unmarshal(data,&resp)

        return resp.BidAccepted,nil 

    }

	return false,nil
}