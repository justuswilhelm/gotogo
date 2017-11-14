package main

import (
	"github.com/justuswilhelm/gotogo/lib"
	"log"
)

func sanityCheck(p *lib.Process) error {
	name, err := p.Name()
	if err != nil {
		return err
	}
	log.Printf("Name: %s", name)
	version, err := p.Version()
	if err != nil {
		return err
	}
	log.Printf("Version: %s", version)
	return nil
}

func logBoard(p *lib.Process) {
	board, err := p.ShowBoard()
	if err != nil {
		log.Fatalf("Error showing board: %+v", err)
	}
	log.Print(board)
}

func logScore(p *lib.Process) {
	score, err := p.FinalScore()
	if err != nil {
		log.Fatalf("Error retrieving score %+v", err)
	}
	log.Printf("Score: %s", score)
}

func initialize(p *lib.Process, boardsize int, komi string) error {
	if err := p.Boardsize(boardsize); err != nil {
		return err
	}
	if err := p.Komi(komi); err != nil {
		return err
	}
	if err := p.ClearBoard(); err != nil {
		return err
	}
	return nil
}

func create() (*lib.Process, *lib.Process) {
	black, err := lib.CreateGnuGo("black")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	white, err := lib.CreateGnuGo("white")
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = black.StartProcess()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = white.StartProcess()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := sanityCheck(black); err != nil {
		log.Fatalf("Black sanity check fatal: %+v", err)
	}
	if err := sanityCheck(white); err != nil {
		log.Fatalf("White sanity check fatal: %+v", err)
	}
	if err := initialize(black, 9, "5.5"); err != nil {
		log.Fatalf("Black sanity check fatal: %+v", err)
	}
	if err := initialize(white, 9, "5.5"); err != nil {
		log.Fatalf("White sanity check fatal: %+v", err)
	}
	return black, white
}

func isPass(move string) bool {
	return move == "PASS"
}

func play(black *lib.Process, white *lib.Process) error {
	firstPass := false
	for {
		blackMove, err := black.GenMove(lib.Black)
		if err != nil {
			return err
		}
		if isPass(blackMove) {
			if firstPass {
				break
			} else {
				firstPass = true
			}
		}
		logBoard(black)
		if err := white.Play(lib.Black, blackMove); err != nil {
			return err
		}

		whiteMove, err := white.GenMove(lib.White)
		if err != nil {
			return err
		}
		if isPass(whiteMove) {
			if firstPass {
				break
			} else {
				firstPass = true
			}
		}
		logBoard(white)
		if err := black.Play(lib.White, whiteMove); err != nil {
			return err
		}
	}

	logScore(black)
	logScore(white)

	if err := black.Close(); err != nil {
		return err
	}

	if err := white.Close(); err != nil {
		return err
	}

	return nil
}

func main() {
	black, white := create()
	if err := play(black, white); err != nil {
		log.Fatalf("Error when playing game: %+v", err)
	}
}
