package i18n

import (
	"io/ioutil"
	"os"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestSprintf(t *testing.T) {
	message.SetString(language.SimplifiedChinese, "%s has %d cat.", "%s有%d只猫。")

	session := RegistPrinter("test_sprintf", language.SimplifiedChinese)
	defer DeletePrinter("test_sprintf")

	if str := session.Sprintf("%s has %d cat.", "tom", 3); str != "tom有3只猫。" {
		t.Errorf("expected %s got %s", "tom有3只猫。", str)
	}
	if str := Session("test_sprintf").Sprintf("%s has %d cat.", "tom", 3); str != "tom有3只猫。" {
		t.Errorf("expected %s got %s", "tom有3只猫。", str)
	}
}

func TestFprintf(t *testing.T) {
	message.SetString(language.SimplifiedChinese, "%s has %d cat.", "%s有%d只猫。")

	session := RegistPrinter("test_fprintf", language.SimplifiedChinese)
	defer DeletePrinter("test_fprintf")

	file, err := ioutil.TempFile("", "test")
	if err != nil {
		panic(err)
	}
	defer os.Remove(file.Name())

	if _, err := session.Fprintf(file, "%s has %d cat.", "tom", 3); err != nil {
		panic(err)
	}
	file.Sync()

	file.Seek(0, 0)
	str, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	if string(str) != "tom有3只猫。" {
		t.Errorf("expected %s got %s", "tom有3只猫。", string(str))
	}
}
