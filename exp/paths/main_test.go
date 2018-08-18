package paths

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/guregu/null"
	"github.com/stellar/go/support/db"
	"github.com/stellar/go/xdr"
)

func TestPublicOffers(t *testing.T) {
	// Get all offers
	session, err := db.Open("postgres", "postgres://localhost/core?sslmode=disable")
	if err != nil {
		panic(err)
	}
	q := &Q{session}

	var offers []OfferDB
	sql := sq.Select("*").From("offers").OrderBy("price asc")

	err = q.Select(&offers, sql)
	if err != nil {
		panic(err)
	}

	// Construct an exchange state
	ex := &Exchange{}
	ex.Nodes = make(map[Asset]*Node)

	for _, offer := range offers {
		selling := Asset{
			Type:   offer.SellingAssetType,
			Code:   offer.SellingAssetCode.String,
			Issuer: offer.SellingIssuer.String,
		}

		buying := Asset{
			Type:   offer.BuyingAssetType,
			Code:   offer.BuyingAssetCode.String,
			Issuer: offer.BuyingIssuer.String,
		}

		o := Offer{
			Amount: int64(offer.Amount),
			Pricen: offer.Pricen,
			Priced: offer.Priced,
		}

		sellingNode := ex.GetNode(selling)
		buyingNode := ex.GetNode(buying)
		sellingNode.Markets[buyingNode] = append(sellingNode.Markets[buyingNode], o)
	}

	source := Asset{Type: xdr.AssetTypeAssetTypeNative}
	destination := Asset{Type: xdr.AssetTypeAssetTypeCreditAlphanum4, Code: "EURT", Issuer: "GAP5LETOV6YIE62YAM56STDANPRDO7ZFDBGSNHJQIYGGKSMOZAHOOS2S"}
	ex.Find(source, destination)
}

// ----------------------------------------------------------------------------------------------------------------
// From horizon internal:
type Q struct {
	*db.Session
}

// Offer is row of data from the `offers` table from stellar-core
type OfferDB struct {
	SellerID string `db:"sellerid"`
	OfferID  int64  `db:"offerid"`

	SellingAssetType xdr.AssetType `db:"sellingassettype"`
	SellingAssetCode null.String   `db:"sellingassetcode"`
	SellingIssuer    null.String   `db:"sellingissuer"`

	BuyingAssetType xdr.AssetType `db:"buyingassettype"`
	BuyingAssetCode null.String   `db:"buyingassetcode"`
	BuyingIssuer    null.String   `db:"buyingissuer"`

	Amount       xdr.Int64 `db:"amount"`
	Pricen       int32     `db:"pricen"`
	Priced       int32     `db:"priced"`
	Price        float64   `db:"price"`
	Flags        int32     `db:"flags"`
	Lastmodified int32     `db:"lastmodified"`
}
