package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"os/exec"
	"strings"
	"time"
)

var push = flag.Bool("p", false, "push to github")

func main() {
	flag.Parse()

	dateString := dateString()
	markdown := dateString + ".md"
	html := dateString + ".html"

	createFile(markdown)
	createFile(html)

	mostStarred(markdown, html, 10)

	if *push {
		gitPull()
		gitAddAll()
		gitCommit(dateString)
		gitPush()
	}
}

func dateString() string {
	y, m, d := time.Now().Date()
	mStr := fmt.Sprintf("%d", m)
	dStr := fmt.Sprintf("%d", d)
	if m < 10 {
		mStr = fmt.Sprintf("0%d", m)
	}
	if d < 10 {
		dStr = fmt.Sprintf("0%d", d)
	}
	return fmt.Sprintf("%d-%s-%s", y, mStr, dStr)

}

func createFile(filename string) {
	fo, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
}

func mostStarred(filename string, txt string, topnum int) {
	var doc *goquery.Document
	var e error

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	t, err := os.OpenFile(txt, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer t.Close()

	if doc, e = goquery.NewDocument("https://github.com/trending"); e != nil {
		panic(e.Error())
	}

	doc.Find("li.repo-list-item").Each(func(i int, s *goquery.Selection) {
		topnum--
		if topnum < 0 {
			return
		}

		title := s.Find("h3 a").Text()
		description := s.Find("p.repo-list-description").Text()
		url, _ := s.Find("h3 a").Attr("href")
		url = "https://github.com" + url
		meta := s.Find("p.repo-list-meta").Text()
		href, _ := s.Find("p.repo-list-meta a").Attr("href")
		var contributors []string

		s.Find("p.repo-list-meta a img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			title, _ := s.Attr("title")
			height, _ := s.Attr("height")
			width, _ := s.Attr("width")
			contributor := "<img src='" + src + "'"
			contributor += "title='" + title + "'"
			contributor += "height='" + height + "'"
			contributor += "width='" + width + "'"
			contributor += ">"
			contributors = append(contributors, contributor)
		})

		title = trimString(title)
		info := strings.Split(meta, "•")

		for i := 0; i < len(info); i++ {
			info[i] = strings.TrimSpace(info[i])
		}

		mdcontent := "* [" + title + "](" + url + "): "
		mdcontent += strings.Join(info, ", ") + " "
		mdcontent += "<a href='https://github.com" + href + "'>"
		mdcontent += strings.Join(contributors, "") + "\n\r"
		mdcontent += "</a>"
		mdcontent += strings.TrimSpace(description) + "\n\r\r"

		if _, err = f.WriteString(mdcontent); err != nil {
			panic(err)
		}

		content := "<b>" + title + "</b>: " + strings.TrimSuffix(strings.Join(info, ", "), ", Built by")
		content += "<br>" + description + "<br><br>"

		if _, err = t.WriteString(content); err != nil {
			panic(err)
		}

	})
}

func trimString(s string) string {
	var trimS string
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\n' {
			trimS += string(s[i])
		}
	}
	return trimS
}

func gitPull() {
	app := "git"
	arg0 := "pull"
	arg1 := "origin"
	arg2 := "master"
	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}

func gitAddAll() {
	app := "git"
	arg0 := "add"
	arg1 := "."
	cmd := exec.Command(app, arg0, arg1)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}

func gitCommit(date string) {
	app := "git"
	arg0 := "commit"
	arg1 := "-am"
	arg2 := date
	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}

func gitPush() {
	app := "git"
	arg0 := "push"
	arg1 := "origin"
	arg2 := "master"
	cmd := exec.Command(app, arg0, arg1, arg2)
	out, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return
	}

	print(string(out))
}

func SendToNote(user, password, host, to, subject, body string) error {

	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	from := mail.Address{user, user}
	toMail := mail.Address{to, to}
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = toMail.String()
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))

	auth := smtp.PlainAuth("", user, password, host)
	err := smtp.SendMail(
		host+":25",
		auth,
		user,
		[]string{toMail.Address},
		[]byte(message),
	)

	if err != nil {
		panic(err)
	}
	return err
}
