package discourse

import(
	"os"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
	"log"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"discourse": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		log.Println("[ERROR]: ",err)
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T)  {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DISCOURSE_API_USERNAME"); v == "" {
		t.Fatal("DISCOURSE_API_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("DISCOURSE_API_KEY"); v == "" {
		t.Fatal("DISCOURSE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("DISCOURSE_BASE_URL"); v == "" {
		t.Fatal("DISCOURSE_BASE_URL must be set for acceptance tests")
	}
}
