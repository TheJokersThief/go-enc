package cli

import ()

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
