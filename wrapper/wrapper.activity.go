package wrapper

// ActivityTrackerInsert Pengingat aktivitas
type ActivityTrackerInsert struct {
	Subject int    `db:"subject"`
	Type    string `db:"type"`
}
