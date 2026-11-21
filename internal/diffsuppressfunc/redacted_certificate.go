package diffsuppressfunc

import (
	"crypto/sha1"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const redactedValue = "[REDACTED]"

var _ schema.SchemaDiffSuppressFunc = RedactedCertificate

func RedactedCertificate(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue != redactedValue {
		return false
	}

	remoteFingerprint, ok := d.Get(k + "_fingerprint").(string)
	if !ok || remoteFingerprint == "" {
		return false
	}

	return strings.EqualFold(certificateFingerprint(newValue), remoteFingerprint)
}

func certificateFingerprint(certificate string) string {
	block, _ := pem.Decode([]byte(strings.TrimSpace(certificate)))
	if block == nil {
		return ""
	}

	sum := sha1.Sum(block.Bytes)
	parts := make([]string, len(sum))
	for i, b := range sum {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, ":")
}
