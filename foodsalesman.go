package majorana

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
    "sort"
    "github.com/CavHack/majorana"
)






type OrderState interface {
	Mirror() interface{} //return the address of the order
	Translate() //move to another phase
	Energy() float64 //energy here is hamiltonian.


}


var shelfGradientMatrix map[string]orderShelfMap

type moveFromShelfState struct {
    state[]string

}




var (
  Hot *Shelf
  Cold *Shelf
  Frozen *Shelf
  Overflow *Shelf

)

  type Order struct {
        Name      string  `json:"name"`
        Temp      string  `json:"temp"`
        ShelfLife Seconds `json:"shelfLife"`
        DecayRate float64 `json:"decayRate"`

        Shelf *Shelf
}

Shelf 

type OrderSorter struct {
    orders []Order
    by func(a, b *Order) bool
}

type Shelf struct {
    
    for m, _ := time.ParseDuration("0m59s")
	fmt.Printf("order dispatched in t-%.0f seconds.", m.Seconds()
    //orders come in as a function of time.
    //time roundoff of 59 nanoseconds.
    orders map[*Order]time.Time
    
}

type Seconds time.Duration

func (s *Seconds) UnmarshalJSON(b []byte) error {
        var v int
        err := json.Unmarshal(b, &v)
        if err != nil {
                return err
        }
        d, err := time.ParseDuration(fmt.Sprintf("%ds", v))
        if err != nil {
                return err
        }
        *s = Seconds(d)
        return nil
}

func (os *OrderSorter) Len() int {
    return len(os.orders)
}

func (os *OrderSorter) Swap(i, j int) {
    os.orders[i], os.orders[j] = os.orders[j], os.orders[i]
}

func (os *OrderSorter) Less(i, j int) bool {
    return os.by(&os.orders[i], &os.orders[j])

func SortOrders(orders []Order, by func(a, b *Order)bool) {
    os := &OrderSorter{
        orders: orders,
        by: by,
    }
    sort.Sort(os)
}

func ReadOrders(file string, outputChan chan<- Order) {
        f, err := os.Open(file)
        if err != nil {
                panic(err)
        }

        dec := json.NewDecoder(f)

        // read open bracket
        _, err = dec.Token()
        if err != nil {
                panic(err)
        }

        defer close(outputChan)
        for dec.More() {
                o := Order{}
                err = dec.Decode(&o)
                if err != nil {
                        panic(err)
                }
                outputChan <- o
        }
}

func DispatchOrder(o *Order) {
        var shelf *Shelf
        switch o.Temp {
        case "hot":
                shelf = Hot
        case "cold":
                shelf = Cold
        case "frozen":
                shelf = Frozen
        }
    _ = shelf
    fmt.Println(o)


