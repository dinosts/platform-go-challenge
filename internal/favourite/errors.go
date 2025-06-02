package favourite

import "errors"

var (
	ErrCouldNotFindFavouriteForAsset = errors.New("Could not find favourite for asset")
	ErrAssetNotFound                 = errors.New("Asset not found")
	ErrCouldNotSaveFavourite         = errors.New("Could not save favourite")
	ErrFavouriteNotUnderGivenUser    = errors.New("Favourite is not under given user")
	ErrFavouriteNotFound             = errors.New("Favourite not found.")
)
