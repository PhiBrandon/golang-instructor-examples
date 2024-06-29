package input

import "fmt"

type Prompt struct {
	Instruction string
	Input       string
}

func (p *Prompt) Createprompt() string {
	return fmt.Sprintf("Instructions:\n%v\n\nInput:\n%v", p.Instruction, p.Input)
}
