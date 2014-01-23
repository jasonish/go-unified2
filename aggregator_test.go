package unified2

import "fmt"
import "testing"

func TestAggregator(*testing.T) {

	//var aggregator EventAggregator
	//aggregator := new(EventAggregator)
	aggregator := NewEventAggregator()
	fmt.Println(aggregator)
	aggregator.buffer.PushBack(1)
	fmt.Println("len:", aggregator.Len())
	aggregator.buffer.PushBack(2)
	fmt.Println("len:", aggregator.Len())
	aggregator.buffer.PushBack(3)
	fmt.Println("len:", aggregator.Len())

	// aggregator.buffer.Remove(aggregator.buffer.Front())
	// fmt.Println("len:", aggregator.Len())
	// aggregator.buffer.Remove(aggregator.buffer.Front())
	// fmt.Println("len:", aggregator.Len())
	// aggregator.buffer.Remove(aggregator.buffer.Front())
	// fmt.Println("len:", aggregator.Len())

	// aggregator.buffer.Front()

	// fmt.Println("len:", aggregator.Len())

	aggregator.Flush()
	aggregator.Flush()

	// for e := aggregator.buffer.Front(); e != nil; e = e.Next() {
	// 	fmt.Println(e.Value)
	// }
	//aggregator.Flush()

	//aggregator0 := EventAggregator{}
	//aggregator0 := new(EventAggregator)
	// aggregator0 := &EventAggregator{}
	// fmt.Println(aggregator0)
	// aggregator0.buffer.PushBack(1)
	// aggregator0.buffer.PushBack(2)
	// aggregator0.buffer.PushBack(3)
	// for e := aggregator0.buffer.Front(); e != nil; e = e.Next() {
	// 	fmt.Println(e.Value)
	// }
}
