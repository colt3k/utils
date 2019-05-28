package freedesktop

func Test() {
	title := "title"
	message := "message"
	n := &freedesktop.Notification{
		Summary:       title,
		Body:          message,
		ExpireTimeout: 500,
		AppIcon:       "utilities-terminal",
	}

	if err := n.Send(); err != nil {
		t.Error(err)
	}
}
