package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"os"
	"strings"
	"testing"
)

func TestAccApsaraStackSlbAclsDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb_acl.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb_acl.default.name}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb_acl.default.name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_slb_acl.default.name}"`,
			"tags":       `{Created = "TF1"}`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_slb_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${apsarastack_slb_acl.default.id}_fake"]`,
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":               `["${apsarastack_slb_acl.default.id}"]`,
			"resource_group_id": `""`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":               `["${apsarastack_slb_acl.default.id}_fake"]`,
			"resource_group_id": fmt.Sprintf(`"%s_fake"`, os.Getenv("APSARASTACK_RESOURCE_GROUP_ID")),
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${apsarastack_slb_acl.default.id}"]`,
			"name_regex": `"${apsarastack_slb_acl.default.name}"`,
			// The resource route tables do not support resource_group_id, so it was set empty.
			"resource_group_id": `""`,
		}),
		fakeConfig: testAccCheckApsaraStackSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":               `["${apsarastack_slb_acl.default.id}_fake"]`,
			"name_regex":        `"${apsarastack_slb_acl.default.name}"`,
			"resource_group_id": `""`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":                     "1",
			"ids.#":                      "1",
			"names.#":                    "1",
			"acls.0.id":                  CHECKSET,
			"acls.0.resource_group_id":   CHECKSET,
			"acls.0.name":                fmt.Sprintf("tf-testAccSlbAclDataSourceBisic-%d", rand),
			"acls.0.ip_version":          "ipv4",
			"acls.0.entry_list.#":        "2",
			"acls.0.related_listeners.#": "0",
			"acls.0.tags.%":              "2",
			"acls.0.tags.Created":        "TF",
			"acls.0.tags.For":            "acceptance test",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":  "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var slbaclsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_slb_acls.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbaclsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, tagsConf, idsConf, resourceGroupIdConf, allConf)
}

func testAccCheckApsaraStackSlbAclsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlbAclDataSourceBisic-%d"
}
variable "ip_version" {
	default = "ipv4"
}

resource "apsarastack_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
    entry = "10.10.10.0/24"
    comment = "first"
  }
  entry_list {
      entry = "168.10.10.0/24"
      comment = "second"
  }
   tags = {
      Created = "TF"
       For     = "acceptance test"
    }
}

data "apsarastack_slb_acls" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
