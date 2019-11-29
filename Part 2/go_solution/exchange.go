package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Exchange struct {
	book map[string]*OrderBook
}

func NewExchange() *Exchange {
	return &Exchange{book: make(map[string]*OrderBook)}
}

func (ex *Exchange) Book(ins string) *OrderBook {
	if ob, ok := ex.book[ins]; ok {
		return ob
	}
	ob := NewOrderBook(ins)
	ex.book[ins] = ob
	return ob
}

func (ex *Exchange) Execute(o *Order) ([]Trade, error) {
	ob := ex.Book(o.Instrument)
	ob.Insert(o)
	return ob.Match()
}

func readOrders() chan *Order {
	orders := make(chan *Order, 10000)

	reader := csv.NewReader(os.Stdin)
	reader.FieldsPerRecord = 4
	reader.Comma = ':'
	reader.ReuseRecord = true

	go func() {
		for {
			defer close(orders)
			rec, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			o, err := ScanOrder(rec)
			if err != nil {
				panic(err)
			}
			orders <- o
		}
	}()

	return orders
}

func executeOrders(orders chan *Order) chan []Trade {
	verbose := false
	exchange := NewExchange()
	trades := make(chan []Trade, 100)
	go func() {
		defer close(trades)
		for o := range orders {
			if verbose {
				fmt.Println(o)
			}
			ts, err := exchange.Execute(o)
			if err != nil {
				panic(err)
			}
			if len(ts) > 0 {
				trades <- ts
			}
			if verbose {
				fmt.Println(exchange.book[o.Instrument])
			}
		}
	}()
	return trades
}

func main() {
	orders := readOrders()
	trades := executeOrders(orders)
	for ts := range trades {
		for t := range ts {
			fmt.Println(t)
		}
	}
}
