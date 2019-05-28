package notifyicon

func Test() {
	n := &notifyicon.Notification{
		BalloonTipTitle: title,
		BalloonTipText:  message,
		BalloonTipIcon:  notifyicon.BalloonTipIconInfo,
	}

	if err := n.Send(); err != nil {
		t.Error(err)
	}
}
