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

	p1 pb.GameSerivce_PlayServer
	p2 pb.GameSerivce_PlayServer
}

func (s *server) Play(stream pb.GameSerivce_PlayServer) error {
	fmt.Println("Игра началась!")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return err
		}

		switch x := req.Action.(type) {

		case *pb.PlayerAction_Join:
			fmt.Println("Присоединился игрок:", x.Join)

			if s.p1 == nil {
				s.p1 = stream
				stream.Send(&pb.ServerResponse{
					Response: &pb.ServerResponse_Init{
						Init: &pb.InitInfo{YourPlayer: 1},
					},
				})
				stream.Send(&pb.ServerResponse{
					Response: &pb.ServerResponse_WaitForSecond{
						WaitForSecond: "Ты Игрок 1 (X). Ждём второго игрока...",
					},
				})
			} else if s.p2 == nil {
				s.p2 = stream
				stream.Send(&pb.ServerResponse{
					Response: &pb.ServerResponse_Init{
						Init: &pb.InitInfo{YourPlayer: 2},
					},
				})

				msg := &pb.ServerResponse{
					Response: &pb.ServerResponse_Board{
						Board: &pb.Board{
							Rows:          s.g.Board,
							Turn:          "Игра началась. Ход Игрока 1",
							CurrentPlayer: int32(s.g.Player),
						},
					},
				}
				if s.p1 != nil {
					s.p1.Send(msg)
				}
				if s.p2 != nil {
					s.p2.Send(msg)
				}
			} else {
				stream.Send(&pb.ServerResponse{
					Response: &pb.ServerResponse_GameOver{
						GameOver: "Комната занята, тут уже играют два игрока.",
					},
				})
				return nil
			}

		case *pb.PlayerAction_Move:
			var player int
			switch stream {
			case s.p1:
				player = 1
			case s.p2:
				player = 2
			default:
				stream.Send(&pb.ServerResponse{
					Response: &pb.ServerResponse_GameOver{
						GameOver: "Сначала нужно Join.",
					},
				})
				return nil
			}

			if player != s.g.Player {
				resp := &pb.ServerResponse{
					Response: &pb.ServerResponse_Board{
						Board: &pb.Board{
							Rows:          s.g.Board,
							Turn:          fmt.Sprintf("Сейчас ход игрока %d, подожди.", s.g.Player),
							CurrentPlayer: int32(s.g.Player),
						},
					},
				}
				if s.p1 != nil {
					s.p1.Send(resp)
				}
				if s.p2 != nil {
					s.p2.Send(resp)
				}
				continue
			}

			ok := s.g.MakeMove(int(x.Move.X), int(x.Move.Y))
			if !ok {
				resp := &pb.ServerResponse{
					Response: &pb.ServerResponse_Board{
						Board: &pb.Board{
							Rows:          s.g.Board,
							Turn:          fmt.Sprintf("Некорректный ход: (%d,%d)", x.Move.X, x.Move.Y),
							CurrentPlayer: int32(s.g.Player),
						},
					},
				}
				if s.p1 != nil {
					s.p1.Send(resp)
				}
				if s.p2 != nil {
					s.p2.Send(resp)
				}
				continue
			}

			winner := s.g.Winner()
			if winner != 0 || s.g.Draw() {
				boardResp := &pb.ServerResponse{
					Response: &pb.ServerResponse_Board{
						Board: &pb.Board{
							Rows:          s.g.Board,
							Turn:          "",
							CurrentPlayer: int32(s.g.Player),
						},
					},
				}
				if s.p1 != nil {
					s.p1.Send(boardResp)
				}
				if s.p2 != nil {
					s.p2.Send(boardResp)
				}

				msg := "Ничья. Игра окончена"
				if winner != 0 {
					msg = fmt.Sprintf("Поздравляем, игрок %d победил!", winner)
				}
				resp := &pb.ServerResponse{
					Response: &pb.ServerResponse_GameOver{
						GameOver: msg,
					},
				}
				if s.p1 != nil {
					s.p1.Send(resp)
				}
				if s.p2 != nil {
					s.p2.Send(resp)
				}
				return nil
			}

			resp := &pb.ServerResponse{
				Response: &pb.ServerResponse_Board{
					Board: &pb.Board{
						Rows:          s.g.Board,
						Turn:          fmt.Sprintf("Игрок %d сделал ход", s.g.Player),
						CurrentPlayer: int32(s.g.Player),
					},
				},
			}
			if s.p1 != nil {
				s.p1.Send(resp)
			}
			if s.p2 != nil {
				s.p2.Send(resp)
			}
			s.g.NextPlayer()
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
	fmt.Println("Сервер запущен на порту 50051...")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v\n", err)
	}
}
