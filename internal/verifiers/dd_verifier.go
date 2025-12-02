package verifiers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	int_logger "github.com/New-Horizons-Team/supply-chain-firewall/internal/logger"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/once_fetch"
	"github.com/New-Horizons-Team/supply-chain-firewall/internal/target"
	"github.com/go-resty/resty/v2"
)

const ddDatasetSamplesUrl = "https://raw.githubusercontent.com/DataDog/malicious-software-packages-dataset/main/samples"

type ddVerifier struct {
	client *resty.Client

	pypi once_fetch.Once[ddManifest]
	npm  once_fetch.Once[ddManifest]
}

// Maps packages to a list of malicious versions,
// or null if the package has always been malicious.
type ddManifest map[string][]string

func init() {
	client := resty.New().
		SetBaseURL(ddDatasetSamplesUrl).
		SetTimeout(time.Second * 5)

	verifier := &ddVerifier{
		client: client,
	}

	verifiers = append(verifiers, verifier)
}

// Name implements FirewallVerifier.
func (*ddVerifier) Name() string {
	return "DatadogMaliciousPackagesVerifier"
}

// Verify implements FirewallVerifier.
func (d *ddVerifier) Verify(ctx context.Context, pkg target.Package) []Finding {
	manifest, err := d.getCachedManifest(ctx, pkg.Ecosystem)
	if err != nil {
		logger := int_logger.GetLogger(ctx)
		logger.Warn(ctx, "Verification failed: returning CRITICAL finding for package %s", pkg)

		summary := fmt.Sprintf("Failed to download list of malicious %s packages: %s", pkg.Ecosystem, err.Error())
		return []Finding{{
			Summary:  summary,
			Severity: SeverityCritical,
		}}
	}

	if _, ok := manifest[pkg.Name]; ok {
		return []Finding{{
			Severity: SeverityCritical,
			Summary:  fmt.Sprintf("Datadog Security Research has determined that package %s is malicious", pkg.Name),
		}}
	}

	return nil
}

// getCachedManifest downloads and caches the requested manifest.
func (d *ddVerifier) getCachedManifest(ctx context.Context, ecosystem target.Ecosystem) (ddManifest, error) {
	switch ecosystem {
	case target.EcosystemPypi:
		return d.pypi.Do(ctx, d.getManifestPypi)
	case target.EcosystemNpm:
		return d.npm.Do(ctx, d.getManifestNpm)
	case target.EcosystemGo:
		return ddManifest{}, nil
	default:
		return nil, ErrInvalidEcosystem
	}
}

// getManifestPypi downloads the PyPI manifest.
func (d *ddVerifier) getManifestPypi(ctx context.Context) (ddManifest, error) {
	return d.getManifest(ctx, target.EcosystemPypi)
}

// getManifestNpm downloads the NPM manifest.
func (d *ddVerifier) getManifestNpm(ctx context.Context) (ddManifest, error) {
	return d.getManifest(ctx, target.EcosystemNpm)
}

// getManifest downloads the requested manifest.
func (d *ddVerifier) getManifest(ctx context.Context, ecosystem target.Ecosystem) (ddManifest, error) {
	// Since content-type is plain/text, we must decode the response manually.
	resp, err := d.client.NewRequest().
		SetContext(ctx).
		Get(fmt.Sprintf("%s/manifest.json", strings.ToLower(string(ecosystem))))
	if err != nil {
		return nil, errors.Join(ErrDownloadDataDogManifest, err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("failed to fetch %s manifest with status %s (%d)", ecosystem, resp.Status(), resp.StatusCode())
	}

	tmp := ddManifest{}
	err = json.Unmarshal(resp.Body(), &tmp)
	if err != nil {
		return nil, errors.Join(ErrDecodeDataDogManifest, err)
	}

	// Normalize every package name.
	normalized := ddManifest{}

	for key, value := range tmp {
		normalized[strings.ToLower(key)] = value
	}

	return normalized, nil
}
