package main

import (
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func Test(t *testing.T) {
	g := Goblin(t)

	//special hook for gomega
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Logger client", func() {
		logger := NewClient("test_set", "test_type")

		g.It("Must be initialised with Log Set and Log Type parameters", func() {
			Expect(logger.LogSet).To(Equal("test_set"))
			Expect(logger.LogType).To(Equal("test_type"))
		})

		var log = map[string]interface{}{
			"message": "hello",
			"code":    400,
			"user_info": map[string]interface{}{
				"name": "tester",
			},
		}

		g.It("Must be able to create a correct message body", func() {
			body := logger.createMessageBody("new_type", log)
			Expect(body["log_set"]).To(Equal(logger.LogSet))
			Expect(body["log_type"]).To(Equal("new_type"))
			Expect(body["log"]).To(Equal(log))
		})

	})

}
