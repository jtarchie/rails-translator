package main

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func TestTranslation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Translation Suite")
}

var _ = Describe("Translations", func() {
	var server *ghttp.Server
	var tempDir string

	BeforeEach(func() {
		server = ghttp.NewServer()

		// Create a temporary directory for test files
		var err error
		tempDir, err = os.MkdirTemp("", "translator_test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()

		// Cleanup: Remove the temporary directory and all its contents
		os.RemoveAll(tempDir)
	})

	Describe("Run function", func() {
		Context("with valid input", func() {
			It("should translate a simple message and save to the specified file", func() {
				// Mock the OpenAI API response
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/chat/completions"),
						ghttp.RespondWithJSONEncoded(200, map[string]interface{}{
							"choices": []map[string]interface{}{
								{
									"message": map[string]string{
										"content": "こんにちは！",
									},
								},
							},
						}),
					),
				)

				// Create a temporary source file
				fromFile := tempDir + "/basic_en.yaml"
				err := os.WriteFile(fromFile, []byte("en:\n  welcome: \"Hello!\""), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				cli := &CLI{
					BaseURL:           server.URL(),
					FromFilename:      fromFile,
					FromLanguage:      "en",
					OpenAIAccessToken: "YOUR_TEST_TOKEN",
					ToFilename:        tempDir + "/jp.yaml",
					ToLanguage:        "Japanese",
				}

				err = cli.Run()
				Expect(err).NotTo(HaveOccurred())

				// Read the translated file and check its contents
				contents, _ := os.ReadFile(tempDir + "/jp.yaml")
				Expect(string(contents)).To(ContainSubstring("jp:\n"))
				Expect(string(contents)).To(ContainSubstring("こんにちは"))
			})
		})
	})
})
