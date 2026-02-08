package main

import (
	"fmt"
	"io"
	"net"

	"github.com/Dzetner/tic-tac-toe-grpc/game"
	pb "github.com/Dzetner/tic-tac-toe-grpc/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGameSerivceServer
	g *game.Game
}

func (s *server) Play(stream pb.GameSerivce_PlayServer) error {
	fmt.Println("Play starting..")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("Error receiving from stream: %v\n", err)
			return err
		}
		switch x := req.Action.(type) {
		case *pb.PlayerAction_Join:
			fmt.Println("Player joined: ", x.Join)
			stream.Send(&pb.ServerResponse{
				Response: &pb.ServerResponse_WaitForSecond{
					WaitForSecond: "Hello. Let`s wait for another player!",
				},
			})
		case *pb.PlayerAction_Move:
			flag := s.g.MakeMove(int(x.Move.X), int(x.Move.Y))
			if flag == false {
				stream.Send(&pb.ServerResponse{
					Response: &pb.ServerResponse_Board{
						Board: &pb.Board{
							Rows: s.g.Board,
							Turn: fmt.Sprintf("Incorrect turn: (%d,%d)", x.Move.X, x.Move.Y),
						},
					},
				})
			} else {
				winner := s.g.Winner()
				if winner != 0 || s.g.Draw() {
					stream.Send(&pb.ServerResponse{
						Response: &pb.ServerResponse_Board{
							Board: &pb.Board{
								Rows: s.g.Board,
								Turn: "",
							},
						},
					})
					msg := "Draw. Game Over"
					if winner != 0 {
						msg = fmt.Sprintf("Congratulations, Player %d WIN!", winner)
					}
					stream.Send(&pb.ServerResponse{
						Response: &pb.ServerResponse_GameOver{
							GameOver: msg,
						},
					})
					return nil
				}
				stream.Send(&pb.ServerResponse{
					Response: &pb.ServerResponse_Board{
						Board: &pb.Board{
							Rows: s.g.Board,
							Turn: fmt.Sprintf("Player %d make a move (%d,%d)", s.g.Player, x.Move.X, x.Move.Y),
						},
					},
				})
				s.g.NextPlayer()
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("Failed to listen: %v\n", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterGameSerivceServer(s, &server{
		g: game.NewGame(),
	})
	fmt.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v\n", err)
	}
}
