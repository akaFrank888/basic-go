package wire

import "fmt"

func UseRepository() {
	repo := initUserRepository()
	fmt.Print(repo)

}
