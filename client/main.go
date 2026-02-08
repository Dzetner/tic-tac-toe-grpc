package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/Dzetner/tic-tac-toe-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Can`t listen to server:%v\n", err)
	}
	defer conn.Close()

	c := pb.NewGameSerivceClient(conn)

	ctx := context.Background()

	stream, err := c.Play(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err := stream.Send(&pb.PlayerAction{
		Action: &pb.PlayerAction_Join{
			Join: "Andrey"},
	}); err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Игра закончена. Спасибо!")
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		switch x := resp.Response.(type) {
		case *pb.ServerResponse_Board:
			matrix := x.Board.Rows
			turn := x.Board.Turn
			fmt.Println(turn)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					fmt.Printf("%s ", matrix[3*i+j])
				}
				fmt.Println()
			}

		case *pb.ServerResponse_WaitForSecond:
			ans := x.WaitForSecond
			fmt.Println(ans)
		case *pb.ServerResponse_GameOver:
			fmt.Println(x.GameOver)
			return
		}

		var xMove, yMove int32
		fmt.Println("Введите ваш ход")
		fmt.Scan(&xMove, &yMove)
		stream.Send(&pb.PlayerAction{
			Action: &pb.PlayerAction_Move{
				Move: &pb.Move{
					X: xMove,
					Y: yMove,
				},
			},
		})

	}
}
