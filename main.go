package main

import (
	"bufio"
	"fmt"
	"github.com/restartfu/bedrock-porter/porter"
	"github.com/restartfu/bedrock-porter/porter/frontend"
	"os"
	"strings"
)

func main() {
	fmt.Print(frontend.Style.Render("Enter the path to the resource pack or drop file: "))
	path, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	path = strings.TrimSpace(path)

	fmt.Print(frontend.ClearTerminal, frontend.CursorUp, frontend.CursorUp)
	fmt.Print(frontend.Style.Render(fmt.Sprintf(" Porting %s...\n", path)))

	pack, err := porter.NewResourcePack(path)
	if err != nil {
		fmt.Println(err)
	}
	pack.Port()
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}
