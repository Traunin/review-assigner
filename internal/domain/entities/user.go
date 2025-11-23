package entities

type User struct {
	user_id   UserID
	username  string
	is_active bool
}

func NewUser(user_id UserID, username string, is_active bool) (*User, error) {
	user := &User{
		user_id:   user_id,
		username:  username,
		is_active: is_active,
	}

	if err := validate(user); err != nil {
		return nil, err
	}
	return user, nil
}

func validate(user *User) error {
	if user.user_id == "" {
		return ErrUserNoID
	}
	if user.username == "" {
		return ErrUserNoUsername
	}
	return nil
}

func (user *User) IsActive() bool {
	return user.is_active
}

func (user *User) Activate() {
	user.is_active = true
}

func (user *User) Deactivate() {
	user.is_active = false
}

func (user *User) ID() UserID {
	return user.user_id
}

func (user *User) Username() string {
	return user.username
}
