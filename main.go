package main

import (
	"flag"
	"github.com/justuswilhelm/gotogo/lib"
	"log"
)

var (
	blackCommand = ""
	whiteCommand = ""
	black        *lib.Process
	white        *lib.Process
)

func init() {
	flag.StringVar(&blackCommand, "b", "gnugo --mode gtp", "Black Command")
	flag.StringVar(&whiteCommand, "w", "gnugo --mode gtp", "White Command")
}

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

func create() error {
	var err error
	black, err = lib.CreateProcess("black", blackCommand)
	if err != nil {
		return err
	}
	white, err = lib.CreateProcess("white", whiteCommand)
	if err != nil {
		return err
	}
	if err := black.StartProcess(); err != nil {
		return err
	}
	if err := white.StartProcess(); err != nil {
		return err
	}
	if err := sanityCheck(black); err != nil {
		return err
	}
	if err := sanityCheck(white); err != nil {
		return err
	}
	if err := initialize(black, 9, "5.5"); err != nil {
		return err
	}
	if err := initialize(white, 9, "5.5"); err != nil {
		return err
	}
	return nil
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
		} else {
			firstPass = false
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
		} else {
			firstPass = false
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
	flag.Parse()

	if err := create(); err != nil {
		log.Fatalf("Error when creating game: %+v", err)
	}
	if err := play(black, white); err != nil {
		log.Fatalf("Error when playing game: %+v", err)
	}
}
