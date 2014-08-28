package main


var _ = fmt.Print // debugging

func main() {
        var tree int
        var apple int
        apple = 6
        tree = 2 + apple + 3
        if tree == 11 {
        	apple = 2
        }
        return apple
}
