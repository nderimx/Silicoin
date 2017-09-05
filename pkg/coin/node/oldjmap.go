//this is made so the program doesn't brake on go 1.6x and below because of the older json package
package node
import ("fmt"
		"strconv")
func ConvMap(input map[string]string) map[int]string{
	iLen:=len(input)
	integers:=make([]int, iLen)
	stringers:=make([]string, iLen)
	output:=make(map[int]string, iLen)
	i:=0
	for key, value:=range input{
		integers[i], _=strconv.Atoi(key)
		stringers[i]=value
		i++
	}
	for i:=0; i<iLen; i++{
		output[integers[i]]=stringers[i]
	}
	return output
}