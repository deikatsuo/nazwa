package wrapper

// Login struk formulir login
// menyimpan data login
type Login struct {
	Loginid  string `json:"loginid" binding:"required"`
	Password string `json:"password" binding:"required"`
}
