package osv_advisory

import (
	"errors"
	"strings"

	cvss20 "github.com/pandatix/go-cvss/20"
	cvss30 "github.com/pandatix/go-cvss/30"
	cvss31 "github.com/pandatix/go-cvss/31"
	cvss40 "github.com/pandatix/go-cvss/40"
)

type OsvAdvisory struct {
	ID       string
	Severity *Severity
}

type OsvVulnerability interface {
	GetID() string
	GetSeverities() []OsvSeverity
}

type OsvSeverity interface {
	GetType() string
	GetScore() string
}

type cvss interface {
	BaseScore() float64
}

// Converts cvss40.CVSS40 to the same interface as the package for the other versions.
type _cvss40 cvss40.CVSS40

// BaseScore implements cvss.
func (c *_cvss40) BaseScore() float64 {
	return ((*cvss40.CVSS40)(c)).Score()
}

// AdvisoryFromOsvSeverity converts the severities to a unified OsvAdvisory.
func AdvisoryFromOsvSeverity(vulnerability OsvVulnerability) (OsvAdvisory, error) {
	var advisory OsvAdvisory

	advisory.ID = vulnerability.GetID()
	if advisory.ID == "" {
		return advisory, ErrMissingAdvisoryID
	}

	severities := vulnerability.GetSeverities()
	if len(severities) == 0 {
		return advisory, nil
	}

	var severity Severity
	for _, value := range severities {
		score := value.GetScore()

		if value.GetType() == "Ubuntu" {
			if score == "Negligible" {
				severity = max(SeverityNone, severity)
			} else {
				var tmp Severity

				err := tmp.UnmarshalText([]byte(score))
				if err != nil {
					return advisory, errors.Join(ErrInvalidUbuntuScore, err)
				}

				severity = max(tmp, severity)
			}

			continue
		}

		var cvss cvss
		var err error
		var isV2 bool

		switch value.GetType() {
		case "CVSS_V2":
			isV2 = true
			cvss, err = cvss20.ParseVector(score)
		case "CVSS_V3":
			if strings.HasPrefix(score, "CVSS:3.1") {
				cvss, err = cvss31.ParseVector(score)
			} else {
				cvss, err = cvss30.ParseVector(score)
			}
		case "CVSS_V4":
			var tmp *cvss40.CVSS40
			tmp, err = cvss40.ParseVector(score)
			cvss = (*_cvss40)(tmp)
		default:
			err = ErrInvalidCvssType
		}
		if err != nil {
			return advisory, errors.Join(ErrParseScore, err)
		}

		// Based on cvss package for Python.
		if baseScore := cvss.BaseScore(); baseScore == 0.0 {
			severity = max(SeverityNone, severity)
		} else if baseScore <= 3.9 {
			severity = max(SeverityLow, severity)
		} else if baseScore <= 6.9 {
			severity = max(SeverityMedium, severity)
		} else if baseScore <= 8.9 || isV2 {
			severity = max(SeverityHigh, severity)
		} else {
			severity = max(SeverityCritical, severity)
		}
	}

	advisory.Severity = &severity

	return advisory, nil
}
