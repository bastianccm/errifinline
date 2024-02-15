package main

import "log"

func error1() error {
	return nil
}

func error2() (*string, error) {
	return nil, nil
}

func main() {
	err := error1()
	if err != nil { // want `inline err assignment in if initializer`
		log.Fatal(err)
	}

	if err := error1(); err != nil {
		log.Fatal(err)
	}

	_, err = error2()
	_, _ = 1, 2
	a, b := 1, 2
	if err != nil { // want `inline err assignment in if initializer`
		log.Fatal(err)
	}
	_, _ = a, b

	_, err = error2() //nolint:ineffassign,staticcheck // test case
	a, b, err = 1, 2, nil
	if err != nil {
		log.Fatal(err)
	}
	_, _ = a, b

	if _, err2 := error2(); err2 != nil {
		log.Fatal(err2)
	}

	_, err2 := error2()
	if err2 != nil { // want `inline err assignment in if initializer`
		log.Fatal(err2)
	}

	if _, err = error2(); err != nil {
		log.Fatal(err)
	}

	something, err := error2()
	if err != nil {
		log.Fatal(err)
	}
	_ = something
}
