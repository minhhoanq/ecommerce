package email

import (
	"testing"

	"github.com/minhhoanq/ecommerce/user_service/config"
	"github.com/stretchr/testify/require"
)

func TestSender(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	cfg, err := config.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(cfg.EmailSenderName, cfg.EmailSenderAddress, cfg.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://techschool.guru">Tech School</a></p>
	`
	to := []string{"minhhoanglost@gmail.com"}

	err = sender.SendMail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
