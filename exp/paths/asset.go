package paths

import (
	"github.com/stellar/go/xdr"
)

func (a *Asset) String() string {
	if a.Type == xdr.AssetTypeAssetTypeNative {
		return "native"
	} else {
		return a.Code + " " + a.Issuer
	}
}
