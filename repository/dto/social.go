package dto

import (
	"app/repository/dao"
)

type ToggleFollow struct {
	UserID string `binding:"required" json:"userID"`
}

func (body ToggleFollow) Follow(id string) (dao.User, error) {
	me, err := dao.FindUser(id, map[string]interface{}{
		"preload": []string{"Followings"},
	})
	if err != nil {
		return me, err
	}
	user, err := dao.FindUser(body.UserID, map[string]interface{}{
		// "preload": []string{"Fans"},
	})
	if err != nil {
		return me, err
	}
	next := make([]dao.User, 0)
	for _, following := range me.Followings {
		if following.ID != user.ID {
			next = append(next, following)
		}
	}
	me.Follow(user)
	me.Followings = append(next, user)
	return me, nil
}

func (body ToggleFollow) Unfollow(id string) (dao.User, error) {
	me, err := dao.FindUser(id, map[string]interface{}{
		"preload": []string{"Followings"},
	})
	if err != nil {
		return me, err
	}
	user, err := dao.FindUser(body.UserID, map[string]interface{}{
		// "preload": []string{"Fans"},
	})
	if err != nil {
		return me, err
	}
	next := make([]dao.User, 0)
	for _, following := range me.Followings {
		if following.ID != user.ID {
			next = append(next, following)
		}
	}
	me.Unfollow(user)
	me.Followings = next
	return me, nil
}
