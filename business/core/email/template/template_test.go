package template_test

import (
	"myc-devices-simulator/business/core/email/template"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestRender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		language  string
		templType string
		data      interface{}
		isError   bool
	}{
		{
			name:      "correct render template account-validation - es",
			language:  "es",
			templType: "account-validation.html",
			data: struct {
				Email, ValidationURI string
			}{Email: "email-es@restmail.net", ValidationURI: "/api/v1/user-service/activation/"},
			isError: false,
		},
		{
			name:      "correct render template account-validation - en",
			language:  "en",
			templType: "account-validation.html",
			data: struct {
				Email, ValidationURI string
			}{Email: "email-en@restmail.net", ValidationURI: "/api/v1/user-service/activation/"},
			isError: false,
		},
		{
			name:      "correct render template account-validation - fr",
			language:  "fr",
			templType: "account-validation.html",
			data: struct {
				Email, ValidationURI string
			}{Email: "email-fr@restmail.net", ValidationURI: "/api/v1/user-service/activation/"},
			isError: false,
		},
		{
			name:      "correct render template account-validation - pt",
			language:  "pt",
			templType: "account-validation.html",
			data: struct {
				Email, ValidationURI string
			}{Email: "email-pt@restmail.net", ValidationURI: "/api/v1/user-service/activation/"},
			isError: false,
		},
		{
			name:      "template not exist",
			language:  faker.RandomChoice([]string{"es", "en", "fr", "pt"}),
			templType: "account.html",
			data: struct {
				Email, ValidationURI string
			}{Email: faker.Internet().Email(), ValidationURI: faker.Internet().DomainName()},
			isError: true,
		},
	}

	cupaloy := cupaloy.New(cupaloy.SnapshotSubdirectory("./.snapshots/"))

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			renderTemplate, err := template.Render(test.language, test.templType, test.data)

			if !test.isError {
				assert.Equal(t, nil, err)
				assert.NotEqualf(t, 0, renderTemplate.Len(), "no equal length")

				assert.Equal(t, nil, cupaloy.SnapshotMulti(getSnapshotFileName(test.name), strings.TrimRight(renderTemplate.String(), "\n")))
			}

			if test.isError {
				assert.Error(t, err)
			}
		})
	}
}

func getSnapshotFileName(snapshot string) string {
	return strings.Replace(strings.ToLower(snapshot), " ", "-", -1)
}
