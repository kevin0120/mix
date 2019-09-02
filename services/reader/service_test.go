package reader

import (
	"fmt"
	"github.com/ebfe/scard"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	for {
		fmt.Println("List available readers")
		ctx, err := scard.EstablishContext()
		if err != nil {
			time.Sleep(SEARCH_ITV)
			continue
		}

		readers, err := ctx.ListReaders()
		fmt.Println(readers)
		if err != nil {
			fmt.Println("reader lost")
			_ = ctx.Release()
			time.Sleep(SEARCH_ITV)
			continue
		}

		for {
			if len(readers) > 0 {
				index, err := waitUntilCardPresent(ctx, readers)
				if err != nil {
					break
				}

				fmt.Println("Connect to card")
				card, err := ctx.Connect(readers[index], scard.ShareExclusive, scard.ProtocolAny)
				if err != nil {
					time.Sleep(SEARCH_ITV)
					continue
				}

				for {
					var cmd = []byte{0xff, 0xca, 0x00, 0x00, 0x00} // SELECT uid
					rsp, err := card.Transmit(cmd)
					if err != nil {
						// card lost
						fmt.Println("card lost")
						_ = card.Disconnect(scard.ResetCard)
						break
					}

					uid := fmt.Sprintf("%x", string(rsp))
					fmt.Println(uid)

					time.Sleep(SEARCH_ITV)
				}
			}
		}
	}
}
