package pkg

import (
	"github.com/LiveRamp/iabconsent"
)

// CheckConsent parses a TCF v2 consent string and checks for personalized ads and vendor consent.
func CheckConsent(consentStr string, consentVendorID int) (bool, error) {
	consent, err := iabconsent.ParseV2(consentStr)

	if err != nil {
		return false, Errorf(EINVALID, "failed to parse consent string: %w", err)
	}

	// Purposes related to personalized ads: are 1 and 4
	if !consent.PurposeAllowed(1) || !consent.PurposeAllowed(4) {
		return false, nil
	}

	if !consent.VendorAllowed(consentVendorID) {
		return false, nil
	}

	return true, nil
}
