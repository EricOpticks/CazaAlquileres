package mailer

// smtpServer data to smtp server
type SmtpServer struct {
	Host string
	Port string
}
// Address URI to smtp server
func (s *SmtpServer) Address() string {
	return s.Host + ":" + s.Port
}
