package checker_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lordmangila/status-checker/pkg/checker"
	"github.com/stretchr/testify/assert"
)

func ExampleSite_Validate() {
	site := checker.Site{
		URL: "sample",
	}
	site.Validate()
	fmt.Println(site.Error)
	fmt.Println(site.Valid)

	site = checker.Site{
		URL: "http://google.com",
	}
	site.Validate()
	fmt.Println(site.Valid)

	site = checker.Site{
		URL: "http://www.google.com",
	}
	site.Validate()
	fmt.Println(site.Valid)

	site = checker.Site{
		URL: "http://www.google.com/",
	}
	site.Validate()
	fmt.Println(site.Valid)

	// Output:
	// Invalid URI: sample
	// false
	// true
	// true
	// true
}

func ExampleSite_Marshal() {
	site := checker.Site{
		URL:        "invalidurl",
		StatusCode: 0,
		Active:     false,
		Valid:      false,
		Error:      "Invalid URI: invalidurl",
	}
	fmt.Println(site.Marshal())
	fmt.Println(string(site.Marshal()))

	site = checker.Site{
		URL:        "http://www.google.com",
		StatusCode: 200,
		Active:     true,
		Valid:      true,
		Error:      "",
	}
	fmt.Println(site.Marshal())
	fmt.Println(string(site.Marshal()))

	// Output:
	// [123 34 85 82 76 34 58 34 105 110 118 97 108 105 100 117 114 108 34 44 34 83 116 97 116 117 115 67 111 100 101 34 58 48 44 34 65 99 116 105 118 101 34 58 102 97 108 115 101 44 34 86 97 108 105 100 34 58 102 97 108 115 101 44 34 69 114 114 111 114 34 58 34 73 110 118 97 108 105 100 32 85 82 73 58 32 105 110 118 97 108 105 100 117 114 108 34 125]
	// {"URL":"invalidurl","StatusCode":0,"Active":false,"Valid":false,"Error":"Invalid URI: invalidurl"}
	// [123 34 85 82 76 34 58 34 104 116 116 112 58 47 47 119 119 119 46 103 111 111 103 108 101 46 99 111 109 34 44 34 83 116 97 116 117 115 67 111 100 101 34 58 50 48 48 44 34 65 99 116 105 118 101 34 58 116 114 117 101 44 34 86 97 108 105 100 34 58 116 114 117 101 44 34 69 114 114 111 114 34 58 34 34 125]
	// {"URL":"http://www.google.com","StatusCode":200,"Active":true,"Valid":true,"Error":""}
}

func TestHealthCheck(t *testing.T) {
	site := checker.Site{
		URL: "https://google.com",
	}
	site.HealthCheck()

	site = checker.Site{
		URL: "https://newrelic.com",
	}
	site.HealthCheck()

	site = checker.Site{
		URL: "https://yahoo.com",
	}
	site.HealthCheck()

	site = checker.Site{
		URL: "https://github.com",
	}
	site.HealthCheck()
}

func assertEqual(t *testing.T, site checker.Site) {
	if "/Users/lord" != os.Getenv("HOME") {
		assert.Equal(t, site.Active, true)
		assert.Equal(t, site.StatusCode, 200)
	}
}
