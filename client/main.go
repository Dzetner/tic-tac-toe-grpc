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

func readMove() (int32, int32) {
	for {
		var xMove, yMove int32
		fmt.Println("Введите ваш ход: два числа от 0 до 2, через пробел:")
		n, err := fmt.Scan(&xMove, &yMove)
		if err != nil || n != 2 {
			fmt.Println("Некорректный ввод, попробуйте ещё раз.")
			var dump string
			fmt.Scanln(&dump)
			continue
		}
		if xMove < 0 || xMove > 2 || yMove < 0 || yMove > 2 {
			fmt.Println("Координаты должны быть от 0 до 2. Попробуйте ещё раз.")
			continue
		}
		return xMove, yMove
	}
}

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Can't connect to server: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewGameSerivceClient(conn)
	ctx := context.Background()

	stream, err := c.Play(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var name string
	fmt.Println("Введите ваше имя:")
	fmt.Scan(&name)

	if err := stream.Send(&pb.PlayerAction{
		Action: &pb.PlayerAction_Join{
			Join: name,
		},
	}); err != nil {
		log.Fatal(err)
	}

	var yourPlayer int32 = 0
	needMove := false

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
		case *pb.ServerResponse_Init:
			yourPlayer = x.Init.YourPlayer
			fmt.Println("Ты Игрок", yourPlayer)

		case *pb.ServerResponse_Board:
			matrix := x.Board.Rows
			turn := x.Board.Turn
			current := x.Board.CurrentPlayer

			fmt.Println(turn)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					fmt.Printf("%s ", matrix[3*i+j])
				}
				fmt.Println()
			}

			if yourPlayer != 0 && yourPlayer == current {
				needMove = true
			} else {
				needMove = false
			}

		case *pb.ServerResponse_WaitForSecond:
			fmt.Println(x.WaitForSecond)
			needMove = false

		case *pb.ServerResponse_GameOver:
			fmt.Println(x.GameOver)
			return
		}

		if needMove {
			xMove, yMove := readMove()
			if err := stream.Send(&pb.PlayerAction{
				Action: &pb.PlayerAction_Move{
					Move: &pb.Move{
						X: xMove,
						Y: yMove,
					},
				},
			}); err != nil {
				log.Fatal(err)
			}
		}
	}
}
