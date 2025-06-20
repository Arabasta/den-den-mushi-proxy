package io

import "regexp"

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))" +
	"|(?:\u001B\\].*?(?:\u0007|\u001B\\\\))" // OSC
var re = regexp.MustCompile(ansi)

func Stripansi(str string) string {
	return re.ReplaceAllString(str, "")
}
