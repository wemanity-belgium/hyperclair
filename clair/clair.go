package clair

import (
	"strconv"
	"strings"
  "math"

	"github.com/coreos/clair/api/v1"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

var uri string
var priority string
var healthPort int

//Report Reporting Config value
var Report ReportConfig

//VulnerabiliesCounts Total count of vulnerabilities
type VulnerabiliesCounts struct {
  Total  int
  High   int
  Medium int
  Low    int
}

//RelativeCount get the percentage of vulnerabilities of a severity
func (vulnerabilityCount VulnerabiliesCounts) RelativeCount(severity string) float64  {
  var count int
  
  switch severity {
  case "High":
    count = vulnerabilityCount.High
  case "Medium":
    count = vulnerabilityCount.Medium
  case "Low":
    count = vulnerabilityCount.Low
  }
  
  return math.Ceil(float64(count) / float64(vulnerabilityCount.Total) * 100 * 100) / 100
}

//ImageAnalysis Full image analysis
type ImageAnalysis struct {
	Registry  string
	ImageName string
	Tag       string
	Layers    []v1.LayerEnvelope
}

func (imageAnalysis ImageAnalysis) String() string {
	return imageAnalysis.Registry + "/" + imageAnalysis.ImageName + ":" + imageAnalysis.Tag
}

func (imageAnalysis ImageAnalysis) ShortName(l v1.Layer) string {
	return xstrings.Substr(l.Name, 0, 12)
}

func (imageAnalysis ImageAnalysis) CountVulnerabilities(l v1.Layer) int {
	count := 0
	for _, f := range l.Features {
		count += len(f.Vulnerabilities)
	}
	return count
}

//CountAllVulnerabilities Total count of vulnerabilities
func (imageAnalysis ImageAnalysis) CountAllVulnerabilities() VulnerabiliesCounts {
  var result VulnerabiliesCounts;
  result.Total = 0;
  result.High = 0;
  result.Medium = 0;
  result.Low = 0;
  
  for _, l := range imageAnalysis.Layers {
    for _, f := range l.Layer.Features {
      result.Total += len(f.Vulnerabilities)
      for _, v := range f.Vulnerabilities {
        switch v.Severity {
        case "High":
          result.High++
        case "Medium":
          result.Medium++
        case "Low":
          result.Low++
        }
      }
    }
  }
  
  return result;
}

type Vulnerability struct {
	Name, Severity, IntroduceBy, Description, Layer string
}

func (imageAnalysis ImageAnalysis) SortVulnerabilities() []Vulnerability {
	low := []Vulnerability{}
	medium := []Vulnerability{}
	high := []Vulnerability{}
	critical := []Vulnerability{}
	defcon1 := []Vulnerability{}

	for _, l := range imageAnalysis.Layers {
		for _, f := range l.Layer.Features {
			for _, v := range f.Vulnerabilities {
				nv := Vulnerability{
					Name:        v.Name,
					Severity:    v.Severity,
					IntroduceBy: f.Name + ":" + f.Version,
					Description: v.Description,
					Layer:       l.Layer.Name,
				}
				switch strings.ToLower(v.Severity) {
				case "low":
					low = append(low, nv)
				case "medium":
					medium = append(medium, nv)
				case "high":
					high = append(high, nv)
				case "critical":
					critical = append(critical, nv)
				case "defcon1":
					defcon1 = append(defcon1, nv)
				}
			}
		}
	}

	return append(defcon1, append(critical, append(high, append(medium, low...)...)...)...)
}

func fmtURI(u string, port int) {
	uri = u
	if port != 0 {
		uri += ":" + strconv.Itoa(port)
	}
	if !strings.HasSuffix(uri, "/v1") {
		uri += "/v1"
	}
	if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
		uri = "http://" + uri
	}
}

//Config configure Clair from configFile
func Config() {
	fmtURI(viper.GetString("clair.uri"), viper.GetInt("clair.port"))
	priority = viper.GetString("clair.priority")
	healthPort = viper.GetInt("clair.healthPort")
	Report.Path = viper.GetString("clair.report.path")
	Report.Format = viper.GetString("clair.report.format")
}
