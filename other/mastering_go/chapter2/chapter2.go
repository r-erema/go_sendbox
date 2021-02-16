package chapter2

//#include <stdio.h>
//long sum(int n)
//{
//  if (n == 0)
//    return 1;
//  else
//    return(n + sum(n-1));
//}
import "C"

func sumC(n int) int {
	return int(C.sum(C.int(n)))
}

func sumGo(n int) int {
	if n == 0 {
		return 1
	} else {
		return n + sumGo(n-1)
	}
}
