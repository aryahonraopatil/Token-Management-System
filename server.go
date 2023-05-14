
package main

import (
	"context"
	"fmt"
	"log"
	"flag"
	"net"
	"math"
	"sync"
  	"encoding/binary"

	"google.golang.org/grpc"
	"crypto/sha256"

	pb "tokenmngr/proto" 
	
)

// generates a hash value
func Hash(name string, nonce uint64) uint64 {
	hashervalue := sha256.New()
	hashervalue.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hashervalue.Sum(nil))
}


type TokenMngrServer struct {
	mu sync.Mutex
	toks map[string]*Token  //stores all tokens
	pb.UnimplementedTokenManagerServer
	
}

func NewTokManager() *TokenMngrServer { 
	return &TokenMngrServer{
		toks: make(map[string]*Token),
	}
}

// creates a token struct
type Token struct {
	id     string
	name   string
	Domain struct {
		low  uint64
		mid  uint64
		high uint64
	}
	State struct {
		Partial uint64
		Final   uint64
	}
}


// New token creation
func (s *TokenMngrServer ) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := req.GetId()
	if _, ok := s.toks[id]; ok {
		return nil, fmt.Errorf("id %s token exists", id)
	}

	s.toks[id] = &Token{
         id: id,
    }
	
	s.Print(id)

	return &pb.CreateResponse{
		Success: true,
	}, nil

}


// Write Token
func (s *TokenMngrServer) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := req.GetId()
	tok, ok := s.toks[id]
	
	if !ok {
		return nil, fmt.Errorf("id %s token exists", id)
	}

	//if tok.Domain.low >= tok.Domain.mid || tok.Domain.mid >= tok.Domain.high {
       // return &pb.WriteResponse{success: false}, fmt.Errorf("invalid low/mid/high values for token with id %s", id)
    //}

	tok.name = req.GetName()
	tok.Domain.low = req.GetLow()
	tok.Domain.mid = req.GetMid()
	tok.Domain.high = req.GetHigh()

	var arg_min uint64 = tok.Domain.low
	var hash_min uint64 = math.MaxUint64
	
	for arg := tok.Domain.low; arg < tok.Domain.mid; arg++ {
		hash := Hash(tok.name, arg); 
		if hash < hash_min{
			arg_min = arg
			hash_min = hash
		}
	}

	tok.State.Partial = arg_min
	tok.State.Final = 0

	s.Print(id)

	return &pb.WriteResponse{ Partial: tok.State.Partial,}, nil
}

// Read Token 
func (s *TokenMngrServer) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	
	s.mu.Lock()
	defer s.mu.Unlock()

	id := req.GetId()
	tok, exists := s.toks[id]

	if !exists {
		return nil, fmt.Errorf("id %s exists", id)
		//use req.GetId() instead of id
	}

	var min_val uint64 = 0
	// put equal to 0 instead of tok.domain.low
	var minimum uint64= math.MaxUint64

	for arg := tok.Domain.mid; arg < tok.Domain.high; arg++{
		hash := Hash(tok.name, arg)
		if hash < minimum {
			minimum = hash
			min_val = arg
		}
	}

	Partial_hash := Hash(tok.name, tok.State.Partial)
	
	
	if minimum < Partial_hash{
		tok.State.Final = min_val
	}else {
		tok.State.Final = tok.State.Partial
	}

	s.Print(id)

	return &pb.ReadResponse{Final: tok.State.Final}, nil
}

// Delete token
func (s *TokenMngrServer ) Drop(ctx context.Context, req *pb.DropRequest) (*pb.DropResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := req.GetId()

	if _, ok := s.toks[id]; !ok {
		return nil, fmt.Errorf("id %s token exists", id)
	}

	delete(s.toks, id)
	s.Print(id)

	return &pb.DropResponse{
		Success: true,
	}, nil
}


//Print info of token
func (s *TokenMngrServer) Print(id string){

    f := s.toks[id] //get token from id
	//f, ok := s.toks[id]

	fmt.Println("Token Id: ",f.id)
	fmt.Println("Name: ", f.name)
	fmt.Println("Low: ", f.Domain.low)
	fmt.Println("Mid: ", f.Domain.mid)
	fmt.Println("High: ", f.Domain.high)
	fmt.Println("Partial: ",f.State.Partial)
	fmt.Println("Final: ", f.State.Final)

	//fmt.Println("Token Id: ",s.toks[id].id)
	//fmt.Println("Name: ", s.toks[id].name)
	//fmt.Println("Low: ",s.toks[id].Domain.low)
	//fmt.Println("Mid: ",s.toks[id].Domain.mid)
	//fmt.Println("High: ",s.toks[id].Domain.high)
	//fmt.Println("Partial: ",s.toks[id].State.Partial_value)
	//fmt.Println("Final: ",s.toks[id].State.Final_value)
	

	for i := range s.toks{
		fmt.Println(i)
	}

	return 
}



func main(){
	var port_num = flag.Int("port", 50051, "Server Port")
	
	flag.Parse()

	port:=fmt.Sprintf(":%d",*port_num)
	
	lis,err := net.Listen("tcp", port)
	if err!=nil{
		log.Fatalf("Listen Error: %v",err)
	}

	g:=grpc.NewServer()
    tok_man := NewTokManager()
	pb.RegisterTokenManagerServer(g, tok_man)

	log.Printf("Server started on %d", *port_num)

	if err := g.Serve(lis); err != nil{
		log.Fatalf("Serving failed: %v",err)
	}

}