package redisprovider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// testAccProviders is a map of provider names to provider instances
var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"redis": testAccProvider,
	}
}

func TestResourceRedisStringCreate_Overridable(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceRedisString().Schema, map[string]interface{}{
		"key":         "test-key",
		"value":       "test-value",
		"overridable": true,
	})

	client, mock := redismock.NewClientMock()
	cfg := &ProviderConfig{RedisClient: client}
	ctx := context.Background()

	// Simulate key exists
	mock.ExpectExists("test-key").SetVal(1)
	mock.ExpectSet("test-key", "test-value", 0).SetVal("OK")
	// After Set, the Create function calls Read, which calls Get
	mock.ExpectGet("test-key").SetVal("test-value")

	diags := resourceRedisStringCreate(ctx, d, cfg)
	assert.Len(t, diags, 0, "should not error when overridable is true and key exists")
}

func TestResourceRedisStringCreate_NotOverridable(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceRedisString().Schema, map[string]interface{}{
		"key":         "test-key",
		"value":       "test-value",
		"overridable": false,
	})

	client, mock := redismock.NewClientMock()
	cfg := &ProviderConfig{RedisClient: client}
	ctx := context.Background()

	// Simulate key exists
	mock.ExpectExists("test-key").SetVal(1)

	diags := resourceRedisStringCreate(ctx, d, cfg)
	assert.Len(t, diags, 1, "should error when overridable is false and key exists")
	assert.Contains(t, diags[0].Summary, "already exists")
}

func TestResourceRedisStringCreate_KeyDoesNotExist(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceRedisString().Schema, map[string]interface{}{
		"key":         "test-key",
		"value":       "test-value",
		"overridable": false,
	})

	client, mock := redismock.NewClientMock()
	cfg := &ProviderConfig{RedisClient: client}
	ctx := context.Background()

	// Simulate key does not exist
	mock.ExpectExists("test-key").SetVal(0)
	mock.ExpectSet("test-key", "test-value", 0).SetVal("OK")
	// After Set, the Create function calls Read, which calls Get
	mock.ExpectGet("test-key").SetVal("test-value")

	diags := resourceRedisStringCreate(ctx, d, cfg)
	assert.Len(t, diags, 0, "should not error when key does not exist")
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("REDIS_URL"); v == "" {
		t.Fatal("REDIS_URL must be set for acceptance tests")
	}
}

func testAccProviderConfig() string {
	return fmt.Sprintf(`
provider "redis" {
  redis_url = "%s"
}
`, os.Getenv("REDIS_URL"))
}

func TestAccRedisString_Basic(t *testing.T) {
	resourceName := "redis_string.test"
	key := "tf_acc_test_key"
	value := "tf_acc_test_value"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig() + fmt.Sprintf(`
resource "redis_string" "test" {
  key   = "%s"
  value = "%s"
}
`, key, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "value", value),
				),
			},
		},
	})
}

func TestAccRedisString_Overridable(t *testing.T) {
	resourceName := "redis_string.test"
	key := "tf_acc_test_overridable"
	value1 := "value1"
	value2 := "value2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create with value1
			{
				Config: testAccProviderConfig() + fmt.Sprintf(`
resource "redis_string" "test" {
  key   = "%s"
  value = "%s"
}
`, key, value1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "value", value1),
				),
			},
			// Try to create with same key, overridable = false (should fail)
			{
				Config: testAccProviderConfig() + fmt.Sprintf(`
resource "redis_string" "test" {
  key   = "%s"
  value = "%s"
}
`, key, value2),
				ExpectError: regexp.MustCompile("already exists"),
			},
			// Create with overridable = true (should override)
			{
				Config: testAccProviderConfig() + fmt.Sprintf(`
resource "redis_string" "test" {
  key         = "%s"
  value       = "%s"
  overridable = true
}
`, key, value2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "value", value2),
					resource.TestCheckResourceAttr(resourceName, "overridable", "true"),
				),
			},
		},
	})
}

func TestAccRedisString_UpdateAndDelete(t *testing.T) {
	resourceName := "redis_string.test"
	key := "tf_acc_test_update"
	value1 := "v1"
	value2 := "v2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccProviderConfig() + fmt.Sprintf(`
resource "redis_string" "test" {
  key   = "%s"
  value = "%s"
}
`, key, value1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "value", value1),
				),
			},
			// Update
			{
				Config: testAccProviderConfig() + fmt.Sprintf(`
resource "redis_string" "test" {
  key   = "%s"
  value = "%s"
}
`, key, value2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "value", value2),
				),
			},
		},
	})
} 