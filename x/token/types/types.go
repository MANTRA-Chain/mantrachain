package types

type NftCollectionCategory string

// TODO: make it works the same way as guard module required privileges kind enum

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
