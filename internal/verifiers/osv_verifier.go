package verifiers

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	int_logger "github.com/New-Horizons-Team/supply-chain-firewall/internal/logger"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/slice_utils"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/target"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/verifiers/osv_advisory"
	"github.com/go-resty/resty/v2"
)

const (
	osvDevQueryUrl      = "https://api.osv.dev/v1/query"
	osvDevVulnUrlPrefix = "https://osv.dev/vulnerability"
	osvDevListUrlPrefix = "https://osv.dev/list"
)

type osvVerifier struct {
	client *resty.Client
}

func init() {
	// The OSV.dev API is sometimes quite slow, hence the generous timeout.
	client := resty.New().
		SetBaseURL(osvDevQueryUrl).
		SetTimeout(time.Second * 10)

	verifier := &osvVerifier{
		client: client,
	}

	verifiers = append(verifiers, verifier)
}

// Name implements FirewallVerifier.
func (*osvVerifier) Name() string {
	return "OsvVerifier"
}

// Verify implements FirewallVerifier.
// OSV.dev advisories with `MAL` IDs are treated as `CRITICAL` findings and all
// others are treated as `WARNING`.  *It is very important to note that most but
// **not all** OSV.dev malicious package advisories have `MAL` IDs.*
func (o *osvVerifier) Verify(ctx context.Context, pkg target.Package) []Finding {
	vulns, err := o.listVulnerabilities(ctx, pkg)
	if err != nil {
		return errorFinding(pkg, err)
	}

	if len(vulns) == 0 {
		return nil
	}

	malicious, nonMalicious, err := listSortedAdvisories(vulns)
	if err != nil {
		logger := int_logger.GetLogger(ctx)
		logger.Warn(ctx, "Verification failed: returning WARNING finding for package %s", pkg)
		return errorFinding(pkg, err)
	}

	fn := func(value osv_advisory.OsvAdvisory) Finding {
		severity := SeverityWarning

		var kind string
		if strings.HasPrefix(value.ID, "MAL") {
			kind = "malicious package "
			severity = SeverityCritical
		}

		var tag string
		if value.Severity != nil {
			tag = fmt.Sprintf("[%s] ", value.Severity.String())
		}

		summary := fmt.Sprintf(
			"An OSV.dev %sadvisory exists for package %s:\n"+
				"  * %s%s/%s",
			kind,
			pkg,
			tag,
			osvDevVulnUrlPrefix,
			value.ID,
		)

		return Finding{
			Summary:  summary,
			Severity: severity,
		}
	}

	ret := make([]Finding, 0, len(malicious)+len(nonMalicious))
	ret = slice_utils.Map(malicious, fn, ret)
	ret = slice_utils.Map(nonMalicious, fn, ret)

	return ret
}

type OsvRequest struct {
	PageToken string     `json:"page_token,omitempty"`
	Version   string     `json:"version"`
	Package   OsvPackage `json:"package"`
}

type OsvPackage struct {
	Name      string           `json:"name"`
	Ecosystem target.Ecosystem `json:"ecosystem"`
}

type OsvResponse struct {
	Vulns    []OsvVulnerability `json:"vulns"`
	NextPage string             `json:"next_page_token"`
}

type OsvVulnerability struct {
	ID       string        `json:"id"`
	Severity []OsvSeverity `json:"severity"`
}

type OsvSeverity struct {
	Type  string `json:"type"`
	Score string `json:"score"`
}

// GetID implements osv_advisory.OsvVulnerability.
func (o OsvVulnerability) GetID() string {
	return o.ID
}

// GetSeverities implements osv_advisory.OsvVulnerability.
func (o OsvVulnerability) GetSeverities() []osv_advisory.OsvSeverity {
	ret := make([]osv_advisory.OsvSeverity, len(o.Severity))

	for i := range o.Severity {
		ret[i] = o.Severity[i]
	}

	return ret
}

// GetType implements osv_advisory.OsvSeverity.
func (s OsvSeverity) GetType() string {
	return s.Type
}

// GetScore implements osv_advisory.OsvSeverity.
func (s OsvSeverity) GetScore() string {
	return s.Score
}

// listVulnerabilities lists every vulnerability associated with the package.
func (o *osvVerifier) listVulnerabilities(ctx context.Context, pkg target.Package) ([]OsvVulnerability, error) {
	var vulns []OsvVulnerability

	payload := OsvRequest{
		Version: pkg.Version,
		Package: OsvPackage{
			Name:      pkg.Name,
			Ecosystem: pkg.Ecosystem,
		},
	}

	for {
		var body OsvResponse

		resp, err := o.client.NewRequest().
			SetHeader("Content-Type", "application/json").
			SetContext(ctx).
			SetBody(payload).
			SetResult(&body).
			Post("")
		if err != nil {
			logger := int_logger.GetLogger(ctx)
			logger.Warn(ctx, "Verification failed: returning WARNING finding for package %s", pkg)

			return nil, err
		} else if resp.IsError() {
			logger := int_logger.GetLogger(ctx)
			logger.Warn(ctx, "Failed to query OSV.dev API: returning WARNING finding for package %s", pkg)

			return nil, fmt.Errorf("OSV.dev failed with status %s (%d)", resp.Status(), resp.StatusCode())
		}

		vulns = append(vulns, body.Vulns...)

		payload.PageToken = body.NextPage
		if payload.PageToken == "" {
			break
		}
	}

	return vulns, nil
}

// listSortedAdvisories converts the list of vunerabilities into two list of advisories,
// one for malicious packages and another for non-malicious ones,
// both sorted in descending vulnerability order.
func listSortedAdvisories(vulns []OsvVulnerability) (malicious, nonMalicious []osv_advisory.OsvAdvisory, err error) {
	for _, value := range vulns {
		var advisory osv_advisory.OsvAdvisory

		if value.ID == "" {
			continue
		}

		advisory, err = osv_advisory.AdvisoryFromOsvSeverity(value)
		if err != nil {
			return
		}

		if strings.HasPrefix(value.ID, "MAL") {
			malicious = append(malicious, advisory)
		} else {
			nonMalicious = append(nonMalicious, advisory)
		}
	}

	sortFn := func(a, b osv_advisory.OsvAdvisory) int {
		if b.Severity == nil {
			return -1
		} else if a.Severity == nil {
			return 1
		}
		return int(*b.Severity - *a.Severity)
	}

	slices.SortStableFunc(malicious, sortFn)
	slices.SortStableFunc(nonMalicious, sortFn)

	return
}

// errorFinding prepares a finding indicating the validation error.
func errorFinding(pkg target.Package, err error) []Finding {
	summary := fmt.Sprintf(
		"Failed to verify package against OSV.dev: %s.\n"+
			"Before proceeding, please check for OSV.dev advisories related to this package.\n"+
			"DO NOT PROCEED if it has an advisory with a MAL ID: it is very likely malicious.\n"+
			"  * %s?q=%s&ecosystem=%s",
		err.Error(),
		osvDevListUrlPrefix,
		pkg.Name,
		pkg.Ecosystem,
	)

	return []Finding{{
		Summary:  summary,
		Severity: SeverityWarning,
	}}
}
