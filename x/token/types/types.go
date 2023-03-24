package types

import fmt "fmt"

type NftCollectionCategory string

const (
	GeneralNftCollectionCat      NftCollectionCategory = "general"
	ArtNftCollectionCat                                = "art"
	CollectiblesNftCollectionCat                       = "collectibles"
	MusicNftCollectionCat                              = "music"
	PhotographyNftCollectionCat                        = "photography"
	SportsNftCollectionCat                             = "sports"
	TradingCardsNftCollectionCat                       = "tradingCards"
	UtilityNftCollectionCat                            = "utility"
)

func ParseNftCollectionCategory(s string) (c NftCollectionCategory, err error) {
	requiredPrivilegesKind := map[NftCollectionCategory]struct{}{
		GeneralNftCollectionCat:      {},
		ArtNftCollectionCat:          {},
		CollectiblesNftCollectionCat: {},
		MusicNftCollectionCat:        {},
		PhotographyNftCollectionCat:  {},
		SportsNftCollectionCat:       {},
		TradingCardsNftCollectionCat: {},
		UtilityNftCollectionCat:      {},
	}

	cap := NftCollectionCategory(s)
	_, ok := requiredPrivilegesKind[cap]
	if !ok {
		return c, fmt.Errorf(`cannot parse:[%s] as nft collection category`, s)
	}
	return cap, nil
}

func (c NftCollectionCategory) String() string {
	return string(c)
}

func (c NftCollectionCategory) Bytes() []byte {
	return []byte(c)
}
