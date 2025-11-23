package entities

type User struct {
	user_id   UserID
	username  string
	is_active bool
	team_id   *TeamID
}

func NewUser(
	user_id UserID,
	username string,
	is_active bool,
	team_id *TeamID,
) (*User, error) {
	user := &User{
		user_id:   user_id,
		username:  username,
		is_active: is_active,
		team_id:   team_id,
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

func (user *User) SetActive(is_active bool) {
	user.is_active = is_active
}

func (user *User) ID() UserID {
	return user.user_id
}

func (user *User) Username() string {
	return user.username
}

func (user *User) TeamID() *TeamID {
	return user.team_id
}

func (user *User) SetTeamID(teamID *TeamID) {
	user.team_id = teamID
}

func (user *User) SetUsername(username string) {
	user.username = username
}
