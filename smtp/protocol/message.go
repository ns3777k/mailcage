// Based on https://github.com/mailhog/smtp/commit/0e36ecc166f43c24cc7b34c9db695119899f0062
package protocol

type Message struct {
	From string
	To   []string
	Data string
	Helo string
}
