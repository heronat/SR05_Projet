package utils
import "fmt"
import "os"
// Retourne le maximum de deux estampilles logiques
func MaxStamp(a, b int) int {
	fmt.Fprintf(os.Stderr, "MaxStamp(%d, %d)\n", a, b)
	if a > b {
		return a
	}
	return b
}

// Retourne le maximum de deux horloges vectorielles
func MaxVectorClock(a []int, b []int) []int {
	var max []int
	

	for i := 0; i < NB_SITES; i++ {
		//fmt.Fprintf(os.Stderr, "MaxStamp(%d, %d)\n", a[i], b[i])
		max = append(max,MaxStamp(a[i], b[i]))
	}
	return max
}


