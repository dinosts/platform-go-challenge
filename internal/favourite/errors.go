package favourite

import "errors"

var (
	ErrCouldNotFindFavouriteForAsset = errors.New("Could not find favourite for asset")
	ErrAssetNotFound                 = errors.New("Asset not found")
	ErrCouldNotSaveFavourite         = errors.New("Could not save favourite")
)
