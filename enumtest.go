// enumtest

package main

//	"fmt"

// 实现枚举例子

type State int

// iota 初始化后会自动递增
const (
	Running    State = iota // value --> 0
	Stopped                 // value --> 1
	Rebooting               // value --> 2
	Terminated              // value --> 3
)

func (tthis State) String() string {
	switch tthis {
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	default:
		return "Unknow"
	}
}

// 输出 state Running
// 没有重载String函数的情况下则输出 state 0
// func main() {
// 	state := Stopped
// 	fmt.Println("state", state)
// }
