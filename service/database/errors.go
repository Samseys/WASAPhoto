package database

import "errors"

var ErrAlreadyBanned = errors.New("the user is already banned")
var ErrUsernameAlreadyTaken = errors.New("another user already has this username")
var ErrCommentNotFound = errors.New("the comment doesn't exist")
var ErrAlreadyFollowed = errors.New("the user is already followed")
var ErrPhotoNotFound = errors.New("Photo doesn't exist")
var ErrUserProfileNotFound = errors.New("user doesn't exist")
var ErrAlreadyLiked = errors.New("the photo is already liked")
var ErrNotBanned = errors.New("the user is not banned")
var ErrNotFollowed = errors.New("the user is not followed")
var ErrNotLiked = errors.New("the photo is not liked")
